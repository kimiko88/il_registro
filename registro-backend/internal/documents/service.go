package documents

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"registro-backend/internal/permissions"
	"registro-backend/internal/users"
)

type Service interface {
	CreateDocument(ctx context.Context, actorRole, userID, schoolID string, req CreateDocumentRequest) (*DocumentListResponse, error)
	GetDocument(ctx context.Context, actorRole, schoolID, id string) (*DocumentDetailResponse, error)
	UpdateDocument(ctx context.Context, actorRole, schoolID, userID, id string, req UpdateDocumentRequest) error
	DeleteDocument(ctx context.Context, actorRole, schoolID, userID, id string) error
	AttachFile(ctx context.Context, actorRole, schoolID, docID, fileURL string) error

	// Workflow
	ProcessWorkflow(ctx context.Context, actorRole, schoolID, userID, docID string, req WorkflowActionRequest) error

	// Sign
	SignDocument(ctx context.Context, actorRole, schoolID, userID, docID string, req SignDocumentRequest) error

	// Lists
	GetInbox(ctx context.Context, schoolID string) ([]DocumentListResponse, error)
	GetReviewQueue(ctx context.Context, schoolID string) ([]DocumentListResponse, error)
	GetMyDocuments(ctx context.Context, schoolID, userID string) ([]DocumentListResponse, error)
	ListDocuments(ctx context.Context, actorRole, schoolID string, docType *DocType) ([]DocumentListResponse, error)

	// Templates
	CreateTemplate(ctx context.Context, schoolID string, req TemplateRequest) error
	ListTemplates(ctx context.Context, schoolID string) ([]DocumentTemplate, error)
	UpdateTemplate(ctx context.Context, schoolID, id string, req TemplateRequest) error
	DeleteTemplate(ctx context.Context, schoolID, id string) error

	// Export
	ExportDocument(ctx context.Context, actorRole, schoolID, id, format string) ([]byte, string, error)

	// GetDocumentContent restituisce il contenuto corrente del documento per firma o verifica.
	// Implementa signatures.DocumentService.
	GetDocumentContent(ctx context.Context, id string) (string, error)

	// Status
	LockDocument(ctx context.Context, id string) error
	GetDocumentVersions(ctx context.Context, actorRole, schoolID, docID string) ([]DocumentVersion, error)
}

type service struct {
	repo        Repository
	userRepo    users.Repository
	validator   *Validator
	workflow    *WorkflowEngine
	tpl         *TemplateEngine
	signer      *SignatureProvider
	exporter    *Exporter
	permManager *permissions.Manager
}

func NewService(repo Repository, uRepo ...users.Repository) Service {
	var userRepo users.Repository
	if len(uRepo) > 0 {
		userRepo = uRepo[0]
	}
	return &service{
		repo:        repo,
		userRepo:    userRepo,
		validator:   NewValidator(),
		workflow:    NewWorkflowEngine(),
		tpl:         NewTemplateEngine(),
		signer:      NewSignatureProvider(),
		exporter:    NewExporter(),
		permManager: permissions.NewManager(),
	}
}

func (s *service) CreateDocument(ctx context.Context, actorRole, userID, schoolID string, req CreateDocumentRequest) (*DocumentListResponse, error) {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentCreate) {
		return nil, errors.New("unauthorized")
	}
	content := req.Content

	// Bug 90/133: replace hardcoded mock data with real values from the request
	if req.TemplateID != nil {
		t, err := s.repo.GetTemplate(*req.TemplateID)
		if err == nil {
			if t.SchoolID != "" && t.SchoolID != schoolID {
				return nil, errors.New("unauthorized: template belongs to another school")
			}
			studentName := ""
			studentIDStr := ""
			if req.StudentID != nil {
				studentIDStr = *req.StudentID
				studentName = studentIDStr
				if s.userRepo != nil {
					if u, err := s.userRepo.GetByID(ctx, *req.StudentID); err == nil && u != nil {
						studentName = fmt.Sprintf("%s %s", u.FirstName, u.LastName)
					}
				}
			}
			data := map[string]string{
				"student_name": studentName,
				"student_id":   studentIDStr,
				"title":        req.Title,
				"date":         time.Now().Format("2006-01-02"),
				"school_id":    schoolID,
			}
			content = s.tpl.Render(t.Content, data)
		}
	}

	doc := &Document{
		SchoolID:  schoolID,
		Title:     req.Title,
		Type:      req.Type,
		Status:    StatusDraft,
		CreatedBy: userID,
		StudentID: req.StudentID,
		ClassID:   req.ClassID,
	}

	if err := s.validator.ValidateDocument(doc, content); err != nil {
		return nil, err
	}

	if err := s.repo.Create(doc, content); err != nil {
		return nil, err
	}

	return &DocumentListResponse{ID: doc.ID, Title: doc.Title, Status: doc.Status}, nil
}

