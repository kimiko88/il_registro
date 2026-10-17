package unit

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"registro-backend/internal/personnel_desk"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDeskRepo struct {
	mu       sync.RWMutex
	requests map[string]*personnel_desk.DeskRequest
}

func newMockDeskRepo() *mockDeskRepo {
	return &mockDeskRepo{
		requests: make(map[string]*personnel_desk.DeskRequest),
	}
}

func (m *mockDeskRepo) Create(ctx context.Context, r *personnel_desk.DeskRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()

	// store a deep copy
	copyReq := *r
	m.requests[r.ID] = &copyReq
	return nil
}

func (m *mockDeskRepo) GetByID(ctx context.Context, schoolID, id string) (*personnel_desk.DeskRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, exists := m.requests[id]
	if !exists || r.SchoolID != schoolID {
		return nil, errors.New("richiesta non trovata")
	}
	copyReq := *r
	return &copyReq, nil
}

func (m *mockDeskRepo) List(ctx context.Context, schoolID, applicantID, status string) ([]*personnel_desk.DeskRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*personnel_desk.DeskRequest
	for _, r := range m.requests {
		if r.SchoolID != schoolID {
			continue
		}
		if applicantID != "" && r.ApplicantID != applicantID {
			continue
		}
		if status != "" && string(r.Status) != status {
			continue
		}
		copyReq := *r
		result = append(result, &copyReq)
	}
	return result, nil
}

func (m *mockDeskRepo) Update(ctx context.Context, r *personnel_desk.DeskRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.requests[r.ID]; !exists {
		return errors.New("richiesta non trovata")
	}
	r.UpdatedAt = time.Now()
	copyReq := *r
	m.requests[r.ID] = &copyReq
	return nil
}

func (m *mockDeskRepo) Delete(ctx context.Context, schoolID, id, applicantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, exists := m.requests[id]
	if !exists || r.SchoolID != schoolID {
		return errors.New("richiesta non trovata")
	}
	if applicantID != "" && r.ApplicantID != applicantID {
		return errors.New("non autorizzato")
	}
	delete(m.requests, id)
	return nil
}

func TestPersonnelDesk_DraftToSubmit_Lifecycle(t *testing.T) {
	ctx := context.Background()
	repo := newMockDeskRepo()
	svc := personnel_desk.NewService(repo)

	schoolID := "school-1"
	applicantID := "teacher-1"

	// 1. Create Draft Request
	input := personnel_desk.CreateDeskRequestInput{
		Category:    personnel_desk.CatPermessoBreve,
		StartDate:   "2026-10-15",
		EndDate:     "2026-10-15",
		Hours:       2,
		Description: "Visita specialistica programmata",
		SubmitNow:   false,
	}

	created, err := svc.CreateRequest(ctx, schoolID, applicantID, input)
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, personnel_desk.StatusDraft, created.Status)
	assert.Equal(t, applicantID, created.ApplicantID)

	// 2. Unauthorized user cannot submit
	_, err = svc.SubmitRequest(ctx, schoolID, created.ID, "teacher-2")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "solo il richiedente può inviare la richiesta")

	// 3. Applicant submits successfully
	submitted, err := svc.SubmitRequest(ctx, schoolID, created.ID, applicantID)
	require.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusSubmitted, submitted.Status)

	// 4. Cannot re-submit already submitted request
	_, err = svc.SubmitRequest(ctx, schoolID, created.ID, applicantID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "la richiesta è già stata inviata")
}

