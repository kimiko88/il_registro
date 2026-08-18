package attendance

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Minimal in-memory mocks
// ---------------------------------------------------------------------------

// mockRepo satisfies attendance.Repository with configurable stubs.
type mockRepo struct {
	isAssigned        bool
	isAssignedErr     error
	isSub             bool
	isSubErr          error
	justification     *Justification
	justificationErr  error
	deleteErr         error
	processTxErr      error
	hasOverlap        bool
	hasOverlapErr     error
	createErr         error
	batchCreateErr    error
	updateErr         error
	findByIDResult    *Attendance
	findByIDErr       error
	findByClassResult []*Attendance
	findByStudResult  []Attendance
}

func (m *mockRepo) IsTeacherAssignedToClass(_ context.Context, _, _ string) (bool, error) {
	return m.isAssigned, m.isAssignedErr
}
func (m *mockRepo) AreStudentsInClass(_ context.Context, studentIDs []string, _ string) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, id := range studentIDs {
		res[id] = true
	}
	return res, nil
}
func (m *mockRepo) IsTeacherSubstitute(_ context.Context, _, _ string, _ time.Time, _ int) (bool, error) {
	return m.isSub, m.isSubErr
}
func (m *mockRepo) Create(_ *Attendance) error        { return m.createErr }
func (m *mockRepo) BatchCreate(_ []*Attendance) error { return m.batchCreateErr }
func (m *mockRepo) Update(_ *Attendance) error        { return m.updateErr }
func (m *mockRepo) FindByID(id string) (*Attendance, error) {
	if m.findByIDResult != nil {
		return m.findByIDResult, m.findByIDErr
	}
	schoolID := "school1"
	hour := 1
	classID := "class1"
	return &Attendance{
		ID:       id,
		SchoolID: schoolID,
		ClassID:  classID,
		Hour:     &hour,
		Date:     time.Now().Add(-24 * time.Hour),
		Status:   StatusPresent,
	}, m.findByIDErr
}
func (m *mockRepo) FindByClassAndDate(_ string, _ time.Time) ([]Attendance, error) {
	var res []Attendance
	for _, r := range m.findByClassResult {
		if r != nil {
			res = append(res, *r)
		}
	}
	return res, nil
}
func (m *mockRepo) FindByStudent(_ string, _, _ time.Time) ([]Attendance, error) {
	return m.findByStudResult, nil
}
func (m *mockRepo) FindUnjustifiedByStudent(_ string) ([]Attendance, error) {
	return nil, nil
}
func (m *mockRepo) JustifyAbsenceByParent(_, _, _, _ string) error { return nil }
func (m *mockRepo) GetStats(_ string) (*SummaryResponse, error) {
	return &SummaryResponse{}, nil
}
func (m *mockRepo) GetStatsBatch(_ context.Context, _ []string) (map[string]*SummaryResponse, error) {
	return map[string]*SummaryResponse{}, nil
}
func (m *mockRepo) CountDistinctDays(_ string) (int, error) {
	return 10, nil
}
func (m *mockRepo) GetAnalytics(_ context.Context, _ string) (*AnalyticsResponse, error) {
	return &AnalyticsResponse{}, nil
}
func (m *mockRepo) DeleteByClassDateHour(_, _ string, _ time.Time, _ int) error { return m.deleteErr }
func (m *mockRepo) FindJustificationByID(_ string) (*Justification, error) {
	if m.justification != nil {
		return m.justification, m.justificationErr
	}
	return nil, m.justificationErr
}
func (m *mockRepo) CreateJustification(_ *Justification) error { return m.createErr }
func (m *mockRepo) UpdateJustification(_ *Justification) error { return nil }
func (m *mockRepo) HasOverlappingJustification(_ context.Context, _ string, _, _ time.Time) (bool, error) {
	return m.hasOverlap, m.hasOverlapErr
}
func (m *mockRepo) FindPendingJustifications(_, _ string) ([]Justification, error) {
	return nil, nil
}
func (m *mockRepo) DeleteJustification(_ string) error        { return m.deleteErr }
func (m *mockRepo) DeletePendingJustification(_ string) error { return m.deleteErr }
func (m *mockRepo) IsStudentInClass(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}
func (m *mockRepo) IsClassInSchool(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}
func (m *mockRepo) ProcessJustificationTx(_ context.Context, _ *Justification, _ string, _ bool) error {
	return m.processTxErr
}
func (m *mockRepo) GetMonthlyBreakdown(_ context.Context, _ string, _ string) ([]MonthlyBreakdownRow, error) {
	return []MonthlyBreakdownRow{}, nil
}
func (m *mockRepo) GetStudentAttendanceStats(_ string) (*AttendanceStats, error) {
	return &AttendanceStats{}, nil
}