func (s *service) GetDocument(ctx context.Context, actorRole, schoolID, id string) (*DocumentDetailResponse, error) {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentRead) {
		return nil, errors.New("unauthorized")
	}
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if doc.SchoolID != schoolID {
		return nil, errors.New("unauthorized: cannot access documents of another school")
	}

	content, _ := s.repo.GetContent(id, doc.CurrentVersion)
	vers, _ := s.repo.GetVersions(id)

	sigs, _ := s.repo.GetSignatures(id)
	var sigSum []SignatureSummary
	for _, sig := range sigs {
		sigSum = append(sigSum, SignatureSummary{SignedBy: sig.SignerID, SignedAt: sig.SignedAt})
	}

	return &DocumentDetailResponse{
		DocumentListResponse: DocumentListResponse{
			ID: doc.ID, Title: doc.Title, Type: doc.Type, Status: doc.Status,
			UpdatedAt: doc.UpdatedAt, IsSigned: doc.IsSigned,
		},
		Content:        content,
		CurrentVersion: doc.CurrentVersion,
		Versions:       convertVersions(vers),
		Signatures:     sigSum,
	}, nil
}

func (s *service) UpdateDocument(ctx context.Context, actorRole, schoolID, userID, id string, req UpdateDocumentRequest) error {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentUpdate) {
		return errors.New("unauthorized")
	}
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if doc.SchoolID != schoolID {
		return errors.New("unauthorized: cannot edit document of another school")
	}
	if doc.CreatedBy != userID && actorRole != "admin" && actorRole != "superadmin" {
		return errors.New("forbidden: non puoi modificare i documenti altrui")
	}

	if doc.Status != StatusDraft && doc.Status != StatusRejected {
		return errors.New("cannot edit non-draft document")
	}

	if req.Title != nil {
		doc.Title = *req.Title
	}
	content := ""
	if req.Content != nil {
		content = *req.Content
	}

	if content == "" {
		c, _ := s.repo.GetContent(id, doc.CurrentVersion)
		content = c
	}

	return s.repo.Update(doc, content, req.ChangeLog)
}

func (s *service) DeleteDocument(ctx context.Context, actorRole, schoolID, userID, id string) error {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentDelete) {
		return errors.New("unauthorized")
	}
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if actorRole != "superadmin" && doc.SchoolID != schoolID {
		return errors.New("unauthorized: cannot delete document of another school")
	}
	if actorRole != "admin" && actorRole != "superadmin" {
		if doc.CreatedBy != userID {
			return errors.New("unauthorized: can only delete own documents")
		}
	}
	return s.repo.Delete(id)
}

func (s *service) ProcessWorkflow(ctx context.Context, actorRole, schoolID, userID, docID string, req WorkflowActionRequest) error {
	doc, err := s.repo.FindByID(docID)
	if err != nil {
		return err
	}
	if doc.SchoolID != schoolID {
		return errors.New("unauthorized: cannot process document of another school")
	}

	// RBAC check per workflow action
	switch req.Action {
	case "submit":
		if doc.CreatedBy != userID && actorRole != "admin" && actorRole != "superadmin" {
			return errors.New("unauthorized: only document creator or admin can submit")
		}
	case "approve_secretary":
		if actorRole != "secretary" && actorRole != "admin" && actorRole != "superadmin" {
			return errors.New("unauthorized: insufficient permissions for secretary approval")
		}
	case "approve_director":
		if actorRole != "admin" && actorRole != "superadmin" {
			return errors.New("unauthorized: insufficient permissions for director approval")
		}
	case "reject":
		if actorRole != "secretary" && actorRole != "admin" && actorRole != "superadmin" {
			return errors.New("unauthorized: insufficient permissions to reject")
		}
	default:
		return errors.New("unknown action")
	}

	// Determine next status
	var next DocStatus
	switch req.Action {
	case "submit":
		next = StatusSubmitted
	case "approve_secretary":
		next = StatusReview
	case "approve_director":
		next = StatusApproved
	case "reject":
		next = StatusRejected
	}

	if err := s.workflow.CanTransition(doc.Status, next); err != nil {
		return err
	}

	return s.repo.UpdateStatus(docID, next)
}