func TestPersonnelDesk_ThreeStepApproval_Success(t *testing.T) {
	ctx := context.Background()
	repo := newMockDeskRepo()
	svc := personnel_desk.NewService(repo)

	schoolID := "school-1"
	applicantID := "teacher-1"

	// Create and submit directly
	created, err := svc.CreateRequest(ctx, schoolID, applicantID, personnel_desk.CreateDeskRequestInput{
		Category:    personnel_desk.CatFerie,
		StartDate:   "2026-11-01",
		EndDate:     "2026-11-05",
		Days:        5,
		Description: "Ferie per motivi personali",
		SubmitNow:   true,
	})
	require.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusSubmitted, created.Status)

	// Step 1: Assistente Amministrativo (AA) Istruttoria
	aaID := "aa-1"
	afterAA, err := svc.AAReview(ctx, schoolID, created.ID, aaID, personnel_desk.AAReviewInput{
		Note:    "Documentazione e disponibilità ore verificate regolarmente",
		Approve: true,
	})
	require.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusDSGAReview, afterAA.Status)
	assert.NotNil(t, afterAA.AAReviewedBy)
	assert.Equal(t, aaID, *afterAA.AAReviewedBy)
	assert.NotNil(t, afterAA.AAReviewedAt)
	assert.Equal(t, "Documentazione e disponibilità ore verificate regolarmente", afterAA.AANote)

	// Step 2: Visto DSGA
	dsgaID := "dsga-1"
	afterDSGA, err := svc.DSGASign(ctx, schoolID, created.ID, dsgaID, personnel_desk.DSGASignInput{
		Note:    "Visto contabile e di copertura favorevole",
		Approve: true,
	})
	require.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusDSReview, afterDSGA.Status)
	assert.NotNil(t, afterDSGA.DSGASignedBy)
	assert.Equal(t, dsgaID, *afterDSGA.DSGASignedBy)
	assert.NotNil(t, afterDSGA.DSGASignedAt)

	// Step 3: Approvazione Dirigente Scolastico con emissione decreto
	dsID := "ds-1"
	decreeNum := "DEC-2026/044-FERIE"
	approved, err := svc.DSApprove(ctx, schoolID, created.ID, dsID, personnel_desk.DSApproveInput{
		DecreeNum: decreeNum,
		Note:      "Concesso ai sensi del CCNL vigente",
		Approve:   true,
	})
	require.NoError(t, err)
	assert.Equal(t, personnel_desk.StatusApproved, approved.Status)
	assert.NotNil(t, approved.DSApprovedBy)
	assert.Equal(t, dsID, *approved.DSApprovedBy)
	assert.NotNil(t, approved.DSDecreeNum)
	assert.Equal(t, decreeNum, *approved.DSDecreeNum)
}

func TestPersonnelDesk_Rejection_Branches(t *testing.T) {
	ctx := context.Background()

	t.Run("AA rejection sets status to rejected", func(t *testing.T) {
		repo := newMockDeskRepo()
		svc := personnel_desk.NewService(repo)
		req, err := svc.CreateRequest(ctx, "school-1", "user-1", personnel_desk.CreateDeskRequestInput{
			Category:    personnel_desk.CatPermessoBreve,
			StartDate:   "2026-10-20",
			EndDate:     "2026-10-20",
			Hours:       1,
			Description: "Test",
			SubmitNow:   true,
		})
		require.NoError(t, err)

		rejected, err := svc.AAReview(ctx, "school-1", req.ID, "aa-1", personnel_desk.AAReviewInput{
			Note:    "Mancanza documentazione giustificativa",
			Approve: false,
		})
		require.NoError(t, err)
		assert.Equal(t, personnel_desk.StatusRejected, rejected.Status)
	})

	t.Run("DSGA rejection sets status to rejected", func(t *testing.T) {
		repo := newMockDeskRepo()
		svc := personnel_desk.NewService(repo)
		req, _ := svc.CreateRequest(ctx, "school-1", "user-1", personnel_desk.CreateDeskRequestInput{
			Category:    personnel_desk.CatPermessoStudio,
			StartDate:   "2026-11-01",
			EndDate:     "2026-11-01",
			Hours:       4,
			Description: "Permesso 150 ore studio",
			SubmitNow:   true,
		})
		afterAA, _ := svc.AAReview(ctx, "school-1", req.ID, "aa-1", personnel_desk.AAReviewInput{Approve: true})

		rejected, err := svc.DSGASign(ctx, "school-1", afterAA.ID, "dsga-1", personnel_desk.DSGASignInput{
			Note:    "Superamento contingente massimo ore",
			Approve: false,
		})
		require.NoError(t, err)
		assert.Equal(t, personnel_desk.StatusRejected, rejected.Status)
	})

	t.Run("DS rejection sets status to rejected", func(t *testing.T) {
		repo := newMockDeskRepo()
		svc := personnel_desk.NewService(repo)
		req, _ := svc.CreateRequest(ctx, "school-1", "user-1", personnel_desk.CreateDeskRequestInput{
			Category:    personnel_desk.CatAspettativa,
			StartDate:   "2026-12-01",
			EndDate:     "2026-12-31",
			Days:        31,
			Description: "Aspettativa non retribuita",
			SubmitNow:   true,
		})
		afterAA, _ := svc.AAReview(ctx, "school-1", req.ID, "aa-1", personnel_desk.AAReviewInput{Approve: true})
		afterDSGA, _ := svc.DSGASign(ctx, "school-1", afterAA.ID, "dsga-1", personnel_desk.DSGASignInput{Approve: true})

		rejected, err := svc.DSApprove(ctx, "school-1", afterDSGA.ID, "ds-1", personnel_desk.DSApproveInput{
			Note:    "Motivazioni di servizio ostative",
			Approve: false,
		})
		require.NoError(t, err)
		assert.Equal(t, personnel_desk.StatusRejected, rejected.Status)
	})
}