// mockUserRepo satisfies users.Repository with configurable stubs.
type mockUserRepo struct {
	isGuardian    bool
	isGuardianErr error
	users         map[string]*users.User
}

func (m *mockUserRepo) IsGuardian(_ context.Context, _, _ string) (bool, error) {
	return m.isGuardian, m.isGuardianErr
}
func (m *mockUserRepo) GetByID(_ context.Context, id string) (*users.User, error) {
	if m.users != nil {
		if u, ok := m.users[id]; ok {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found: %s", id)
}
func (m *mockUserRepo) Create(_ context.Context, _ *users.User) error               { return nil }
func (m *mockUserRepo) GetByEmail(_ context.Context, _ string) (*users.User, error) { return nil, nil }
func (m *mockUserRepo) Update(_ context.Context, _ *users.User) error               { return nil }
func (m *mockUserRepo) GetPasswordHistory(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}
func (m *mockUserRepo) AddPasswordHistory(_ context.Context, _, _ string) error { return nil }
func (m *mockUserRepo) Delete(_ context.Context, _ string) error                { return nil }
func (m *mockUserRepo) Restore(_ context.Context, _ string) error               { return nil }
func (m *mockUserRepo) List(_ context.Context, _ users.UserFilter) ([]users.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) ListByIDs(_ context.Context, _ []string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepo) LogAudit(_ context.Context, _ *users.AuditLog) error { return nil }
func (m *mockUserRepo) GetAuditLogs(_ context.Context, _ string, _, _ int) ([]users.AuditLog, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) BulkCreate(_ context.Context, _ []users.User) (int, []string, error) {
	return 0, nil, nil
}
func (m *mockUserRepo) BulkDelete(_ context.Context, _ []string) (int, error) { return 0, nil }
func (m *mockUserRepo) HardDelete(_ context.Context, _ string) error          { return nil }
func (m *mockUserRepo) RevokeAllUserTokens(_ context.Context, _ string) error { return nil }
func (m *mockUserRepo) ClearTempMFASecret(_ context.Context, _ string) error  { return nil }
func (m *mockUserRepo) AddGuardian(_ context.Context, _, _, _ string) error   { return nil }
func (m *mockUserRepo) GetChildren(_ context.Context, _ string) ([]users.StudentChild, error) {
	return nil, nil
}
func (m *mockUserRepo) GetStudentsByClass(_ context.Context, _ string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetStudentProfile(_ context.Context, _ string) (string, error) { return "", nil }
func (m *mockUserRepo) GetParentProfile(_ context.Context, _ string) (string, error)  { return "", nil }
func (m *mockUserRepo) RemoveGuardian(_ context.Context, _, _ string) error           { return nil }
func (m *mockUserRepo) GetGuardians(_ context.Context, _ string) ([]users.GuardianInfo, error) {
	return nil, nil
}
func (m *mockUserRepo) GetFascicoloSummary(_ context.Context, _ string, _ bool) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockUserRepo) IsActive(_ context.Context, _ string) (bool, error) { return true, nil }

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func strPtr(s string) *string { return &s }

func makeService(r *mockRepo, u *mockUserRepo) Service {
	return NewService(r, u, nil, nil)
}

// ---------------------------------------------------------------------------
// MarkAttendance
// ---------------------------------------------------------------------------

func TestMarkAttendance_MissingTeacherID_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: true}, &mockUserRepo{})
	err := svc.MarkAttendance(context.Background(), "", "school1", CreateAttendanceRequest{
		ClassID:   "class1",
		StudentID: "stu1",
		Date:      "2024-01-10",
		Status:    StatusPresent,
		Hour:      1,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestMarkAttendance_TeacherNotAssigned_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: false}, &mockUserRepo{})
	err := svc.MarkAttendance(context.Background(), "teacher1", "school1", CreateAttendanceRequest{
		ClassID:   "class1",
		StudentID: "stu1",
		Date:      "2024-01-10",
		Status:    StatusPresent,
		Hour:      1,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestMarkAttendance_AssignedTeacher_Passes(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: true, createErr: nil}, &mockUserRepo{})
	err := svc.MarkAttendance(context.Background(), "teacher1", "school1", CreateAttendanceRequest{
		ClassID:   "class1",
		StudentID: "stu1",
		Date:      time.Now().Format("2006-01-02"),
		Status:    StatusPresent,
		Hour:      1,
	})
	assert.NoError(t, err)
}

