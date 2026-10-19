package documents

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ExtDocRepoMock struct {
	mock.Mock
}

func (m *ExtDocRepoMock) Create(doc *Document, initialContent string) error {
	args := m.Called(doc, initialContent)
	doc.ID = "doc-new"
	return args.Error(0)
}

func (m *ExtDocRepoMock) Update(doc *Document, newContent, changeLog string) error {
	args := m.Called(doc, newContent, changeLog)
	return args.Error(0)
}

func (m *ExtDocRepoMock) UpdateStatus(docID string, status DocStatus) error {
	args := m.Called(docID, status)
	return args.Error(0)
}

func (m *ExtDocRepoMock) Delete(docID string) error {
	args := m.Called(docID)
	return args.Error(0)
}

func (m *ExtDocRepoMock) ListAll(schoolID string, docType *DocType) ([]Document, error) {
	args := m.Called(schoolID, docType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Document), args.Error(1)
}

func (m *ExtDocRepoMock) FindByID(id string) (*Document, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Document), args.Error(1)
}

func (m *ExtDocRepoMock) GetContent(docID string, version int) (string, error) {
	args := m.Called(docID, version)
	return args.String(0), args.Error(1)
}

func (m *ExtDocRepoMock) GetVersions(docID string) ([]DocumentVersion, error) {
	args := m.Called(docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentVersion), args.Error(1)
}

func (m *ExtDocRepoMock) FindByClass(classID string) ([]Document, error) {
	args := m.Called(classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Document), args.Error(1)
}

func (m *ExtDocRepoMock) FindByStudent(studentID string) ([]Document, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Document), args.Error(1)
}

func (m *ExtDocRepoMock) GetInbox(schoolID string) ([]Document, error) {
	args := m.Called(schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Document), args.Error(1)
}

func (m *ExtDocRepoMock) GetReviewQueue(schoolID string) ([]Document, error) {
	args := m.Called(schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Document), args.Error(1)
}

func (m *ExtDocRepoMock) AddSignature(sig *DocumentSignature) error {
	args := m.Called(sig)
	return args.Error(0)
}

func (m *ExtDocRepoMock) GetSignatures(docID string) ([]DocumentSignature, error) {
	args := m.Called(docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentSignature), args.Error(1)
}

func (m *ExtDocRepoMock) GetTemplates(schoolID string) ([]DocumentTemplate, error) {
	args := m.Called(schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocumentTemplate), args.Error(1)
}

func (m *ExtDocRepoMock) GetTemplate(id string) (*DocumentTemplate, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DocumentTemplate), args.Error(1)
}

func (m *ExtDocRepoMock) CreateTemplate(tpl *DocumentTemplate) error {
	args := m.Called(tpl)
	return args.Error(0)
}

func (m *ExtDocRepoMock) UpdateTemplate(tpl *DocumentTemplate) error {
	args := m.Called(tpl)
	return args.Error(0)
}