func TestPersonnelDesk_InvalidStateTransitions_Guards(t *testing.T) {
	ctx := context.Background()
	repo := newMockDeskRepo()
	svc := personnel_desk.NewService(repo)

	req, err := svc.CreateRequest(ctx, "school-1", "user-1", personnel_desk.CreateDeskRequestInput{
		Category:    personnel_desk.CatMalattia,
		StartDate:   "2026-10-10",
		EndDate:     "2026-10-12",
		Days:        3,
		Description: "Certificato medico telematico",
		SubmitNow:   true,
	})
	require.NoError(t, err)

	// Guard 1: DSGA cannot sign while still in submitted (needs AA review first)
	_, err = svc.DSGASign(ctx, "school-1", req.ID, "dsga-1", personnel_desk.DSGASignInput{Approve: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "la richiesta non è in attesa di visto DSGA")

	// Guard 2: DS cannot approve while still in submitted
	_, err = svc.DSApprove(ctx, "school-1", req.ID, "ds-1", personnel_desk.DSApproveInput{Approve: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "la richiesta non è in attesa di approvazione DS")

	// Pass AA review
	afterAA, err := svc.AAReview(ctx, "school-1", req.ID, "aa-1", personnel_desk.AAReviewInput{Approve: true})
	require.NoError(t, err)

	// Guard 3: AA cannot review twice once moved to dsga_review
	_, err = svc.AAReview(ctx, "school-1", afterAA.ID, "aa-2", personnel_desk.AAReviewInput{Approve: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "la richiesta non si trova nello stato di istruttoria")
}

func TestPersonnelDesk_AccessControl_RBAC(t *testing.T) {
	ctx := context.Background()
	repo := newMockDeskRepo()
	svc := personnel_desk.NewService(repo)

	schoolID := "school-1"
	applicantID := "docente-rossi"

	req, err := svc.CreateRequest(ctx, schoolID, applicantID, personnel_desk.CreateDeskRequestInput{
		Category:    personnel_desk.CatFerie,
		StartDate:   "2026-11-10",
		EndDate:     "2026-11-12",
		Days:        3,
		Description: "Ferie",
		SubmitNow:   true,
	})
	require.NoError(t, err)

	// 1. Applicant can view their own request
	viewOwn, err := svc.GetRequest(ctx, schoolID, req.ID, applicantID, "teacher")
	require.NoError(t, err)
	assert.Equal(t, req.ID, viewOwn.ID)

	// 2. Other teacher cannot view applicant's request
	_, err = svc.GetRequest(ctx, schoolID, req.ID, "docente-bianchi", "teacher")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "accesso non autorizzato")

	// 3. Reviewer roles CAN view applicant's request
	reviewerRoles := []string{
		"assistente_amministrativo",
		"dsga",
		"principal",
		"vice_principal",
		"secretary",
		"admin",
		"superadmin",
		"collaboratore_ds",
	}

	for _, role := range reviewerRoles {
		viewByReviewer, err := svc.GetRequest(ctx, schoolID, req.ID, "staff-"+role, role)
		require.NoError(t, err, "Role %s should have access to view request", role)
		assert.Equal(t, req.ID, viewByReviewer.ID)
	}

	// 4. List requests filters by applicant for non-reviewers
	listDocente, err := svc.ListRequests(ctx, schoolID, "docente-bianchi", "teacher", "")
	require.NoError(t, err)
	assert.Len(t, listDocente, 0) // docente-bianchi has 0 requests

	listAA, err := svc.ListRequests(ctx, schoolID, "staff-aa", "assistente_amministrativo", "")
	require.NoError(t, err)
	assert.Len(t, listAA, 1) // AA sees all requests of the school
}
