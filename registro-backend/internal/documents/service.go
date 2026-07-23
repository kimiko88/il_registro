package documents

import (
	"context"
	"errors"

	"registro-backend/internal/permissions"
)

type Service interface {
	CreateDocument(ctx context.Context, actorRole, userID, schoolID string, req CreateDocumentRequest) (*DocumentListResponse, error)
	GetDocument(ctx context.Context, actorRole, schoolID, id string) (*DocumentDetailResponse, error)
	UpdateDocument(ctx context.Context, actorRole, userID, id string, req UpdateDocumentRequest) error
	DeleteDocument(ctx context.Context, actorRole, userID, id string) error
	AttachFile(ctx context.Context, actorRole, schoolID, docID, fileURL string) error

	// Workflow
	ProcessWorkflow(ctx context.Context, userID, docID string, req WorkflowActionRequest) error

	// Sign
	SignDocument(ctx context.Context, userID, docID string, req SignDocumentRequest) error

	// Lists
	GetInbox(ctx context.Context, schoolID string) ([]DocumentListResponse, error)
	GetReviewQueue(ctx context.Context, schoolID string) ([]DocumentListResponse, error)
	GetMyDocuments(ctx context.Context, userID string) ([]DocumentListResponse, error)
	ListDocuments(ctx context.Context, actorRole, schoolID string, docType *DocType) ([]DocumentListResponse, error)

	// Templates
	CreateTemplate(ctx context.Context, schoolID string, req TemplateRequest) error
	ListTemplates(ctx context.Context, schoolID string) ([]DocumentTemplate, error)
	UpdateTemplate(ctx context.Context, schoolID, id string, req TemplateRequest) error
	DeleteTemplate(ctx context.Context, schoolID, id string) error

	// Export
	ExportDocument(ctx context.Context, actorRole, schoolID, id, format string) ([]byte, string, error)

	// Status
	LockDocument(ctx context.Context, id string) error
	GetDocumentVersions(ctx context.Context, actorRole, schoolID, docID string) ([]DocumentVersion, error)
}

type service struct {
	repo        Repository
	validator   *Validator
	workflow    *WorkflowEngine
	tpl         *TemplateEngine
	signer      *SignatureProvider
	exporter    *Exporter
	permManager *permissions.Manager
}

func NewService(repo Repository) Service {
	return &service{
		repo:        repo,
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

	// Template usage
	if req.TemplateID != nil {
		t, err := s.repo.GetTemplate(*req.TemplateID)
		if err == nil {
			// Mock data map
			data := map[string]string{"student_name": "Mario Rossi", "date": "2025-01-01"}
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

func (s *service) UpdateDocument(ctx context.Context, actorRole, userID, id string, req UpdateDocumentRequest) error {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentUpdate) {
		return errors.New("unauthorized")
	}
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return err
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

	// Assuming content is always needed for version update or we fetch old?
	// Simplified: Update requires content
	if content == "" {
		c, _ := s.repo.GetContent(id, doc.CurrentVersion)
		content = c
	}

	return s.repo.Update(doc, content, req.ChangeLog)
}

func (s *service) DeleteDocument(ctx context.Context, actorRole, userID, id string) error {
	if !s.permManager.HasPermission(actorRole, permissions.DocumentDelete) {
		return errors.New("unauthorized")
	}
	doc, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "director" && actorRole != "principal" {
		if doc.CreatedBy != userID {
			return errors.New("unauthorized: can only delete own documents")
		}
	}
	return s.repo.Delete(id)
}

func (s *service) ProcessWorkflow(ctx context.Context, userID, docID string, req WorkflowActionRequest) error {
	doc, err := s.repo.FindByID(docID)
	if err != nil {
		return err
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
		next = StatusRejected // or Draft
	default:
		return errors.New("unknown action")
	}

	if err := s.workflow.CanTransition(doc.Status, next); err != nil {
		return err
	}

	return s.repo.UpdateStatus(docID, next)
}

func (s *service) LockDocument(ctx context.Context, id string) error {
	// StatusSigned or StatusLocked
	return s.repo.UpdateStatus(id, StatusSigned)
}

func (s *service) SignDocument(ctx context.Context, userID, docID string, req SignDocumentRequest) error {
	doc, err := s.repo.FindByID(docID)
	if err != nil {
		return err
	}

	if doc.Status != StatusApproved {
		return errors.New("document not valid for signing")
	}

	if err := s.signer.Verify(docID, req.SignatureData, req.CertificateData); err != nil {
		return err
	}

	sig := &DocumentSignature{
		DocumentID:    docID,
		VersionID:     "latest", // Needs fetch version ID logic
		SignerID:      userID,
		SignatureData: req.SignatureData,
		Certificate:   req.CertificateData,
	}
	// Fetch version ID for current
	// Omitted for brevity, assume "GetVersions" gives IDs
	vers, _ := s.repo.GetVersions(docID)
	if len(vers) > 0 {
		sig.VersionID = vers[0].ID
	}

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

func (s *service) GetMyDocuments(ctx context.Context, userID string) ([]DocumentListResponse, error) {
	// Filter by CreatedBy or StudentID depending on role?
	// Assume teacher: CreatedBy
	// Stub
	return nil, nil
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

	content, _ := s.repo.GetContent(docID, doc.CurrentVersion)
	updatedContent := content
	if updatedContent == "" {
		updatedContent = "Attachment: " + fileURL
	} else {
		updatedContent += "\nAttachment: " + fileURL
	}

	return s.repo.Update(doc, updatedContent, "Attached file: "+fileURL)
}
