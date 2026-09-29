package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"registro-backend/internal/substitutions"
)

// mockSubstitutionsRepo implements substitutions.Repository for unit tests
type mockSubstitutionsRepo struct {
	substitutions.Repository
	subs             map[string]*substitutions.Substitution
	available        []substitutions.TeacherCandidate
	classTeachers    map[string]bool
	subjectTeachers  map[string]bool
	subCounts        map[string]int
	profileIDMapping map[string]string
}

func newMockSubstitutionsRepo() *mockSubstitutionsRepo {
	return &mockSubstitutionsRepo{
		subs:             make(map[string]*substitutions.Substitution),
		classTeachers:    make(map[string]bool),
		subjectTeachers:  make(map[string]bool),
		subCounts:        make(map[string]int),
		profileIDMapping: make(map[string]string),
	}
}

func (m *mockSubstitutionsRepo) Create(ctx context.Context, sub *substitutions.Substitution) error {
	if sub.ID == "" {
		sub.ID = fmt.Sprintf("sub-%d", len(m.subs)+1)
	}
	sub.CreatedAt = time.Now()
	m.subs[sub.ID] = sub
	return nil
}

func (m *mockSubstitutionsRepo) GetByID(ctx context.Context, id string) (*substitutions.Substitution, error) {
	sub, ok := m.subs[id]
	if !ok {
		return nil, fmt.Errorf("substitution not found")
	}
	return sub, nil
}

func (m *mockSubstitutionsRepo) AssignSubstitute(ctx context.Context, id, substituteTeacherID, notes string) error {
	sub, ok := m.subs[id]
	if !ok {
		return fmt.Errorf("substitution not found")
	}
	sub.SubstituteTeacherID = &substituteTeacherID
	sub.Status = substitutions.StatusAssigned
	sub.Notes = notes
	return nil
}

func (m *mockSubstitutionsRepo) ConfirmSubstitution(ctx context.Context, id, teacherID string) error {
	sub, ok := m.subs[id]
	if !ok {
		return fmt.Errorf("substitution not found")
	}
	sub.Status = substitutions.StatusConfirmed
	return nil
}

func (m *mockSubstitutionsRepo) SignRegister(ctx context.Context, id, sigHash, notes string) error {
	sub, ok := m.subs[id]
	if !ok {
		return fmt.Errorf("substitution not found")
	}
	now := time.Now()
	sub.SignedBySubstitute = true
	sub.SignatureHash = sigHash
	sub.SignatureTimestamp = &now
	sub.OfficialRegisterNotes = notes
	return nil
}

func (m *mockSubstitutionsRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	if p, ok := m.profileIDMapping[userID]; ok {
		return p, nil
	}
	return userID, nil
}

func (m *mockSubstitutionsRepo) GetAvailableTeachers(ctx context.Context, schoolID string) ([]substitutions.TeacherCandidate, error) {
	return m.available, nil
}

func (m *mockSubstitutionsRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return m.classTeachers[teacherID+":"+classID], nil
}

func (m *mockSubstitutionsRepo) IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error) {
	return m.subjectTeachers[teacherID+":"+subjectID], nil
}

func (m *mockSubstitutionsRepo) GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error) {
	return m.subCounts[teacherID], nil
}

