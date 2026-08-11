package grades

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------------------------------------------------------------------------
// CreateGrade — role enforcement
// ---------------------------------------------------------------------------

func TestCreateGrade_RoleEnforcement(t *testing.T) {
	const schoolID = "school-1"

	allowedRoles := []string{"teacher", "admin", "superadmin"}
	deniedRoles := []string{"student", "parent", "coordinator", "secretary", ""}

	for _, role := range allowedRoles {
		role := role
		t.Run("allowed_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			repo.On("Create", mock.Anything).Return(nil).Once()

			err := svc.CreateGrade("actor", role, schoolID, CreateGradeRequest{
				StudentID:   "s1",
				Subject:     "math",
				GradeValue:  8,
				GradeType:   "numeric",
				SchoolID:    schoolID,
				AcademicYear: "2024-2025",
			})
			assert.NoError(t, err, "role %s should be able to create grades", role)
			repo.AssertExpectations(t)
		})
	}

	for _, role := range deniedRoles {
		role := role
		t.Run("denied_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			err := svc.CreateGrade("actor", role, schoolID, CreateGradeRequest{
				StudentID:   "s1",
				Subject:     "math",
				GradeValue:  8,
				GradeType:   "numeric",
				SchoolID:    schoolID,
				AcademicYear: "2024-2025",
			})
			assert.ErrorIs(t, err, ErrUnauthorized,
				"role %s should be denied", role)
		})
	}
}

// ---------------------------------------------------------------------------
// UpdateGrade — role enforcement
// ---------------------------------------------------------------------------

func TestUpdateGrade_RoleEnforcement(t *testing.T) {
	const schoolID = "school-1"
	const gradeID = "grade-1"

	allowedRoles := []string{"teacher", "admin", "superadmin"}
	deniedRoles := []string{"student", "parent", "coordinator", "unknown", ""}

	for _, role := range allowedRoles {
		role := role
		t.Run("allowed_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			existing := &Grade{ID: gradeID, SchoolID: schoolID, GradeValue: 7, GradeType: "numeric"}
			repo.On("FindByID", gradeID).Return(existing, nil).Once()
			repo.On("Update", mock.Anything).Return(nil).Once()

			newVal := float64(9)
			err := svc.UpdateGrade("actor", role, schoolID, gradeID,
				UpdateGradeRequest{GradeValue: &newVal})
			assert.NoError(t, err, "role %s should be able to update grades", role)
			repo.AssertExpectations(t)
		})
	}

	for _, role := range deniedRoles {
		role := role
		t.Run("denied_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			newVal := float64(9)
			err := svc.UpdateGrade("actor", role, schoolID, gradeID,
				UpdateGradeRequest{GradeValue: &newVal})
			assert.ErrorIs(t, err, ErrUnauthorized,
				"role %s should be denied", role)
		})
	}
}

// UpdateGrade: school isolation — non-superadmin cannot update a grade from
// another school.
func TestUpdateGrade_CrossSchoolDenied(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	existing := &Grade{ID: "grade-1", SchoolID: "other-school", GradeValue: 7, GradeType: "numeric"}
	repo.On("FindByID", "grade-1").Return(existing, nil).Once()

	newVal := float64(9)
	err := svc.UpdateGrade("teacher-1", "teacher", "school-1", "grade-1",
		UpdateGradeRequest{GradeValue: &newVal})
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteGrade — role enforcement
// ---------------------------------------------------------------------------

func TestDeleteGrade_RoleEnforcement(t *testing.T) {
	const schoolID = "school-1"
	const gradeID = "grade-1"

	allowedRoles := []string{"admin", "superadmin"}
	deniedRoles := []string{"teacher", "student", "parent", "secretary", "principal", ""}

	for _, role := range allowedRoles {
		role := role
		t.Run("allowed_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			existing := &Grade{ID: gradeID, SchoolID: schoolID, GradeValue: 7, GradeType: "numeric"}
			repo.On("FindByID", gradeID).Return(existing, nil).Once()
			repo.On("Delete", gradeID).Return(nil).Once()

			err := svc.DeleteGrade("actor", role, schoolID, gradeID)
			assert.NoError(t, err, "role %s should be able to delete grades", role)
			repo.AssertExpectations(t)
		})
	}

	for _, role := range deniedRoles {
		role := role
		t.Run("denied_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			err := svc.DeleteGrade("actor", role, schoolID, gradeID)
			assert.ErrorIs(t, err, ErrUnauthorized,
				"role %s should be denied", role)
		})
	}
}