func (m *ExtDocRepoMock) DeleteTemplate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestDocumentsService_Extended(t *testing.T) {
	repo := new(ExtDocRepoMock)
	svc := NewService(repo)
	ctx := context.Background()

	// 1. GetDocument guards and success
	_, err := svc.GetDocument(ctx, "teacher", "sch-1", "")
	assert.ErrorContains(t, err, "required")

	_, err = svc.GetDocument(ctx, "invalid_role", "sch-1", "doc-1") // invalid_role does not have DocumentRead
	assert.Error(t, err)

	repo.On("FindByID", "doc-missing").Return(nil, errors.New("sql: no rows")).Once()
	_, err = svc.GetDocument(ctx, "teacher", "sch-1", "doc-missing")
	assert.Error(t, err)

	repo.On("FindByID", "doc-nil").Return(nil, nil).Once()
	_, err = svc.GetDocument(ctx, "teacher", "sch-1", "doc-nil")
	assert.ErrorContains(t, err, "not found")

	repo.On("FindByID", "doc-wrong-school").Return(&Document{ID: "doc-wrong-school", SchoolID: "other-school"}, nil).Once()
	_, err = svc.GetDocument(ctx, "teacher", "sch-1", "doc-wrong-school")
	assert.ErrorContains(t, err, "unauthorized")

	validDoc := &Document{
		ID:             "doc-1",
		SchoolID:       "sch-1",
		Title:          "Verbale 1",
		Type:           TypeGeneric,
		Status:         StatusDraft,
		CurrentVersion: 1,
		IsSigned:       false,
	}
	repo.On("FindByID", "doc-1").Return(validDoc, nil).Once()
	repo.On("GetContent", "doc-1", 1).Return("test content", nil).Once()
	repo.On("GetVersions", "doc-1").Return([]DocumentVersion{{ID: "v1", VersionNum: 1, Content: "test content"}}, nil).Once()
	repo.On("GetSignatures", "doc-1").Return([]DocumentSignature{{ID: "s1", SignerID: "usr-1", SignedAt: time.Now()}}, nil).Once()

	detail, err := svc.GetDocument(ctx, "teacher", "sch-1", "doc-1")
	assert.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Equal(t, "Verbale 1", detail.Title)
	assert.Equal(t, "test content", detail.Content)

	// 2. UpdateDocument guards and checks
	err = svc.UpdateDocument(ctx, "admin", "sch-1", "usr-1", "", UpdateDocumentRequest{})
	assert.ErrorContains(t, err, "required")

	repo.On("FindByID", "doc-approved").Return(&Document{ID: "doc-approved", SchoolID: "sch-1", Status: StatusApproved, CreatedBy: "usr-1"}, nil).Once()
	newTitle := "Nuovo Titolo"
	err = svc.UpdateDocument(ctx, "admin", "sch-1", "usr-1", "doc-approved", UpdateDocumentRequest{Title: &newTitle})
	assert.ErrorContains(t, err, "cannot edit non-draft")

	repo.On("FindByID", "doc-other-author").Return(&Document{ID: "doc-other-author", SchoolID: "sch-1", Status: StatusDraft, CreatedBy: "usr-other"}, nil).Once()
	err = svc.UpdateDocument(ctx, "secretary", "sch-1", "usr-1", "doc-other-author", UpdateDocumentRequest{Title: &newTitle})
	assert.ErrorContains(t, err, "forbidden")

	repo.On("FindByID", "doc-editable").Return(&Document{ID: "doc-editable", SchoolID: "sch-1", Status: StatusDraft, CreatedBy: "usr-1", CurrentVersion: 1}, nil).Once()
	repo.On("GetContent", "doc-editable", 1).Return("original content", nil).Once()
	repo.On("Update", mock.Anything, "original content", "Aggiornato titolo").Return(nil).Once()
	err = svc.UpdateDocument(ctx, "admin", "sch-1", "usr-1", "doc-editable", UpdateDocumentRequest{Title: &newTitle, ChangeLog: "Aggiornato titolo"})
	assert.NoError(t, err)

	// 3. DeleteDocument guards
	err = svc.DeleteDocument(ctx, "admin", "sch-1", "usr-1", "")
	assert.ErrorContains(t, err, "required")

	repo.On("FindByID", "doc-delete-other").Return(&Document{ID: "doc-delete-other", SchoolID: "sch-1", CreatedBy: "usr-2"}, nil).Once()
	err = svc.DeleteDocument(ctx, "secretary", "sch-1", "usr-1", "doc-delete-other")
	assert.ErrorContains(t, err, "unauthorized")

	repo.On("FindByID", "doc-delete-ok").Return(&Document{ID: "doc-delete-ok", SchoolID: "sch-1", CreatedBy: "usr-1"}, nil).Once()
	repo.On("Delete", "doc-delete-ok").Return(nil).Once()
	err = svc.DeleteDocument(ctx, "admin", "sch-1", "usr-1", "doc-delete-ok")
	assert.NoError(t, err)

	// 4. ProcessWorkflow actions
	err = svc.ProcessWorkflow(ctx, "teacher", "sch-1", "usr-1", "", WorkflowActionRequest{Action: "submit"})
	assert.ErrorContains(t, err, "required")

	repo.On("FindByID", "doc-wf").Return(&Document{ID: "doc-wf", SchoolID: "sch-1", CreatedBy: "usr-1", Status: StatusDraft}, nil).Once()
	repo.On("UpdateStatus", "doc-wf", StatusSubmitted).Return(nil).Once()
	err = svc.ProcessWorkflow(ctx, "teacher", "sch-1", "usr-1", "doc-wf", WorkflowActionRequest{Action: "submit"})
	assert.NoError(t, err)

	repo.On("FindByID", "doc-wf-sub").Return(&Document{ID: "doc-wf-sub", SchoolID: "sch-1", CreatedBy: "usr-1", Status: StatusSubmitted}, nil).Once()
	repo.On("UpdateStatus", "doc-wf-sub", StatusReview).Return(nil).Once()
	err = svc.ProcessWorkflow(ctx, "secretary", "sch-1", "sec-1", "doc-wf-sub", WorkflowActionRequest{Action: "approve_secretary"})
	assert.NoError(t, err)

	repo.On("FindByID", "doc-wf-rev").Return(&Document{ID: "doc-wf-rev", SchoolID: "sch-1", CreatedBy: "usr-1", Status: StatusReview}, nil).Once()
	repo.On("UpdateStatus", "doc-wf-rev", StatusApproved).Return(nil).Once()
	err = svc.ProcessWorkflow(ctx, "principal", "sch-1", "p-1", "doc-wf-rev", WorkflowActionRequest{Action: "approve_director"})
	assert.NoError(t, err)

	// 5. Template operations
	err = svc.CreateTemplate(ctx, "", TemplateRequest{Name: "T1", Content: "C1"})
	assert.ErrorContains(t, err, "schoolID is required")

	err = svc.CreateTemplate(ctx, "sch-1", TemplateRequest{Name: "", Content: "C1"})
	assert.ErrorContains(t, err, "template name is required")

	repo.On("CreateTemplate", mock.Anything).Return(nil).Once()
	err = svc.CreateTemplate(ctx, "sch-1", TemplateRequest{Name: "T1", Content: "C1", Type: TypeGeneric})
	assert.NoError(t, err)

	repo.On("GetTemplates", "sch-1").Return([]DocumentTemplate{{ID: "tpl-1", Name: "T1"}}, nil).Once()
	tpls, err := svc.ListTemplates(ctx, "sch-1")
	assert.NoError(t, err)
	assert.Len(t, tpls, 1)

	err = svc.UpdateTemplate(ctx, "sch-1", "", TemplateRequest{Name: "T1"})
	assert.ErrorContains(t, err, "template id is required")

	repo.On("GetTemplate", "tpl-1").Return(&DocumentTemplate{ID: "tpl-1", SchoolID: "sch-1", IsActive: true}, nil).Once()
	repo.On("UpdateTemplate", mock.Anything).Return(nil).Once()
	err = svc.UpdateTemplate(ctx, "sch-1", "tpl-1", TemplateRequest{Name: "T1 Updated", Content: "C1 Updated"})
	assert.NoError(t, err)

	err = svc.DeleteTemplate(ctx, "sch-1", "")
	assert.ErrorContains(t, err, "template id is required")

	repo.On("GetTemplate", "tpl-1").Return(&DocumentTemplate{ID: "tpl-1", SchoolID: "sch-1"}, nil).Once()
	repo.On("DeleteTemplate", "tpl-1").Return(nil).Once()
	err = svc.DeleteTemplate(ctx, "sch-1", "tpl-1")
	assert.NoError(t, err)

	// 6. ReviewQueue & MyDocuments & ListDocuments
	repo.On("GetReviewQueue", "sch-1").Return([]Document{{ID: "doc-rev"}}, nil).Once()
	revDocs, err := svc.GetReviewQueue(ctx, "sch-1")
	assert.NoError(t, err)
	assert.Len(t, revDocs, 1)

	repo.On("ListAll", "sch-1", (*DocType)(nil)).Return([]Document{{ID: "doc-all"}}, nil).Once()
	allDocs, err := svc.ListDocuments(ctx, "teacher", "sch-1", nil)
	assert.NoError(t, err)
	assert.Len(t, allDocs, 1)

	// 7. LockDocument & GetDocumentContent & GetDocumentVersions
	repo.On("UpdateStatus", "doc-1", StatusSigned).Return(nil).Once()
	err = svc.LockDocument(ctx, "doc-1")
	assert.NoError(t, err)

	repo.On("FindByID", "doc-1").Return(&Document{ID: "doc-1", CurrentVersion: 2}, nil).Once()
	repo.On("GetContent", "doc-1", 2).Return("content v2", nil).Once()
	content, err := svc.GetDocumentContent(ctx, "doc-1")
	assert.NoError(t, err)
	assert.Equal(t, "content v2", content)

	repo.On("FindByID", "doc-1").Return(&Document{ID: "doc-1", SchoolID: "sch-1"}, nil).Once()
	repo.On("GetVersions", "doc-1").Return([]DocumentVersion{{ID: "v1", VersionNum: 1}}, nil).Once()
	versions, err := svc.GetDocumentVersions(ctx, "teacher", "sch-1", "doc-1")
	assert.NoError(t, err)
	assert.Len(t, versions, 1)

	// 8. ExportDocument
	repo.On("FindByID", "doc-1").Return(validDoc, nil).Twice()
	repo.On("GetContent", "doc-1", 1).Return("test export content", nil).Twice()
	repo.On("GetVersions", "doc-1").Return([]DocumentVersion{}, nil).Twice()
	repo.On("GetSignatures", "doc-1").Return([]DocumentSignature{}, nil).Twice()

	pdfBytes, mimePDF, err := svc.ExportDocument(ctx, "teacher", "sch-1", "doc-1", "pdf")
	assert.NoError(t, err)
	assert.Equal(t, "application/pdf", mimePDF)
	assert.NotEmpty(t, pdfBytes)

	docxBytes, mimeDOCX, err := svc.ExportDocument(ctx, "teacher", "sch-1", "doc-1", "docx")
	assert.NoError(t, err)
	assert.Contains(t, mimeDOCX, "wordprocessingml")
	assert.NotEmpty(t, docxBytes)

	// 9. GetMyDocuments & AttachFile
	repo.On("ListAll", "sch-1", (*DocType)(nil)).Return([]Document{
		{ID: "doc-mine", CreatedBy: "usr-1", Title: "Mio Doc"},
		{ID: "doc-not-mine", CreatedBy: "usr-other", Title: "Altro Doc"},
	}, nil).Once()
	myDocs, err := svc.GetMyDocuments(ctx, "sch-1", "usr-1")
	assert.NoError(t, err)
	assert.Len(t, myDocs, 1)
	assert.Equal(t, "Mio Doc", myDocs[0].Title)

	repo.On("FindByID", "doc-1").Return(validDoc, nil).Once()
	repo.On("GetContent", "doc-1", 1).Return("existing content", nil).Once()
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	err = svc.AttachFile(ctx, "admin", "sch-1", "doc-1", "https://storage.example.com/file.pdf")
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}