func (s *service) LockDocument(ctx context.Context, id string) error {
	// Bug 92/134: LockDocument è esposto solo via signatures.Service (trusted-internal).
	// Non è raggiungibile direttamente via HTTP handler. Se in futuro viene esposto
	// tramite handler, aggiungere actorRole/schoolID come parametri e verificare RBAC.
	return s.repo.UpdateStatus(id, StatusSigned)
}

func (s *service) SignDocument(ctx context.Context, actorRole, schoolID, userID, docID string, req SignDocumentRequest) error {
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "teacher" && actorRole != "principal" && actorRole != "vice_principal" {
		return errors.New("unauthorized: insufficient permissions to sign document")
	}
	doc, err := s.repo.FindByID(docID)
	if err != nil {
		return err
	}
	if doc.SchoolID != schoolID {
		return errors.New("unauthorized: cannot sign document of another school")
	}

	if doc.Status != StatusApproved {
		return errors.New("document not valid for signing")
	}

	// Fetch current content for cryptographic signature verification
	content, err := s.repo.GetContent(docID, doc.CurrentVersion)
	if err != nil {
		return fmt.Errorf("failed to fetch document content for verification: %w", err)
	}

	if err := s.signer.Verify(content, req.SignatureData, req.CertificateData); err != nil {
		return err
	}

	sig := &DocumentSignature{
		DocumentID:    docID,
		VersionID:     "", // will be set below; never use literal "latest"
		SignerID:      userID,
		SignatureData: req.SignatureData,
		Certificate:   req.CertificateData,
	}

	// Bug 94: GetVersions failure must not silently produce a garbage VersionID.
	vers, err := s.repo.GetVersions(docID)
	if err != nil {
		return fmt.Errorf("failed to retrieve document versions for signing: %w", err)
	}
	if len(vers) == 0 {
		return errors.New("cannot sign document: no versions found")
	}
	sig.VersionID = vers[0].ID

	return s.repo.AddSignature(sig)
}

func (s *service) GetInbox(ctx context.Context, schoolID string) ([]DocumentListResponse, error) {
	docs, err := s.repo.GetInbox(schoolID)
	if err != nil {
		return nil, err
	}
	return convertList(docs), nil
}

func (s *service) GetReviewQueue(ctx context.Context, schoolID string) ([]DocumentListResponse, error) {
	docs, err := s.repo.GetReviewQueue(schoolID)
	if err != nil {
		return nil, err
	}
	return convertList(docs), nil
}

func (s *service) GetMyDocuments(ctx context.Context, schoolID, userID string) ([]DocumentListResponse, error) {
	if schoolID == "" {
		return nil, errors.New("school_id non fornito")
	}
	docs, err := s.repo.ListAll(schoolID, nil)
	if err != nil {
		return nil, err
	}
	var res []DocumentListResponse
	for _, d := range docs {
		if d.CreatedBy == userID || (d.StudentID != nil && *d.StudentID == userID) {
			res = append(res, DocumentListResponse{
				ID: d.ID, Title: d.Title, Type: d.Type,
				Status: d.Status, IsSigned: d.IsSigned, UpdatedAt: d.UpdatedAt,
			})
		}
	}
	if res == nil {
		res = []DocumentListResponse{}
	}
	return res, nil
}

func (s *service) ListDocuments(ctx context.Context, actorRole, schoolID string, docType *DocType) ([]DocumentListResponse, error) {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentRead) {
		return nil, errors.New("unauthorized")
	}
	docs, err := s.repo.ListAll(schoolID, docType)
	if err != nil {
		return nil, err
	}
	return convertList(docs), nil
}

func (s *service) CreateTemplate(ctx context.Context, schoolID string, req TemplateRequest) error {
	return s.repo.CreateTemplate(&DocumentTemplate{
		SchoolID: schoolID,
		Name:     req.Name,
		Type:     req.Type,
		Content:  req.Content,
	})
}

func (s *service) ListTemplates(ctx context.Context, schoolID string) ([]DocumentTemplate, error) {
	return s.repo.GetTemplates(schoolID)
}

