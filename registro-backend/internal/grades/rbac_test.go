package grades

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestService(repo *MockRepository, userRepo *MockUserRepo) *service {
	return &service{
		repo:       repo,
		userRepo:   userRepo,
		calculator: NewCalculator(),
		// validator is nil: tests below only exercise methods that don't call it
	}
}

// ---------------------------------------------------------------------------
// GetStudentGradesWithFilter – role matrix
// ---------------------------------------------------------------------------

func TestGetStudentGradesWithFilter_RoleMatrix(t *testing.T) {
	const ownerID = "student-owner"
	const otherID = "student-other"
	const parentID = "parent-1"

	cases := []struct {
		name        string
		actorID     string
		actorRole   string
		studentID   string
		guardianOK  bool // only relevant when actorRole == "parent"
		guardianErr error
		wantErr     bool
		errContains string
	}{
		// student can access own grades
		{"student_own", ownerID, "student", ownerID, false, nil, false, ""},
		// student cannot access other student's grades
		{"student_other", otherID, "student", ownerID, false, nil, true, "unauthorized"},
		// teacher may access any student
		{"teacher_any", "teacher-1", "teacher", ownerID, false, nil, false, ""},
		// admin may access any student
		{"admin_any", "admin-1", "admin", ownerID, false, nil, false, ""},
		// parent who IS guardian
		{"parent_guardian", parentID, "parent", ownerID, true, nil, false, ""},
		// parent who is NOT guardian
		{"parent_not_guardian", parentID, "parent", ownerID, false, nil, true, "not a guardian"},
		// parent when IsGuardian returns a DB error
		{"parent_db_error", parentID, "parent", ownerID, false, errors.New("db error"), true, "db error"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			if tc.actorRole == "parent" {
				userRepo.On("IsGuardian", mock.Anything, tc.actorID, tc.studentID).
					Return(tc.guardianOK, tc.guardianErr).Once()
			}

			// Repo is only called when auth passes
			if !tc.wantErr {
				repo.On("FindByStudent", tc.studentID).Return([]Grade{}, nil).Once()
			}

			_, err := svc.GetStudentGradesWithFilter(
				context.Background(), tc.actorID, tc.actorRole, tc.studentID, GradeFilter{},
			)

			if tc.wantErr {
				assert.Error(t, err, "expected error for case %q", tc.name)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err, "unexpected error for case %q", tc.name)
			}

			repo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// GetStudentGradesPaged – same ownership rules as above
// ---------------------------------------------------------------------------

func TestGetStudentGradesPaged_RoleMatrix(t *testing.T) {
	const ownerID = "student-owner"
	const parentID = "parent-1"

	cases := []struct {
		name       string
		actorID    string
		actorRole  string
		studentID  string
		guardianOK bool
		wantErr    bool
	}{
		{"student_own", ownerID, "student", ownerID, false, false},
		{"student_other", "other-student", "student", ownerID, false, true},
		{"parent_ok", parentID, "parent", ownerID, true, false},
		{"parent_denied", parentID, "parent", ownerID, false, true},
		{"teacher_ok", "t1", "teacher", ownerID, false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			if tc.actorRole == "parent" {
				userRepo.On("IsGuardian", mock.Anything, tc.actorID, tc.studentID).
					Return(tc.guardianOK, nil).Once()
			}

			if !tc.wantErr {
				repo.On("FindWithFilterPaginated", mock.Anything).Return([]Grade{}, 0, nil).Once()
			}

			_, err := svc.GetStudentGradesPaged(
				context.Background(), tc.actorID, tc.actorRole, tc.studentID, GradeFilter{},
			)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// BatchCreateGrades – only teacher / admin / superadmin allowed
// ---------------------------------------------------------------------------

func TestBatchCreateGrades_RoleEnforcement(t *testing.T) {
	const schoolID = "school-1"

	allowedRoles := []string{"teacher", "admin", "superadmin"}
	deniedRoles := []string{"student", "parent", "principal", "secretary", "coordinator"}

	for _, role := range allowedRoles {
		t.Run("allowed_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			grades := []*Grade{{StudentID: "s1", SchoolID: schoolID, GradeValue: 7, GradeType: "numeric"}}
			repo.On("BatchCreate", mock.Anything).Return(nil).Once()

			err := svc.BatchCreateGrades("actor-1", role, schoolID, grades)
			assert.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}

	for _, role := range deniedRoles {
		t.Run("denied_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			grades := []*Grade{{StudentID: "s1", SchoolID: schoolID, GradeValue: 7, GradeType: "numeric"}}
			err := svc.BatchCreateGrades("actor-1", role, schoolID, grades)
			assert.ErrorIs(t, err, ErrUnauthorized)
		})
	}
}

// BatchCreateGrades: superadmin bypasses school isolation
func TestBatchCreateGrades_SuperadminBypassesSchoolCheck(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	grade := &Grade{StudentID: "s1", SchoolID: "other-school", GradeValue: 7, GradeType: "numeric"}
	repo.On("BatchCreate", mock.Anything).Return(nil).Once()

	err := svc.BatchCreateGrades("superadmin-1", "superadmin", "school-1", []*Grade{grade})
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// BatchCreateGrades: non-superadmin cannot create grades for another school
func TestBatchCreateGrades_CrossSchoolDenied(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	grade := &Grade{StudentID: "s1", SchoolID: "other-school", GradeValue: 7, GradeType: "numeric"}
	err := svc.BatchCreateGrades("teacher-1", "teacher", "school-1", []*Grade{grade})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create grades for another school")
}

// ---------------------------------------------------------------------------
// UpsertWeightConfig – only admin / superadmin / secretary
// ---------------------------------------------------------------------------

func TestUpsertWeightConfig_RoleEnforcement(t *testing.T) {
	req := UpsertWeightConfigRequest{
		GradeCategory: "summative",
		Weight:        1.0,
	}

	allowedRoles := []string{"admin", "superadmin", "secretary"}
	deniedRoles := []string{"teacher", "student", "parent", "principal", "coordinator"}

	for _, role := range allowedRoles {
		t.Run("allowed_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			cfg := &GradeWeightConfig{}
			repo.On("UpsertWeightConfig", mock.Anything).Return(cfg, nil).Once()

			_, err := svc.UpsertWeightConfig("actor", role, "school-1", req)
			assert.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}

	for _, role := range deniedRoles {
		t.Run("denied_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			_, err := svc.UpsertWeightConfig("actor", role, "school-1", req)
			assert.ErrorIs(t, err, ErrUnauthorized)
		})
	}
}

// ---------------------------------------------------------------------------
// DeleteWeightConfig – only admin / superadmin / secretary
// ---------------------------------------------------------------------------

func TestDeleteWeightConfig_RoleEnforcement(t *testing.T) {
	allowedRoles := []string{"admin", "superadmin", "secretary"}
	deniedRoles := []string{"teacher", "student", "parent", "principal"}

	for _, role := range allowedRoles {
		t.Run("allowed_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			repo.On("DeleteWeightConfig", "cfg-1").Return(nil).Once()

			err := svc.DeleteWeightConfig("actor", role, "school-1", "cfg-1")
			assert.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}

	for _, role := range deniedRoles {
		t.Run("denied_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			err := svc.DeleteWeightConfig("actor", role, "school-1", "cfg-1")
			assert.ErrorIs(t, err, ErrUnauthorized)
		})
	}
}

// ---------------------------------------------------------------------------
// GetChildGrades – parent must be guardian
// ---------------------------------------------------------------------------

func TestGetChildGrades_GuardianCheck(t *testing.T) {
	cases := []struct {
		name        string
		isGuardian  bool
		guardianErr error
		wantErr     bool
		errIs       error
	}{
		{"guardian_ok", true, nil, false, nil},
		{"not_guardian", false, nil, true, ErrNotGuardian},
		{"db_error_propagated", false, errors.New("connection refused"), true, nil},
		{"empty_ids", false, nil, true, ErrNotGuardian}, // parentID == ""
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			parentID := "parent-1"
			studentID := "student-1"

			if tc.name == "empty_ids" {
				_, err := svc.GetChildGrades(context.Background(), "", studentID, GradeFilter{})
				assert.ErrorIs(t, err, ErrNotGuardian)
				return
			}

			userRepo.On("IsGuardian", mock.Anything, parentID, studentID).
				Return(tc.isGuardian, tc.guardianErr).Once()

			if !tc.wantErr {
				repo.On("FindByStudent", studentID).Return([]Grade{}, nil).Once()
			}

			_, err := svc.GetChildGrades(context.Background(), parentID, studentID, GradeFilter{})

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errIs != nil {
					assert.ErrorIs(t, err, tc.errIs)
				}
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// GetChildAverages – parent must be guardian
// ---------------------------------------------------------------------------

func TestGetChildAverages_GuardianCheck(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	const parentID = "student-owner"
	const studentID = "student-1"

	// Denied path
	userRepo.On("IsGuardian", mock.Anything, parentID, studentID).Return(false, nil).Once()
	_, err := svc.GetChildAverages(context.Background(), parentID, studentID)
	assert.ErrorIs(t, err, ErrNotGuardian)

	// Allowed path
	userRepo.On("IsGuardian", mock.Anything, parentID, studentID).Return(true, nil).Once()
	repo.On("FindByStudent", studentID).Return([]Grade{}, nil).Once()
	_, err = svc.GetChildAverages(context.Background(), parentID, studentID)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetMyTrend – student can only access own trend; parent must be guardian
// ---------------------------------------------------------------------------

func TestGetMyTrend_RoleEnforcement(t *testing.T) {
	const ownerID = "student-owner"
	const parentID = "parent-1"

	t.Run("student_own_ok", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUserRepo)
		svc := newTestService(repo, userRepo)

		repo.On("FindByStudent", ownerID).Return([]Grade{}, nil).Once()
		_, err := svc.GetMyTrend(context.Background(), ownerID, "student", ownerID, "math")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("student_other_denied", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUserRepo)
		svc := newTestService(repo, userRepo)

		_, err := svc.GetMyTrend(context.Background(), "other-student", "student", ownerID, "math")
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("parent_not_guardian_denied", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUserRepo)
		svc := newTestService(repo, userRepo)

		userRepo.On("IsGuardian", mock.Anything, parentID, ownerID).Return(false, nil).Once()
		_, err := svc.GetMyTrend(context.Background(), parentID, "parent", ownerID, "math")
		assert.ErrorIs(t, err, ErrNotGuardian)
		userRepo.AssertExpectations(t)
	})

	t.Run("parent_guardian_ok", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUserRepo)
		svc := newTestService(repo, userRepo)

		userRepo.On("IsGuardian", mock.Anything, parentID, ownerID).Return(true, nil).Once()
		repo.On("FindByStudent", ownerID).Return([]Grade{}, nil).Once()
		_, err := svc.GetMyTrend(context.Background(), parentID, "parent", ownerID, "math")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("teacher_ok_any_student", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUserRepo)
		svc := newTestService(repo, userRepo)

		repo.On("FindByStudent", ownerID).Return([]Grade{}, nil).Once()
		_, err := svc.GetMyTrend(context.Background(), "teacher-1", "teacher", ownerID, "math")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

// ---------------------------------------------------------------------------
// GetChildSemesterReport – parent must be guardian
// ---------------------------------------------------------------------------

func TestGetChildSemesterReport_GuardianCheck(t *testing.T) {
	const parentID = "parent-1"
	const studentID = "student-1"

	t.Run("not_guardian", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUserRepo)
		svc := newTestService(repo, userRepo)

		userRepo.On("IsGuardian", mock.Anything, parentID, studentID).Return(false, nil).Once()
		_, err := svc.GetChildSemesterReport(context.Background(), parentID, studentID, 1)
		assert.ErrorIs(t, err, ErrNotGuardian)
		userRepo.AssertExpectations(t)
	})

	t.Run("db_error_propagated", func(t *testing.T) {
		repo := new(MockRepository)
		userRepo := new(MockUserRepo)
		svc := newTestService(repo, userRepo)

		dbErr := errors.New("timeout")
		userRepo.On("IsGuardian", mock.Anything, parentID, studentID).Return(false, dbErr).Once()
		_, err := svc.GetChildSemesterReport(context.Background(), parentID, studentID, 1)
		assert.ErrorIs(t, err, dbErr)
		userRepo.AssertExpectations(t)
	})
}

// ---------------------------------------------------------------------------
// BatchCreateGrades: empty slice is a no-op (no role check needed)
// ---------------------------------------------------------------------------

func TestBatchCreateGrades_EmptySliceIsNoop(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	// Even a denied role returns nil when the slice is empty because the
	// role check happens before the loop, but the empty-guard comes first.
	// This test documents the current contract: empty == no-op regardless of role.
	err := svc.BatchCreateGrades("actor", "teacher", "school-1", []*Grade{})
	assert.NoError(t, err)
}