// ---------------------------------------------------------------------------
// MarkBulk
// ---------------------------------------------------------------------------

func TestMarkBulk_MissingTeacherID_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: true}, &mockUserRepo{})
	err := svc.MarkBulk(context.Background(), "", "school1", BulkAttendanceRequest{
		ClassID: "class1",
		Date:    time.Now().Format("2006-01-02"),
		Hour:    1,
		Statuses: []StudentStatusRequest{
			{StudentID: "stu1", Status: StatusPresent},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestMarkBulk_TeacherNotAssigned_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: false}, &mockUserRepo{})
	err := svc.MarkBulk(context.Background(), "teacher1", "school1", BulkAttendanceRequest{
		ClassID: "class1",
		Date:    time.Now().Format("2006-01-02"),
		Hour:    1,
		Statuses: []StudentStatusRequest{
			{StudentID: "stu1", Status: StatusPresent},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

// FIX regression: malformed date in substitution path must return error (not zero-time bypass).
func TestMarkBulk_Substitution_MalformedDate_ReturnsError(t *testing.T) {
	svc := makeService(&mockRepo{isSub: false, isAssigned: false}, &mockUserRepo{})
	err := svc.MarkBulk(context.Background(), "teacher1", "school1", BulkAttendanceRequest{
		ClassID:        "class1",
		Date:           "NOT-A-DATE",
		Hour:           1,
		IsSubstitution: true,
		Statuses: []StudentStatusRequest{
			{StudentID: "stu1", Status: StatusPresent},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data non valida")
}

// FIX regression: DB error on IsTeacherAssignedToClass must propagate, not grant access.
func TestMarkBulk_Substitution_DBErrorOnAssignment_Propagates(t *testing.T) {
	repo := &mockRepo{
		isSub:         false,
		isSubErr:      errors.New("db timeout"),
		isAssigned:    false,
		isAssignedErr: errors.New("db timeout"),
	}
	svc := makeService(repo, &mockUserRepo{})
	err := svc.MarkBulk(context.Background(), "teacher1", "school1", BulkAttendanceRequest{
		ClassID:        "class1",
		Date:           time.Now().Format("2006-01-02"),
		Hour:           1,
		IsSubstitution: true,
		Statuses: []StudentStatusRequest{
			{StudentID: "stu1", Status: StatusPresent},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "errore verifica docente")
}

// ---------------------------------------------------------------------------
// UpdateAttendance
// ---------------------------------------------------------------------------

func TestUpdateAttendance_SchoolMismatch_Forbidden(t *testing.T) {
	hour := 1
	repo := &mockRepo{
		findByIDResult: &Attendance{
			ID:       "att1",
			SchoolID: "other-school",
			ClassID:  "class1",
			Hour:     &hour,
			Date:     time.Now().Add(-24 * time.Hour),
			Status:   StatusPresent,
		},
		isAssigned: true,
	}
	svc := makeService(repo, &mockUserRepo{})
	statusAbsent := StatusAbsent
	err := svc.UpdateAttendance(context.Background(), "teacher1", "school1", "att1",
		UpdateAttendanceRequest{Status: &statusAbsent})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestUpdateAttendance_TeacherNotAssigned_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: false}, &mockUserRepo{})
	statusAbsent := StatusAbsent
	err := svc.UpdateAttendance(context.Background(), "teacher1", "school1", "att1",
		UpdateAttendanceRequest{Status: &statusAbsent})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestUpdateAttendance_MissingTeacherID_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: true}, &mockUserRepo{})
	statusLate := StatusLate
	err := svc.UpdateAttendance(context.Background(), "", "school1", "att1",
		UpdateAttendanceRequest{Status: &statusLate})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

// ---------------------------------------------------------------------------
// DeleteClassAttendanceHour — role matrix
// ---------------------------------------------------------------------------

func TestDeleteClassAttendanceHour_AuthorisedRoles_Pass(t *testing.T) {
	authorisedRoles := []string{"admin", "superadmin", "secretary", "principal", "vice_principal"}
	for _, role := range authorisedRoles {
		t.Run(role, func(t *testing.T) {
			svc := makeService(&mockRepo{deleteErr: nil}, &mockUserRepo{})
			err := svc.DeleteClassAttendanceHour(context.Background(),
				"actor1", role, "school1", "class1", "2024-01-10", 1)
			assert.NoError(t, err, "role %s should be authorised", role)
		})
	}
}

func TestDeleteClassAttendanceHour_AssignedTeacher_Pass(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: true, deleteErr: nil}, &mockUserRepo{})
	err := svc.DeleteClassAttendanceHour(context.Background(),
		"teacher1", "teacher", "school1", "class1", "2024-01-10", 1)
	assert.NoError(t, err)
}

func TestDeleteClassAttendanceHour_TeacherNotAssigned_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: false}, &mockUserRepo{})
	err := svc.DeleteClassAttendanceHour(context.Background(),
		"teacher1", "teacher", "school1", "class1", "2024-01-10", 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestDeleteClassAttendanceHour_UnauthorisedRoles_Forbidden(t *testing.T) {
	unauthorisedRoles := []string{"student", "parent", "unknown", ""}
	for _, role := range unauthorisedRoles {
		t.Run("role="+role, func(t *testing.T) {
			svc := makeService(&mockRepo{}, &mockUserRepo{})
			err := svc.DeleteClassAttendanceHour(context.Background(),
				"actor1", role, "school1", "class1", "2024-01-10", 1)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "forbidden")
		})
	}
}

// ---------------------------------------------------------------------------
// GetClassAttendance — role matrix
// ---------------------------------------------------------------------------

func TestGetClassAttendance_AuthorisedRoles_Pass(t *testing.T) {
	authorisedRoles := []string{"admin", "superadmin", "secretary", "principal", "vice_principal"}
	for _, role := range authorisedRoles {
		t.Run(role, func(t *testing.T) {
			svc := makeService(&mockRepo{}, &mockUserRepo{})
			_, err := svc.GetClassAttendance(context.Background(),
				"actor1", role, "school1", "class1", "2024-01-10")
			assert.NoError(t, err)
		})
	}
}

func TestGetClassAttendance_AssignedTeacher_Pass(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: true}, &mockUserRepo{})
	_, err := svc.GetClassAttendance(context.Background(),
		"teacher1", "teacher", "school1", "class1", "2024-01-10")
	assert.NoError(t, err)
}

