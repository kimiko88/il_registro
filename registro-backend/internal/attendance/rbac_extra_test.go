package attendance

import (
	"context"
	"errors"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// MarkAttendance — substitute path
// ---------------------------------------------------------------------------

// A substitute teacher (isSub=true) must be able to mark attendance even if
// they are not the regular assigned teacher for the class.
func TestMarkAttendance_SubstituteTeacher_Passes(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: false, isSub: true}, &mockUserRepo{})
	err := svc.MarkAttendance(context.Background(), "sub-teacher", "school1",
		CreateAttendanceRequest{
			ClassID:        "class1",
			StudentID:      "stu1",
			Date:           time.Now().Format("2006-01-02"),
			Status:         StatusPresent,
			Hour:           1,
			IsSubstitution: true,
		})
	assert.NoError(t, err)
}

// A teacher who is neither assigned nor a substitute must be rejected.
func TestMarkAttendance_NotAssignedNotSubstitute_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: false, isSub: false}, &mockUserRepo{})
	err := svc.MarkAttendance(context.Background(), "random-teacher", "school1",
		CreateAttendanceRequest{
			ClassID:        "class1",
			StudentID:      "stu1",
			Date:           time.Now().Format("2006-01-02"),
			Status:         StatusPresent,
			Hour:           1,
			IsSubstitution: true,
		})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

// ---------------------------------------------------------------------------
// MarkBulk — regular teacher path
// ---------------------------------------------------------------------------

// An assigned (non-substitute) teacher must pass the bulk-mark flow.
func TestMarkBulk_AssignedTeacher_Passes(t *testing.T) {
	svc := makeService(&mockRepo{isAssigned: true, batchCreateErr: nil}, &mockUserRepo{})
	err := svc.MarkBulk(context.Background(), "teacher1", "school1",
		BulkAttendanceRequest{
			ClassID: "class1",
			Date:    time.Now().Format("2006-01-02"),
			Hour:    1,
			Statuses: []BulkStudentStatus{
				{StudentID: "stu1", Status: StatusPresent},
			},
		})
	assert.NoError(t, err)
}

// A DB error on IsTeacherAssignedToClass in the non-substitute path must also
// propagate and must NOT silently grant access.
func TestMarkBulk_RegularPath_DBErrorOnAssignment_Propagates(t *testing.T) {
	repo := &mockRepo{
		isAssigned:    false,
		isAssignedErr: errors.New("db timeout"),
	}
	svc := makeService(repo, &mockUserRepo{})
	err := svc.MarkBulk(context.Background(), "teacher1", "school1",
		BulkAttendanceRequest{
			ClassID: "class1",
			Date:    time.Now().Format("2006-01-02"),
			Hour:    1,
			Statuses: []BulkStudentStatus{
				{StudentID: "stu1", Status: StatusPresent},
			},
		})
	require.Error(t, err)
	// Must not grant access silently
	assert.NotContains(t, err.Error(), "success")
}

// ---------------------------------------------------------------------------
// UpdateAttendance — success path
// ---------------------------------------------------------------------------

// An assigned teacher updating a record that belongs to the same school and
// class they teach must succeed.
func TestUpdateAttendance_AssignedTeacherSameSchoolClass_Passes(t *testing.T) {
	hour := 1
	repo := &mockRepo{
		findByIDResult: &Attendance{
			ID:       "att1",
			SchoolID: "school1",
			ClassID:  "class1",
			Hour:     &hour,
			Date:     time.Now().Add(-2 * time.Hour),
			Status:   StatusPresent,
		},
		isAssigned: true,
		updateErr:  nil,
	}
	svc := makeService(repo, &mockUserRepo{})
	statusLate := StatusLate
	err := svc.UpdateAttendance(context.Background(), "teacher1", "school1", "att1",
		UpdateAttendanceRequest{Status: &statusLate})
	assert.NoError(t, err)
}

// ---------------------------------------------------------------------------
// GetStudentAttendance — additional role checks
// ---------------------------------------------------------------------------

// Admin, principal, and secretary should be able to access any student's
// attendance without a school-scoping check on their own user record.
func TestGetStudentAttendance_AdminRoles_Pass(t *testing.T) {
	for _, role := range []string{"admin", "superadmin", "principal", "secretary"} {
		role := role
		t.Run(role, func(t *testing.T) {
			svc := makeService(&mockRepo{}, &mockUserRepo{})
			_, err := svc.GetStudentAttendance(context.Background(),
				"actor1", role, "school1", "stu1",
				time.Now().Add(-7*24*time.Hour), time.Now())
			assert.NoError(t, err, "role %s should be authorised", role)
		})
	}
}

// An unknown role must be rejected.
func TestGetStudentAttendance_UnknownRole_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{}, &mockUserRepo{})
	_, err := svc.GetStudentAttendance(context.Background(),
		"actor1", "janitor", "school1", "stu1",
		time.Now().Add(-7*24*time.Hour), time.Now())
	require.Error(t, err)
}