func TestEmergencySubstitutions_RecommendationScoringAndSorting(t *testing.T) {
	repo := newMockSubstitutionsRepo()
	repo.available = []substitutions.TeacherCandidate{
		{TeacherID: "t-low", TeacherName: "Docente Low"},
		{TeacherID: "t-perfect", TeacherName: "Docente Perfect"},
		{TeacherID: "t-medium", TeacherName: "Docente Medium"},
	}

	// t-perfect: teaches class (+25), teaches subject (+15), 0 weekly subs (+10) => 50+25+15+10 = 100 capped at 95
	repo.classTeachers["t-perfect:class-1"] = true
	repo.subjectTeachers["t-perfect:subj-math"] = true
	repo.subCounts["t-perfect"] = 0

	// t-medium: teaches subject (+15), 1 weekly sub (+5) => 50+15+5 = 70
	repo.subjectTeachers["t-medium:subj-math"] = true
	repo.subCounts["t-medium"] = 1

	// t-low: not teaching class or subject, 5 weekly subs (0 bonus) => 50
	repo.subCounts["t-low"] = 5

	svc := substitutions.NewService(repo)
	recs, err := svc.RecommendSubstitutes(context.Background(), "school-1", "class-1", "subj-math", "2026-10-15", 2)
	require.NoError(t, err)
	require.Len(t, recs, 3)

	// First recommendation must be the highest score (t-perfect with 95)
	assert.Equal(t, "t-perfect", recs[0].TeacherID)
	assert.Equal(t, 95, recs[0].Score)
	assert.True(t, recs[0].TeachesClass)
	assert.True(t, recs[0].TeachesSubject)
	assert.Contains(t, recs[0].Reason, "Docente della stessa classe")

	// Second must be t-medium with 70
	assert.Equal(t, "t-medium", recs[1].TeacherID)
	assert.Equal(t, 70, recs[1].Score)
	assert.False(t, recs[1].TeachesClass)
	assert.True(t, recs[1].TeachesSubject)

	// Third must be t-low with 50
	assert.Equal(t, "t-low", recs[2].TeacherID)
	assert.Equal(t, 50, recs[2].Score)
}

func TestEmergencySubstitutions_RoleAuthorization(t *testing.T) {
	repo := newMockSubstitutionsRepo()
	svc := substitutions.NewService(repo)

	req := substitutions.CreateSubstitutionRequest{
		ClassID:         "class-1",
		AbsentTeacherID: "t-absent",
		Date:            "2026-10-15",
		Slot:            2,
		Subject:         "Matematica",
	}

	// Unauthorized roles: student, parent, external
	for _, role := range []string{"student", "parent", "external"} {
		_, err := svc.CreateSubstitution(context.Background(), role, "school-1", req)
		assert.Error(t, err, "role %s should be rejected", role)
		assert.Contains(t, err.Error(), "insufficient permissions")
	}

	// Authorized roles: admin, principal, vice_principal, secretary
	for _, role := range []string{"admin", "principal", "vice_principal", "secretary"} {
		sub, err := svc.CreateSubstitution(context.Background(), role, "school-1", req)
		assert.NoError(t, err, "role %s should be authorized", role)
		assert.NotNil(t, sub)
		assert.Equal(t, substitutions.StatusPending, sub.Status)
	}
}

func TestEmergencySubstitutions_SignRegisterLifecycle(t *testing.T) {
	repo := newMockSubstitutionsRepo()
	svc := substitutions.NewService(repo)

	// Create pending substitution
	sub, err := svc.CreateSubstitution(context.Background(), "admin", "school-1", substitutions.CreateSubstitutionRequest{
		ClassID:         "class-10",
		AbsentTeacherID: "t-absent",
		Date:            "2026-10-15",
		Slot:            1,
		Subject:         "Italiano",
	})
	require.NoError(t, err)

	// Assign substitute
	err = svc.AssignSubstitute(context.Background(), sub.ID, "vice_principal", substitutions.AssignSubstituteRequest{
		SubstituteTeacherID: "t-substitute-1",
		Notes:               "Supplenza prima ora urgente",
	})
	require.NoError(t, err)

	// Another teacher tries to sign the register -> must be rejected
	err = svc.SignRegister(context.Background(), sub.ID, "t-impostor", "Appunti lezione")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "solo il docente sostituto assegnato può firmare il registro")

	// Authorized substitute teacher signs the register
	err = svc.SignRegister(context.Background(), sub.ID, "t-substitute-1", "Lezione svolta regolarmente su I Promessi Sposi cap. 3")
	require.NoError(t, err)

	// Check register signature on stored substitution
	updatedSub, err := repo.GetByID(context.Background(), sub.ID)
	require.NoError(t, err)
	assert.True(t, updatedSub.SignedBySubstitute)
	assert.NotEmpty(t, updatedSub.SignatureHash)
	assert.Equal(t, "Lezione svolta regolarmente su I Promessi Sposi cap. 3", updatedSub.OfficialRegisterNotes)
	assert.NotNil(t, updatedSub.SignatureTimestamp)
}