func TestGetClassAttendance_TeacherNotAssigned_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: false}, &mockUserRepo{})
	_, err := svc.GetClassAttendance(context.Background(),
		"teacher1", "teacher", "school1", "class1", "2024-01-10")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestGetClassAttendance_UnauthorisedRoles_Forbidden(t *testing.T) {
	unauthorisedRoles := []string{"student", "parent", "unknown"}
	for _, role := range unauthorisedRoles {
		t.Run(role, func(t *testing.T) {
			svc := makeService(&mockRepo{}, &mockUserRepo{})
			_, err := svc.GetClassAttendance(context.Background(),
				"actor1", role, "school1", "class1", "2024-01-10")
			require.Error(t, err)
			assert.Contains(t, err.Error(), "forbidden")
		})
	}
}

func TestGetClassAttendance_MissingRole_Unauthorized(t *testing.T) {
	svc := makeService(&mockRepo{}, &mockUserRepo{})
	_, err := svc.GetClassAttendance(context.Background(),
		"actor1", "", "school1", "class1", "2024-01-10")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

// ---------------------------------------------------------------------------
// GetStudentAttendance — identity & role checks
// ---------------------------------------------------------------------------

func TestGetStudentAttendance_StudentOwnRecord_Pass(t *testing.T) {
	svc := makeService(&mockRepo{}, &mockUserRepo{})
	_, err := svc.GetStudentAttendance(context.Background(),
		"stu1", "student", "school1", "stu1", time.Now().Add(-7*24*time.Hour), time.Now())
	assert.NoError(t, err)
}

func TestGetStudentAttendance_StudentOtherRecord_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{}, &mockUserRepo{})
	_, err := svc.GetStudentAttendance(context.Background(),
		"stu1", "student", "school1", "stu2", time.Now().Add(-7*24*time.Hour), time.Now())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