// ---------------------------------------------------------------------------
// GetClassGrades — role enforcement
// ---------------------------------------------------------------------------

func TestGetClassGrades_RoleEnforcement(t *testing.T) {
	const schoolID = "school-1"
	const classID = "class-1"

	allowedRoles := []string{"teacher", "admin", "superadmin", "principal", "secretary"}
	deniedRoles := []string{"student", "parent", ""}

	for _, role := range allowedRoles {
		role := role
		t.Run("allowed_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			repo.On("FindByClass", classID).Return([]Grade{}, nil).Once()

			_, err := svc.GetClassGrades(context.Background(), "actor", role, schoolID, classID)
			assert.NoError(t, err, "role %s should be able to read class grades", role)
			repo.AssertExpectations(t)
		})
	}

	for _, role := range deniedRoles {
		role := role
		t.Run("denied_"+role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			_, err := svc.GetClassGrades(context.Background(), "actor", role, schoolID, classID)
			assert.ErrorIs(t, err, ErrUnauthorized,
				"role %s should be denied", role)
		})
	}
}

// ---------------------------------------------------------------------------
// GetStudentStats — ownership + admin bypass
// ---------------------------------------------------------------------------

func TestGetStudentStats_StudentOwnRecord_Pass(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	repo.On("GetStudentStats", "student-1").Return(&StudentStats{}, nil).Once()

	_, err := svc.GetStudentStats(context.Background(), "student-1", "student", "student-1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetStudentStats_StudentOtherRecord_Forbidden(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	_, err := svc.GetStudentStats(context.Background(), "other-student", "student", "student-1")
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestGetStudentStats_Admin_AnyStudent_Pass(t *testing.T) {
	for _, role := range []string{"admin", "superadmin", "teacher"} {
		role := role
		t.Run(role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			repo.On("GetStudentStats", "student-1").Return(&StudentStats{}, nil).Once()

			_, err := svc.GetStudentStats(context.Background(), "actor", role, "student-1")
			assert.NoError(t, err, "role %s should access any student stats", role)
			repo.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// ListWeightConfigs — any authenticated role should be able to read configs
// ---------------------------------------------------------------------------

func TestListWeightConfigs_AnyAuthenticatedRole_Pass(t *testing.T) {
	for _, role := range []string{"teacher", "admin", "superadmin", "student", "parent", "secretary", "principal"} {
		role := role
		t.Run(role, func(t *testing.T) {
			repo := new(MockRepository)
			userRepo := new(MockUserRepo)
			svc := newTestService(repo, userRepo)

			repo.On("ListWeightConfigs", "school-1").Return([]*GradeWeightConfig{}, nil).Once()

			_, err := svc.ListWeightConfigs("actor", role, "school-1")
			assert.NoError(t, err, "role %s should be able to list weight configs", role)
			repo.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// GetChildSemesterReport — empty parentID short-circuits before DB call
// ---------------------------------------------------------------------------

func TestGetChildSemesterReport_EmptyParentID_Forbidden(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	// IsGuardian must NOT be called for an empty parentID
	_, err := svc.GetChildSemesterReport(context.Background(), "", "student-1", 1)
	assert.ErrorIs(t, err, ErrNotGuardian)
	userRepo.AssertNotCalled(t, "IsGuardian", mock.Anything, mock.Anything, mock.Anything)
}

// Guardian with a valid parentID and studentID must succeed.
func TestGetChildSemesterReport_Guardian_Pass(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	userRepo.On("IsGuardian", mock.Anything, "parent-1", "student-1").Return(true, nil).Once()
	repo.On("FindByStudent", "student-1").Return([]Grade{}, nil).Once()

	_, err := svc.GetChildSemesterReport(context.Background(), "parent-1", "student-1", 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetStudentGradesWithFilter — DB error on guardian check propagates
// ---------------------------------------------------------------------------

func TestGetStudentGradesWithFilter_GuardianDBError_Propagates(t *testing.T) {
	repo := new(MockRepository)
	userRepo := new(MockUserRepo)
	svc := newTestService(repo, userRepo)

	dbErr := errors.New("connection refused")
	userRepo.On("IsGuardian", mock.Anything, "parent-1", "student-1").
		Return(false, dbErr).Once()

	_, err := svc.GetStudentGradesWithFilter(
		context.Background(), "parent-1", "parent", "student-1", GradeFilter{},
	)
	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr, "DB error must propagate, not be swallowed")
	userRepo.AssertExpectations(t)
}
