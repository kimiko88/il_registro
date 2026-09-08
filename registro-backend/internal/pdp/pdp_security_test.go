package pdp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"registro-backend/internal/users"
)

type mockUserRepoForPdp struct {
	guardians map[string]string // parentID:studentID -> exists
}

func (m *mockUserRepoForPdp) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	_, ok := m.guardians[parentID+":"+studentID]
	return ok, nil
}
func (m *mockUserRepoForPdp) ListByIDs(ctx context.Context, ids []string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) GetByID(ctx context.Context, id string) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) Create(ctx context.Context, u *users.User) error { return nil }
func (m *mockUserRepoForPdp) Update(ctx context.Context, u *users.User) error { return nil }
func (m *mockUserRepoForPdp) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) GetPasswordHistory(ctx context.Context, id string) ([]string, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) AddPasswordHistory(ctx context.Context, id, hash string) error {
	return nil
}
func (m *mockUserRepoForPdp) RevokeAllUserTokens(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForPdp) ClearTempMFASecret(ctx context.Context, id string) error  { return nil }
func (m *mockUserRepoForPdp) GetAllActive(ctx context.Context, schoolID string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) SoftDelete(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForPdp) Delete(ctx context.Context, id string) error     { return nil }
func (m *mockUserRepoForPdp) Restore(ctx context.Context, id string) error    { return nil }
func (m *mockUserRepoForPdp) HardDelete(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForPdp) List(ctx context.Context, filter users.UserFilter) ([]users.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForPdp) LogAudit(ctx context.Context, log *users.AuditLog) error { return nil }
func (m *mockUserRepoForPdp) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]users.AuditLog, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForPdp) AddGuardian(ctx context.Context, parentID, studentID, relType string) error {
	return nil
}
func (m *mockUserRepoForPdp) RemoveGuardian(ctx context.Context, studentProfileID, parentProfileID string) error {
	return nil
}
func (m *mockUserRepoForPdp) GetGuardians(ctx context.Context, studentProfileID string) ([]users.GuardianInfo, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) GetChildren(ctx context.Context, parentID string) ([]users.StudentChild, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) GetStudentsByClass(ctx context.Context, classID string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) GetStudentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *mockUserRepoForPdp) GetParentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *mockUserRepoForPdp) GetFascicoloSummary(ctx context.Context, studentID string, isActive bool) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockUserRepoForPdp) IsActive(ctx context.Context, id string) (bool, error) { return true, nil }
func (m *mockUserRepoForPdp) BulkCreate(ctx context.Context, users []users.User) (int, []string, error) {
	return 0, nil, nil
}
func (m *mockUserRepoForPdp) BulkDelete(ctx context.Context, ids []string) (int, error) {
	return 0, nil
}
func (m *mockUserRepoForPdp) ChangePasswordTx(ctx context.Context, userID, newPasswordHash string) error {
	return nil
}
func (m *mockUserRepoForPdp) ApplyDataRetention(ctx context.Context, schoolID *string, cutoffDate time.Time) (int, error) {
	return 0, nil
}

type mockPdpRepoForSecTest struct {
	plans map[string]*PdpPlan
}

func (m *mockPdpRepoForSecTest) Create(ctx context.Context, plan *PdpPlan) (*PdpPlan, error) {
	plan.ID = "pdp-1"
	m.plans[plan.ID] = plan
	return plan, nil
}
func (m *mockPdpRepoForSecTest) GetByStudent(ctx context.Context, studentID, academicYear string) ([]*PdpPlan, error) {
	var res []*PdpPlan
	for _, p := range m.plans {
		if p.StudentID == studentID {
			res = append(res, p)
		}
	}
	return res, nil
}
func (m *mockPdpRepoForSecTest) GetByClass(ctx context.Context, classID, academicYear string) ([]*PdpPlan, error) {
	return nil, nil
}
func (m *mockPdpRepoForSecTest) GetByID(ctx context.Context, id string) (*PdpPlan, error) {
	if p, ok := m.plans[id]; ok {
		return p, nil
	}
	return nil, nil
}
func (m *mockPdpRepoForSecTest) Update(ctx context.Context, id string, req *UpdatePdpRequest) (*PdpPlan, error) {
	return m.plans[id], nil
}
func (m *mockPdpRepoForSecTest) SetSharedWithFamily(ctx context.Context, id string, shared bool) error {
	if p, ok := m.plans[id]; ok {
		p.SharedWithFamily = shared
	}
	return nil
}
func (m *mockPdpRepoForSecTest) ApproveByFamily(ctx context.Context, id, approvedBy string) error {
	if p, ok := m.plans[id]; ok {
		now := time.Now()
		p.FamilyApprovedAt = &now
		p.FamilyApprovedBy = &approvedBy
	}
	return nil
}
func (m *mockPdpRepoForSecTest) Delete(ctx context.Context, id string) error {
	delete(m.plans, id)
	return nil
}

func TestGetByStudent_ParentNotGuardian_Forbidden(t *testing.T) {
	uRepo := &mockUserRepoForPdp{
		guardians: map[string]string{
			"parent-1:student-1": "ok",
		},
	}
	pRepo := &mockPdpRepoForSecTest{
		plans: map[string]*PdpPlan{
			"pdp-1": {
				ID: "pdp-1", StudentID: "student-2", SchoolID: "school-1",
				SharedWithFamily: true,
			},
		},
	}
	svc := NewService(pRepo, uRepo)

	// Parent-1 attempts to view PDP of student-2 (not their child)
	_, err := svc.GetByStudent(context.Background(), "parent-1", "parent", "school-1", "student-2", "")
	assert.ErrorIs(t, err, ErrNotGuardian, "Parent not linked as guardian must be rejected")
}

func TestGetByStudent_ParentGuardian_Allowed(t *testing.T) {
	uRepo := &mockUserRepoForPdp{
		guardians: map[string]string{
			"parent-1:student-1": "ok",
		},
	}
	pRepo := &mockPdpRepoForSecTest{
		plans: map[string]*PdpPlan{
			"pdp-1": {
				ID: "pdp-1", StudentID: "student-1", SchoolID: "school-1",
				SharedWithFamily: true, Diagnosis: "Secret F81.0",
			},
		},
	}
	svc := NewService(pRepo, uRepo)

	plans, err := svc.GetByStudent(context.Background(), "parent-1", "parent", "school-1", "student-1", "")
	assert.NoError(t, err)
	assert.Len(t, plans, 1)
	assert.Empty(t, plans[0].Diagnosis, "Diagnosis must be stripped for parent")
}

func TestApproveByFamily_NonGuardian_Forbidden(t *testing.T) {
	uRepo := &mockUserRepoForPdp{
		guardians: map[string]string{
			"parent-1:student-1": "ok",
		},
	}
	pRepo := &mockPdpRepoForSecTest{
		plans: map[string]*PdpPlan{
			"pdp-2": {
				ID: "pdp-2", StudentID: "student-2", SchoolID: "school-1",
				SharedWithFamily: true,
			},
		},
	}
	svc := NewService(pRepo, uRepo)

	// Parent-1 attempts to approve PDP of student-2
	err := svc.ApproveByFamily(context.Background(), "parent", "parent-1", "school-1", "pdp-2")
	assert.ErrorIs(t, err, ErrNotGuardian)
}