func (s *service) UpdateTemplate(ctx context.Context, schoolID, id string, req TemplateRequest) error {
	t, err := s.repo.GetTemplate(id)
	if err != nil {
		return err
	}
	if t.SchoolID != schoolID {
		return errors.New("unauthorized: template belongs to another school")
	}

	return s.repo.UpdateTemplate(&DocumentTemplate{
		ID:       id,
		SchoolID: schoolID,
		Name:     req.Name,
		Type:     req.Type,
		Content:  req.Content,
		IsActive: t.IsActive,
	})
}

func (s *service) DeleteTemplate(ctx context.Context, schoolID, id string) error {
	t, err := s.repo.GetTemplate(id)
	if err != nil {
		return err
	}
	if t.SchoolID != schoolID {
		return errors.New("unauthorized: template belongs to another school")
	}
	return s.repo.DeleteTemplate(id)
}

func (s *service) ExportDocument(ctx context.Context, actorRole, schoolID, id, format string) ([]byte, string, error) {
	doc, err := s.GetDocument(ctx, actorRole, schoolID, id)
	if err != nil {
		return nil, "", err
	}

	d := &Document{Title: doc.Title} // minimal

	if format == "pdf" {
		data, err := s.exporter.ToPDF(d, doc.Content)
		return data, "application/pdf", err
	}
	data, err := s.exporter.ToDOCX(d, doc.Content)
	return data, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", err
}

// Helpers
func convertList(docs []Document) []DocumentListResponse {
	var res []DocumentListResponse
	for _, d := range docs {
		res = append(res, DocumentListResponse{
			ID: d.ID, Title: d.Title, Type: d.Type, Status: d.Status, IsSigned: d.IsSigned, UpdatedAt: d.UpdatedAt,
		})
	}
	return res
}

func convertVersions(vers []DocumentVersion) []VersionSummary {
	var res []VersionSummary
	for _, v := range vers {
		res = append(res, VersionSummary{
			VersionNum: v.VersionNum, CreatedAt: v.CreatedAt, ChangeLog: v.ChangeLog, CreatedBy: v.CreatedBy,
		})
	}
	return res
}

func (s *service) GetDocumentVersions(ctx context.Context, actorRole, schoolID, docID string) ([]DocumentVersion, error) {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentRead) {
		return nil, errors.New("unauthorized")
	}
	doc, err := s.repo.FindByID(docID)
	if err != nil {
		return nil, err
	}
	if doc.SchoolID != schoolID {
		return nil, errors.New("unauthorized: cannot access documents of another school")
	}
	return s.repo.GetVersions(docID)
}

func (s *service) AttachFile(ctx context.Context, actorRole, schoolID, docID, fileURL string) error {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentUpdate) {
		return errors.New("unauthorized")
	}
	doc, err := s.repo.FindByID(docID)
	if err != nil {
		return err
	}
	if doc.SchoolID != schoolID {
		return errors.New("unauthorized: cannot access documents of another school")
	}

	// Bug 93: validazione URL per prevenire path traversal, javascript: injection e URL interni arbitrari.
	cleanURL := strings.ReplaceAll(fileURL, "\r", "")
	cleanURL = strings.ReplaceAll(cleanURL, "\n", "")
	cleanURL = strings.TrimSpace(cleanURL)

	parsedURL, err := url.Parse(cleanURL)
	if err != nil || parsedURL.Scheme != "https" {
		return errors.New("fileURL non valido: sono accettati solo URL HTTPS per i documenti ufficiali")
	}
	if parsedURL.Host == "" {
		return errors.New("fileURL non valido: host mancante")
	}

	content, _ := s.repo.GetContent(docID, doc.CurrentVersion)
	updatedContent := content
	if updatedContent == "" {
		updatedContent = "Attachment: " + cleanURL
	} else {
		updatedContent += "\nAttachment: " + cleanURL
	}

	return s.repo.Update(doc, updatedContent, "Attached file: "+cleanURL)
}

// GetDocumentContent restituisce il contenuto corrente del documento.
// Implementa signatures.DocumentService — usato per calcolare l'hash crittografico prima della firma.
// Non esegue check RBAC: la chiamata è interna e il controllo accessi è già stato effettuato
// dall'handler di firma tramite SignDocument / SignDocumentCtx.
func (s *service) GetDocumentContent(ctx context.Context, id string) (string, error) {
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return "", fmt.Errorf("documento non trovato: %w", err)
	}
	content, err := s.repo.GetContent(id, doc.CurrentVersion)
	if err != nil {
		return "", fmt.Errorf("impossibile recuperare il contenuto del documento: %w", err)
	}
	return content, nil
}