// FIX regression: teacher from a different school must be rejected.
func TestGetStudentAttendance_TeacherWrongSchool_Forbidden(t *testing.T) {
	schoolID := "other-school"
	userRepo := &mockUserRepo{
		users: map[string]*users.User{
			"teacher1": {ID: "teacher1", Role: "teacher", SchoolID: &schoolID},
		},
	}
	svc := makeService(&mockRepo{}, userRepo)
	_, err := svc.GetStudentAttendance(context.Background(),
		"teacher1", "teacher", "school1", "stu1",
		time.Now().Add(-7*24*time.Hour), time.Now())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestGetStudentAttendance_TeacherSameSchool_Pass(t *testing.T) {
	schoolID := "school1"
	userRepo := &mockUserRepo{
		users: map[string]*users.User{
			"teacher1": {ID: "teacher1", Role: "teacher", SchoolID: &schoolID},
			"stu1":     {ID: "stu1", Role: "student", SchoolID: &schoolID},
		},
	}
	svc := makeService(&mockRepo{}, userRepo)
	_, err := svc.GetStudentAttendance(context.Background(),
		"teacher1", "teacher", "school1", "stu1",
		time.Now().Add(-7*24*time.Hour), time.Now())
	assert.NoError(t, err)
}

func TestGetStudentAttendance_Guardian_Pass(t *testing.T) {
	userRepo := &mockUserRepo{isGuardian: true}
	svc := makeService(&mockRepo{}, userRepo)
	_, err := svc.GetStudentAttendance(context.Background(),
		"parent1", "parent", "school1", "stu1",
		time.Now().Add(-7*24*time.Hour), time.Now())
	assert.NoError(t, err)
}

func TestGetStudentAttendance_NonGuardianParent_Forbidden(t *testing.T) {
	userRepo := &mockUserRepo{isGuardian: false}
	svc := makeService(&mockRepo{}, userRepo)
	_, err := svc.GetStudentAttendance(context.Background(),
		"parent1", "parent", "school1", "stu1",
		time.Now().Add(-7*24*time.Hour), time.Now())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

// ---------------------------------------------------------------------------
// ProcessJustification — teacher role check
// ---------------------------------------------------------------------------

func TestProcessJustification_TeacherNotAssigned_Forbidden(t *testing.T) {
	classID := "class1"
	repo := &mockRepo{
		justification: &Justification{
			ID:        "just1",
			StudentID: "stu1",
			ParentID:  "parent1",
			Status:    JustificationPending,
			StartDate: time.Now().Add(-48 * time.Hour),
			EndDate:   time.Now().Add(-24 * time.Hour),
		},
		isAssigned: false,
	}
	schoolID := "school1"
	userRepo := &mockUserRepo{
		users: map[string]*users.User{
			"stu1":     {ID: "stu1", Role: "student", ClassID: &classID, SchoolID: &schoolID},
			"teacher1": {ID: "teacher1", Role: "teacher", SchoolID: &schoolID},
		},
	}
	svc := makeService(repo, userRepo)
	err := svc.ProcessJustification(context.Background(), "teacher1", "just1", true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestProcessJustification_AssignedTeacher_Pass(t *testing.T) {
	classID := "class1"
	repo := &mockRepo{
		justification: &Justification{
			ID:        "just1",
			StudentID: "stu1",
			ParentID:  "parent1",
			Status:    JustificationPending,
			StartDate: time.Now().Add(-48 * time.Hour),
			EndDate:   time.Now().Add(-24 * time.Hour),
		},
		isAssigned: true,
	}
	schoolID := "school1"
	userRepo := &mockUserRepo{
		users: map[string]*users.User{
			"stu1":     {ID: "stu1", Role: "student", ClassID: &classID, SchoolID: &schoolID},
			"teacher1": {ID: "teacher1", Role: "teacher", SchoolID: &schoolID},
		},
	}
	svc := makeService(repo, userRepo)
	err := svc.ProcessJustification(context.Background(), "teacher1", "just1", true)
	assert.NoError(t, err)
}

// ---------------------------------------------------------------------------
// DeleteJustification — role & ownership checks
// ---------------------------------------------------------------------------

func TestDeleteJustification_TeacherCannotDeleteParentJustification(t *testing.T) {
	schoolID := "school1"
	repo := &mockRepo{
		justification: &Justification{
			ID:        "just1",
			StudentID: "stu1",
			ParentID:  "parent1",
			Status:    JustificationPending,
			StartDate: time.Now().Add(-48 * time.Hour),
			EndDate:   time.Now().Add(-24 * time.Hour),
		},
	}
	userRepo := &mockUserRepo{
		isGuardian: false,
		users: map[string]*users.User{
			"teacher1": {ID: "teacher1", Role: "teacher", SchoolID: &schoolID},
			"stu1":     {ID: "stu1", Role: "student", SchoolID: &schoolID},
		},
	}
	svc := makeService(repo, userRepo)
	err := svc.DeleteJustification(context.Background(), "teacher1", "just1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestDeleteJustification_AdminSameSchool_Pass(t *testing.T) {
	schoolID := "school1"
	repo := &mockRepo{
		justification: &Justification{
			ID:        "just1",
			StudentID: "stu1",
			ParentID:  "parent1",
			Status:    JustificationPending,
			StartDate: time.Now().Add(-48 * time.Hour),
			EndDate:   time.Now().Add(-24 * time.Hour),
		},
		deleteErr: nil,
	}
	userRepo := &mockUserRepo{
		isGuardian: false,
		users: map[string]*users.User{
			"admin1": {ID: "admin1", Role: "admin", SchoolID: &schoolID},
			"stu1":   {ID: "stu1", Role: "student", SchoolID: &schoolID},
		},
	}
	svc := makeService(repo, userRepo)
	err := svc.DeleteJustification(context.Background(), "admin1", "just1")
	assert.NoError(t, err)
}

func TestDeleteJustification_AdminDifferentSchool_Forbidden(t *testing.T) {
	adminSchool := "school-admin"
	stuSchool := "school-stu"
	repo := &mockRepo{
		justification: &Justification{
			ID:        "just1",
			StudentID: "stu1",
			ParentID:  "parent1",
			Status:    JustificationPending,
			StartDate: time.Now().Add(-48 * time.Hour),
			EndDate:   time.Now().Add(-24 * time.Hour),
		},
	}
	userRepo := &mockUserRepo{
		isGuardian: false,
		users: map[string]*users.User{
			"admin1": {ID: "admin1", Role: "admin", SchoolID: &adminSchool},
			"stu1":   {ID: "stu1", Role: "student", SchoolID: &stuSchool},
		},
	}
	svc := makeService(repo, userRepo)
	err := svc.DeleteJustification(context.Background(), "admin1", "just1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestDeleteJustification_ParentOwner_Pass(t *testing.T) {
	repo := &mockRepo{
		justification: &Justification{
			ID:        "just1",
			StudentID: "stu1",
			ParentID:  "parent1",
			Status:    JustificationPending,
			StartDate: time.Now().Add(-48 * time.Hour),
			EndDate:   time.Now().Add(-24 * time.Hour),
		},
		deleteErr: nil,
	}
	userRepo := &mockUserRepo{isGuardian: true}
	svc := makeService(repo, userRepo)
	err := svc.DeleteJustification(context.Background(), "parent1", "just1")
	assert.NoError(t, err)
}

// ---------------------------------------------------------------------------
// GetChildMonthlyBreakdown — guardianship + DB error propagation
// ---------------------------------------------------------------------------

// FIX regression: DB error must propagate distinctly (not masked as 'access denied').
func TestGetChildMonthlyBreakdown_DBError_Propagates(t *testing.T) {
	userRepo := &mockUserRepo{
		isGuardian:    false,
		isGuardianErr: errors.New("connection refused"),
	}
	svc := makeService(&mockRepo{}, userRepo)
	_, err := svc.GetChildMonthlyBreakdown(context.Background(), "parent1", "stu1", "2024-2025")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused",
		"DB error should propagate, not be swallowed as generic 'access denied'")
}

func TestGetChildMonthlyBreakdown_NonGuardian_Forbidden(t *testing.T) {
	userRepo := &mockUserRepo{isGuardian: false, isGuardianErr: nil}
	svc := makeService(&mockRepo{}, userRepo)
	_, err := svc.GetChildMonthlyBreakdown(context.Background(), "parent1", "stu1", "2024-2025")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestGetChildMonthlyBreakdown_Guardian_Pass(t *testing.T) {
	userRepo := &mockUserRepo{isGuardian: true}
	svc := makeService(&mockRepo{}, userRepo)
	_, err := svc.GetChildMonthlyBreakdown(context.Background(), "parent1", "stu1", "2024-2025")
	assert.NoError(t, err)
}

// ---------------------------------------------------------------------------
// RequestJustification — guardianship check
// ---------------------------------------------------------------------------

func TestRequestJustification_NonGuardian_Forbidden(t *testing.T) {
	userRepo := &mockUserRepo{isGuardian: false}
	svc := makeService(&mockRepo{}, userRepo)
	err := svc.RequestJustification(context.Background(), "parent1", JustificationRequest{
		StudentID: "stu1",
		StartDate: "2024-01-08",
		EndDate:   "2024-01-10",
		Reason:    "Illness",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestRequestJustification_Guardian_Pass(t *testing.T) {
	userRepo := &mockUserRepo{isGuardian: true}
	svc := makeService(&mockRepo{hasOverlap: false}, userRepo)
	err := svc.RequestJustification(context.Background(), "parent1", JustificationRequest{
		StudentID: "stu1",
		StartDate: "2024-01-08",
		EndDate:   "2024-01-10",
		Reason:    "Illness",
	})
	assert.NoError(t, err)
}
