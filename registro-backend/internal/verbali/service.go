package verbali

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrUnauthorized                      = errors.New("unauthorized action on verbali")
	ErrNotFound                          = errors.New("verbale or meeting not found")
	ErrVerbaleLocked                     = errors.New("il verbale è firmato ed è in sola lettura: modifiche bloccate")
	ErrOnlyCoordinatorOrSecretaryCanEdit = errors.New("solo il docente coordinatore di classe o il verbalista possono modificare la bozza del verbale")
	ErrDraftHiddenFromPrincipal          = errors.New("il verbale è in bozza e sarà visibile alla Dirigente Scolastica solo dopo l'avvenuta firma")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("verbali.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

// CreateMeeting creates a new council or general school meeting.
func (s *Service) CreateMeeting(ctx context.Context, actorID, schoolID string, req CreateMeetingRequest) (*CouncilMeeting, error) {
	if schoolID == "" {
		schoolID = s.repo.ResolveSchoolID(ctx, actorID)
	}
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: use YYYY-MM-DD")
	}

	if req.ClassID != "" {
		belongs, err := s.repo.ClassBelongsToSchool(ctx, req.ClassID, schoolID)
		if err != nil {
			return nil, fmt.Errorf("error verifying class tenant association: %w", err)
		}
		if !belongs {
			return nil, fmt.Errorf("class %s does not belong to school %s", req.ClassID, schoolID)
		}
	}

	meetingType := req.MeetingType
	if meetingType == "" {
		meetingType = "consiglio_classe"
	}

	m := &CouncilMeeting{
		SchoolID:    schoolID,
		ClassID:     req.ClassID,
		MeetingType: meetingType,
		Title:       req.Title,
		Date:        d,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Agenda:      req.Agenda,
		CreatedBy:   actorID,
	}

	if err := s.repo.CreateMeeting(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) ListMeetings(ctx context.Context, schoolID, classID string) ([]*CouncilMeeting, error) {
	return s.repo.ListMeetings(ctx, schoolID, classID)
}

func (s *Service) GetMeeting(ctx context.Context, meetingID string) (*CouncilMeeting, error) {
	return s.repo.GetMeetingByID(ctx, meetingID)
}

// CreateVerbale creates a new draft verbale for a meeting.
// Can be created by the designated secretary, coordinator, or school admin.
func (s *Service) CreateVerbale(ctx context.Context, actorID, actorRole string, req CreateVerbaleRequest) (*MeetingVerbale, error) {
	meeting, coordinatorID, err := s.repo.GetMeetingWithCoordinator(ctx, req.MeetingID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Check permissions to create draft
	isSec := (req.SecretaryID != nil && *req.SecretaryID == actorID)
	isPres := (req.PresidentID != nil && *req.PresidentID == actorID)
	isCoord := (coordinatorID != "" && coordinatorID == actorID)
	isAdmin := (actorRole == "admin" || actorRole == "superadmin")

	if !isSec && !isPres && !isCoord && !isAdmin {
		// Only coordinator or secretary can initiate the draft verbale
		return nil, ErrOnlyCoordinatorOrSecretaryCanEdit
	}

	v := &MeetingVerbale{
		MeetingID:   meeting.ID,
		Title:       req.Title,
		Content:     req.Content,
		SecretaryID: req.SecretaryID,
		PresidentID: req.PresidentID,
		IsPublished: req.IsPublished,
		IsSigned:    false,
		Status:      "draft",
		CanEdit:     true,
	}

	if err := s.repo.CreateVerbale(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdateVerbale updates a draft verbale.
// Rule 1: Only the class coordinator or the verbalista (secretary) can edit the draft.
// Rule 2: If the verbale is already signed, it is strictly READ-ONLY for everyone.
func (s *Service) UpdateVerbale(ctx context.Context, actorID, actorRole, verbaleID string, req UpdateVerbaleRequest) (*MeetingVerbale, error) {
	v, err := s.repo.GetVerbaleByID(ctx, verbaleID, actorID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Rule 2: Post-signing tamper-proof lock
	if v.IsSigned || v.Status == "signed" {
		return nil, ErrVerbaleLocked
	}

	// Rule 1: Only class coordinator or verbalista can edit
	_, coordinatorID, err := s.repo.GetMeetingWithCoordinator(ctx, v.MeetingID)
	if err != nil {
		return nil, err
	}

	isSec := (v.SecretaryID != nil && *v.SecretaryID == actorID)
	isPres := (v.PresidentID != nil && *v.PresidentID == actorID)
	isCoord := (coordinatorID != "" && coordinatorID == actorID)
	isAdmin := (actorRole == "admin" || actorRole == "superadmin")

	if !isSec && !isPres && !isCoord && !isAdmin {
		return nil, ErrOnlyCoordinatorOrSecretaryCanEdit
	}

	v.Title = req.Title
	v.Content = req.Content
	v.SecretaryID = req.SecretaryID
	v.PresidentID = req.PresidentID
	v.IsPublished = req.IsPublished

	if err := s.repo.UpdateVerbale(ctx, v); err != nil {
		return nil, err
	}
	v.CanEdit = true
	return v, nil
}

// DeleteVerbale removes an unsigned draft verbale. Signed verbali can NEVER be deleted.
func (s *Service) DeleteVerbale(ctx context.Context, actorID, actorRole, verbaleID string) error {
	v, err := s.repo.GetVerbaleByID(ctx, verbaleID, actorID)
	if err != nil {
		return ErrNotFound
	}

	if v.IsSigned || v.Status == "signed" {
		return ErrVerbaleLocked
	}

	_, coordinatorID, err := s.repo.GetMeetingWithCoordinator(ctx, v.MeetingID)
	if err != nil {
		return err
	}

	isSec := (v.SecretaryID != nil && *v.SecretaryID == actorID)
	isCoord := (coordinatorID != "" && coordinatorID == actorID)
	isAdmin := (actorRole == "admin" || actorRole == "superadmin")

	if !isSec && !isCoord && !isAdmin {
		return ErrOnlyCoordinatorOrSecretaryCanEdit
	}

	return s.repo.DeleteVerbale(ctx, verbaleID)
}

// GetVerbale retrieves a single verbale.
// Rule 3: Hidden from the school principal (Dirigente Scolastica) until it is signed!
func (s *Service) GetVerbale(ctx context.Context, verbaleID, userID, role string) (*MeetingVerbale, error) {
	v, err := s.repo.GetVerbaleByID(ctx, verbaleID, userID)
	if err != nil {
		return nil, ErrNotFound
	}

	isPrincipal := (role == "principal" || role == "vice_principal")
	if isPrincipal && !v.IsSigned {
		return nil, ErrDraftHiddenFromPrincipal
	}

	// Calculate CanEdit flag
	if v.IsSigned || v.Status == "signed" {
		v.CanEdit = false // Read-only for everyone post-signing
	} else {
		_, coordinatorID, _ := s.repo.GetMeetingWithCoordinator(ctx, v.MeetingID)
		isSec := (v.SecretaryID != nil && *v.SecretaryID == userID)
		isPres := (v.PresidentID != nil && *v.PresidentID == userID)
		isCoord := (coordinatorID != "" && coordinatorID == userID)
		isAdmin := (role == "admin" || role == "superadmin")
		v.CanEdit = (isSec || isPres || isCoord || isAdmin)
	}

	return v, nil
}

// ListVerbali returns verbali for a specific meeting.
// If caller is Dirigente Scolastica, drafts are completely omitted until signed.
func (s *Service) ListVerbali(ctx context.Context, meetingID, userID, role string) ([]*MeetingVerbale, error) {
	onlySigned := (role == "principal" || role == "vice_principal")
	list, err := s.repo.ListVerbali(ctx, meetingID, userID, onlySigned)
	if err != nil {
		return nil, err
	}

	_, coordinatorID, _ := s.repo.GetMeetingWithCoordinator(ctx, meetingID)
	for _, v := range list {
		if v.IsSigned || v.Status == "signed" {
			v.CanEdit = false
		} else {
			isSec := (v.SecretaryID != nil && *v.SecretaryID == userID)
			isPres := (v.PresidentID != nil && *v.PresidentID == userID)
			isCoord := (coordinatorID != "" && coordinatorID == userID)
			isAdmin := (role == "admin" || role == "superadmin")
			v.CanEdit = (isSec || isPres || isCoord || isAdmin)
		}
	}
	return list, nil
}

// ListAllVerbali returns all verbali for a school / class.
// For Dirigente Scolastica, returns ONLY signed verbali.
func (s *Service) ListAllVerbali(ctx context.Context, schoolID, classID, userID, role string) ([]*MeetingVerbale, error) {
	onlySigned := (role == "principal" || role == "vice_principal")
	list, err := s.repo.ListAllVerbali(ctx, schoolID, classID, userID, onlySigned)
	if err != nil {
		return nil, err
	}

	for _, v := range list {
		if v.IsSigned || v.Status == "signed" {
			v.CanEdit = false
		} else {
			_, coordinatorID, _ := s.repo.GetMeetingWithCoordinator(ctx, v.MeetingID)
			isSec := (v.SecretaryID != nil && *v.SecretaryID == userID)
			isPres := (v.PresidentID != nil && *v.PresidentID == userID)
			isCoord := (coordinatorID != "" && coordinatorID == userID)
			isAdmin := (role == "admin" || role == "superadmin")
			v.CanEdit = (isSec || isPres || isCoord || isAdmin)
		}
	}
	return list, nil
}

// SignVerbale allows the verbalista or coordinator/president to sign the verbale.
// Once signed, the verbale is locked and becomes visible to the Dirigente Scolastica.
func (s *Service) SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error {
	v, err := s.repo.GetVerbaleByID(ctx, verbaleID, userID)
	if err != nil {
		return ErrNotFound
	}

	_, coordinatorID, _ := s.repo.GetMeetingWithCoordinator(ctx, v.MeetingID)
	isSec := (v.SecretaryID != nil && *v.SecretaryID == userID)
	isPres := (v.PresidentID != nil && *v.PresidentID == userID)
	isCoord := (coordinatorID != "" && coordinatorID == userID)

	if !isSec && !isPres && !isCoord {
		return errors.New("unauthorized: non sei il verbalista o il coordinatore/presidente autorizzato a firmare questo verbale")
	}

	// Record the signature with IP and timestamp
	if err := s.repo.SignVerbale(ctx, verbaleID, userID, ipAddress); err != nil {
		return err
	}

	// Mark as officially signed, tamper-proof, and visible to Principal
	return s.repo.MarkVerbaleSigned(ctx, verbaleID)
}

func (s *Service) GetSignatures(ctx context.Context, verbaleID string) ([]VerbaleSignature, error) {
	return s.repo.GetSignatures(ctx, verbaleID)
}

// ----------------- Template Service Methods -----------------

func isPrincipalOrAdmin(role string) bool {
	switch strings.ToLower(role) {
	case "principal", "vice_principal", "admin", "superadmin", "dsga", "collaboratore_ds", "assistente_amministrativo", "secretary":
		return true
	default:
		return false
	}
}

func (s *Service) ResolveSchoolID(ctx context.Context, userID string) string {
	return s.repo.ResolveSchoolID(ctx, userID)
}

func (s *Service) CreateTemplate(ctx context.Context, actorID, actorRole, schoolID string, req CreateTemplateRequest) (*MeetingVerbaleTemplate, error) {
	if !isPrincipalOrAdmin(actorRole) {
		return nil, ErrUnauthorized
	}
	if schoolID == "" {
		schoolID = s.repo.ResolveSchoolID(ctx, actorID)
	}
	if schoolID == "" {
		return nil, fmt.Errorf("school_id required")
	}

	t := &MeetingVerbaleTemplate{
		SchoolID:        schoolID,
		Title:           req.Title,
		MeetingType:     req.MeetingType,
		Description:     req.Description,
		DefaultAgenda:   req.DefaultAgenda,
		TemplateContent: req.TemplateContent,
		CreatedBy:       actorID,
	}

	if err := s.repo.CreateTemplate(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTemplates(ctx context.Context, schoolID, meetingType string) ([]*MeetingVerbaleTemplate, error) {
	return s.repo.ListTemplates(ctx, schoolID, meetingType)
}

func (s *Service) GetTemplateByID(ctx context.Context, id string) (*MeetingVerbaleTemplate, error) {
	return s.repo.GetTemplateByID(ctx, id)
}

func (s *Service) UpdateTemplate(ctx context.Context, actorRole, schoolID, id string, req UpdateTemplateRequest) (*MeetingVerbaleTemplate, error) {
	if !isPrincipalOrAdmin(actorRole) {
		return nil, ErrUnauthorized
	}

	t, err := s.repo.GetTemplateByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	if schoolID != "" && t.SchoolID != "" && t.SchoolID != schoolID {
		return nil, ErrUnauthorized
	}

	t.Title = req.Title
	t.MeetingType = req.MeetingType
	t.Description = req.Description
	t.DefaultAgenda = req.DefaultAgenda
	t.TemplateContent = req.TemplateContent

	if err := s.repo.UpdateTemplate(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, actorRole, schoolID, id string) error {
	if !isPrincipalOrAdmin(actorRole) {
		return ErrUnauthorized
	}
	return s.repo.DeleteTemplate(ctx, id, schoolID)
}