// Empty actorID for a student role must not resolve to another student's record.
func TestGetStudentAttendance_EmptyActorID_Forbidden(t *testing.T) {
	svc := makeService(&mockRepo{}, &mockUserRepo{})
	_, err := svc.GetStudentAttendance(context.Background(),
		"", "student", "school1", "stu1",
		time.Now().Add(-7*24*time.Hour), time.Now())
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// GetClassAttendance — edge cases
// ---------------------------------------------------------------------------

// An empty actorID even with a privileged role must be refused to prevent
// anonymous reads.
func TestGetClassAttendance_EmptyActorID_Unauthorized(t *testing.T) {
	svc := makeService(&mockRepo{}, &mockUserRepo{})
	_, err := svc.GetClassAttendance(context.Background(),
		"", "admin", "school1", "class1", "2024-01-10")
	require.Error(t, err)
}

// A DB error while checking teacher assignment must propagate, not silently
// grant or deny in an ambiguous way.
func TestGetClassAttendance_TeacherAssignmentDBError_Propagates(t *testing.T) {
	repo := &mockRepo{
		isAssigned:    false,
		isAssignedErr: errors.New("connection refused"),
	}
	svc := makeService(repo, &mockUserRepo{})
	_, err := svc.GetClassAttendance(context.Background(),
		"teacher1", "teacher", "school1", "class1", "2024-01-10")
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "forbidden",
		"DB error should surface as an error, not a generic forbidden")
}

// ---------------------------------------------------------------------------
// DeleteClassAttendanceHour — DB error on teacher assignment check
// ---------------------------------------------------------------------------

// When the DB returns an error checking teacher assignment the call must fail
// and must NOT silently grant the delete.
func TestDeleteClassAttendanceHour_TeacherAssignmentDBError_Propagates(t *testing.T) {
	repo := &mockRepo{
		isAssigned:    false,
		isAssignedErr: errors.New("db error"),
	}
	svc := makeService(repo, &mockUserRepo{})
	err := svc.DeleteClassAttendanceHour(context.Background(),
		"teacher1", "teacher", "school1", "class1", "2024-01-10", 1)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "success")
}

// ---------------------------------------------------------------------------
// ProcessJustification — already-processed and privileged roles
// ---------------------------------------------------------------------------

// Attempting to process a justification that is already approved must return
// an error, preventing double-processing.
func TestProcessJustification_AlreadyApproved_ReturnsError(t *testing.T) {
	classID := "class1"
	repo := &mockRepo{
		justification: &Justification{
			ID:        "just2",
			StudentID: "stu1",
			ParentID:  "parent1",
			Status:    JustificationApproved, // already processed
			StartDate: time.Now().Add(-72 * time.Hour),
			EndDate:   time.Now().Add(-24 * time.Hour),
		},
		isAssigned: true,
	}
	userRepo := &mockUserRepo{
		users: map[string]*users.User{
			"stu1": {ID: "stu1", Role: "student", ClassID: &classID},
		},
	}
	svc := makeService(repo, userRepo)
	err := svc.ProcessJustification(context.Background(), "teacher1", "just2", true)
	require.Error(t, err)
}

// Principal and secretary must be able to process justifications without
// being directly assigned to the class.
func TestProcessJustification_PrincipalAndSecretary_Pass(t *testing.T) {
	for _, role := range []string{"principal", "secretary"} {
		role := role
		t.Run(role, func(t *testing.T) {
			classID := "class1"
			repo := &mockRepo{
				justification: &Justification{
					ID:        "just3",
					StudentID: "stu1",
					ParentID:  "parent1",
					Status:    JustificationPending,
					StartDate: time.Now().Add(-48 * time.Hour),
					EndDate:   time.Now().Add(-24 * time.Hour),
				},
				isAssigned: false, // irrelevant for these roles
			}
			userRepo := &mockUserRepo{
				users: map[string]*users.User{
					"stu1": {ID: "stu1", Role: "student", ClassID: &classID},
					"actor1": {ID: "actor1", Role: role},
				},
			}
			svc := makeService(repo, userRepo)
			err := svc.ProcessJustification(context.Background(), "actor1", "just3", true)
			assert.NoError(t, err, "role %s should be able to process justifications", role)
		})
	}
}

// ---------------------------------------------------------------------------
// RequestJustification — overlap and DB error
// ---------------------------------------------------------------------------

// When a guardian submits a justification for dates already covered by a
// pending/approved justification the request must be rejected.
func TestRequestJustification_OverlapDetected_ReturnsError(t *testing.T) {
	userRepo := &mockUserRepo{isGuardian: true}
	svc := makeService(&mockRepo{hasOverlap: true, hasOverlapErr: nil}, userRepo)
	err := svc.RequestJustification(context.Background(), "parent1",
		JustificationRequest{
			StudentID: "stu1",
			StartDate: "2024-01-08",
			EndDate:   "2024-01-10",
			Reason:    "Illness",
		})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "overlap")
}

// A DB error on overlap check must propagate and not be swallowed.
func TestRequestJustification_OverlapDBError_Propagates(t *testing.T) {
	userRepo := &mockUserRepo{isGuardian: true}
	svc := makeService(
		&mockRepo{hasOverlap: false, hasOverlapErr: errors.New("db error")},
		userRepo,
	)
	err := svc.RequestJustification(context.Background(), "parent1",
		JustificationRequest{
			StudentID: "stu1",
			StartDate: "2024-01-08",
			EndDate:   "2024-01-10",
			Reason:    "Illness",
		})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}
