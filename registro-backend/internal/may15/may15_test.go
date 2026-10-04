package may15

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMay15Repo struct {
	mock.Mock
}

func (m *MockMay15Repo) GetByClassAndYear(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error) {
	args := m.Called(ctx, classID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClassMay15Document), args.Error(1)
}

func (m *MockMay15Repo) Upsert(ctx context.Context, doc *ClassMay15Document) error {
	args := m.Called(ctx, doc)
	doc.ID = "doc-123"
	return args.Error(0)
}

func (m *MockMay15Repo) Publish(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error) {
	args := m.Called(ctx, classID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClassMay15Document), args.Error(1)
}

func TestMay15Service_GetDocument_DefaultDraft(t *testing.T) {
	repo := new(MockMay15Repo)
	svc := NewService(repo)

	repo.On("GetByClassAndYear", mock.Anything, "class-5a", "2025/2026").Return(nil, nil).Once()

	doc, err := svc.GetDocument(context.Background(), "class-5a", "2025/2026")
	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, StatusBozza, doc.Status)
	assert.Equal(t, "class-5a", doc.ClassID)
}

func TestMay15Service_SaveDocument(t *testing.T) {
	repo := new(MockMay15Repo)
	svc := NewService(repo)

	req := SaveMay15Request{
		AcademicYear:       "2025/2026",
		Status:             StatusApprovatoCdC,
		ClassPresentation:  "Classe composta da 22 alunni, clima cooperativo",
		TeachingContinuity: "Continuità garantita su Italiano, Matematica, Lingua Inglese",
		PCTOPathways:       "Percorsi svolti presso aziende del territorio (media 110 ore)",
		ExamSimulations:    "Svolta 1ª simulazione prova nazionale il 15 Febbraio",
		EvaluationRubrics:  "Griglie ministeriali O.M. Esami di Stato",
		CLILModules:        "Modulo di Scienze Naturali in lingua inglese (30 ore)",
	}

	repo.On("Upsert", mock.Anything, mock.MatchedBy(func(d *ClassMay15Document) bool {
		return d.ClassID == "class-5a" && d.Status == StatusApprovatoCdC && d.CLILModules != ""
	})).Return(nil).Once()

	doc, err := svc.SaveDocument(context.Background(), "school-1", "class-5a", req)
	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, "doc-123", doc.ID)
	assert.Equal(t, "Modulo di Scienze Naturali in lingua inglese (30 ore)", doc.CLILModules)
}

func TestMay15Service_PublishDocument(t *testing.T) {
	repo := new(MockMay15Repo)
	svc := NewService(repo)

	now := time.Now()
	expected := &ClassMay15Document{
		ID:           "doc-123",
		ClassID:      "class-5a",
		AcademicYear: "2025/2026",
		Status:       StatusPubblicato,
		PublishedAt:  &now,
	}

	repo.On("Publish", mock.Anything, "class-5a", "2025/2026").Return(expected, nil).Once()

	doc, err := svc.PublishDocument(context.Background(), "class-5a", "2025/2026")
	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, StatusPubblicato, doc.Status)
}
