package personnel_desk_test

import (
	"context"
	"testing"

	"registro-backend/internal/personnel_desk"
)

type mockDeskRepo struct {
	items map[string]*personnel_desk.DeskRequest
}

func newMockDeskRepo() *mockDeskRepo {
	return &mockDeskRepo{
		items: make(map[string]*personnel_desk.DeskRequest),
	}
}

func (m *mockDeskRepo) Create(ctx context.Context, r *personnel_desk.DeskRequest) error {
	r.ID = "req-1"
	m.items[r.ID] = r
	return nil
}

func (m *mockDeskRepo) GetByID(ctx context.Context, schoolID, id string) (*personnel_desk.DeskRequest, error) {
	if item, ok := m.items[id]; ok {
		return item, nil
	}
	return nil, context.DeadlineExceeded
}

func (m *mockDeskRepo) List(ctx context.Context, schoolID, applicantID, status string) ([]*personnel_desk.DeskRequest, error) {
	var res []*personnel_desk.DeskRequest
	for _, it := range m.items {
		if applicantID != "" && it.ApplicantID != applicantID {
			continue
		}
		if status != "" && string(it.Status) != status {
			continue
		}
		res = append(res, it)
	}
	return res, nil
}

func (m *mockDeskRepo) Update(ctx context.Context, r *personnel_desk.DeskRequest) error {
	m.items[r.ID] = r
	return nil
}

func (m *mockDeskRepo) Delete(ctx context.Context, schoolID, id, applicantID string) error {
	delete(m.items, id)
	return nil
}

func TestPersonnelDeskWorkflow(t *testing.T) {
	repo := newMockDeskRepo()
	svc := personnel_desk.NewService(repo)
	ctx := context.Background()

	// 1. Employee creates request
	input := personnel_desk.CreateDeskRequestInput{
		Category:    personnel_desk.CatFerie,
		StartDate:   "2026-10-01",
		EndDate:     "2026-10-05",
		Days:        5,
		Description: "Ferie autunnali",
		SubmitNow:   true,
	}
	req, err := svc.CreateRequest(ctx, "sch-1", "user-emp-1", input)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}
	if req.Status != personnel_desk.StatusSubmitted {
		t.Errorf("expected status submitted, got %s", req.Status)
	}

	// 2. AA review
	aaInput := personnel_desk.AAReviewInput{
		Note:    "Documentazione verificata e completa",
		Approve: true,
	}
	req, err = svc.AAReview(ctx, "sch-1", req.ID, "user-aa-1", aaInput)
	if err != nil {
		t.Fatalf("unexpected error in AA review: %v", err)
	}
	if req.Status != personnel_desk.StatusDSGAReview {
		t.Errorf("expected status dsga_review, got %s", req.Status)
	}

	// 3. DSGA sign
	dsgaInput := personnel_desk.DSGASignInput{
		Note:    "Visto favorevole, compatibile con monte ore",
		Approve: true,
	}
	req, err = svc.DSGASign(ctx, "sch-1", req.ID, "user-dsga-1", dsgaInput)
	if err != nil {
		t.Fatalf("unexpected error in DSGA sign: %v", err)
	}
	if req.Status != personnel_desk.StatusDSReview {
		t.Errorf("expected status ds_review, got %s", req.Status)
	}

	// 4. DS approval
	dsInput := personnel_desk.DSApproveInput{
		DecreeNum: "DEC-2026/42",
		Note:      "Approvato",
		Approve:   true,
	}
	req, err = svc.DSApprove(ctx, "sch-1", req.ID, "user-ds-1", dsInput)
	if err != nil {
		t.Fatalf("unexpected error in DS approve: %v", err)
	}
	if req.Status != personnel_desk.StatusApproved {
		t.Errorf("expected status approved, got %s", req.Status)
	}
	if req.DSDecreeNum == nil || *req.DSDecreeNum != "DEC-2026/42" {
		t.Errorf("expected decree DEC-2026/42, got %v", req.DSDecreeNum)
	}
}
