package parents

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/communications"
	"registro-backend/internal/grades"
	"registro-backend/internal/users"
)

type mockParentRepoForTest struct {
	mock.Mock
}

func (m *mockParentRepoForTest) GetChildrenByParentUserID(ctx context.Context, parentUserID string) ([]string, error) {
	args := m.Called(ctx, parentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockParentRepoForTest) GetDashboardStats(ctx context.Context, parentUserID string, childUserIDs []string) (*ParentDashboardStatsResponse, error) {
	args := m.Called(ctx, parentUserID, childUserIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ParentDashboardStatsResponse), args.Error(1)
}

func TestParentsInitialization(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	assert.NotNil(t, svc)
}

func TestGetChildGradesAverage(t *testing.T) {
	ctx := context.Background()

	t.Run("not guardian returns ErrUnauthorized", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, nil, nil)

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(false, nil).Once()

		avg, err := svc.GetChildGradesAverage(ctx, "parent-1", "student-1")
		assert.ErrorIs(t, err, ErrUnauthorized)
		assert.Equal(t, 0.0, avg)
	})

	t.Run("guardian check error returns ErrUnauthorized", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, nil, nil)

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(false, errors.New("db error")).Once()

		avg, err := svc.GetChildGradesAverage(ctx, "parent-1", "student-1")
		assert.ErrorIs(t, err, ErrUnauthorized)
		assert.Equal(t, 0.0, avg)
	})

	t.Run("grades repo error returns error", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, nil, nil)

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(true, nil).Once()
		gRepo.On("FindByStudent", ctx, "student-1").Return(nil, errors.New("grade query failed")).Once()

		avg, err := svc.GetChildGradesAverage(ctx, "parent-1", "student-1")
		assert.Error(t, err)
		assert.Equal(t, 0.0, avg)
	})

	t.Run("summative grades weighted average calculation", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, nil, nil)

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(true, nil).Once()
		now := time.Now()
		gList := []grades.Grade{
			{
				ID:            "g1",
				GradeValue:    8.0,
				Weight:        1.0,
				GradeCategory: grades.GradeCategorySummative,
				IsPublished:   true,
			},
			{
				ID:            "g2",
				GradeValue:    6.0,
				Weight:        0.5,
				GradeCategory: grades.GradeCategorySummative,
				IsPublished:   true,
			},
			{
				ID:            "g3-deleted",
				GradeValue:    2.0,
				Weight:        1.0,
				GradeCategory: grades.GradeCategorySummative,
				IsPublished:   true,
				DeletedAt:     &now,
			},
			{
				ID:            "g4-unpublished",
				GradeValue:    4.0,
				Weight:        1.0,
				GradeCategory: grades.GradeCategorySummative,
				IsPublished:   false,
			},
			{
				ID:            "g5-formative",
				GradeValue:    10.0,
				Weight:        1.0,
				GradeCategory: "formative",
				IsPublished:   true,
			},
		}
		gRepo.On("FindByStudent", ctx, "student-1").Return(gList, nil).Once()

		avg, err := svc.GetChildGradesAverage(ctx, "parent-1", "student-1")
		assert.NoError(t, err)
		// weighted sum = (8.0 * 1.0) + (6.0 * 0.5) = 8.0 + 3.0 = 11.0
		// total weight = 1.0 + 0.5 = 1.5
		// avg = 11.0 / 1.5 = 7.3333... -> 7.33
		assert.Equal(t, 7.33, avg)
	})

	t.Run("fallback to all published grades when no summative grades", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, nil, nil)

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(true, nil).Once()
		gList := []grades.Grade{
			{
				ID:            "g1",
				GradeValue:    7.0,
				Weight:        0, // weight <= 0 should default to 1.0
				GradeCategory: "oral",
				IsPublished:   true,
			},
			{
				ID:            "g2",
				GradeValue:    9.0,
				Weight:        1.0,
				GradeCategory: "written",
				IsPublished:   true,
			},
		}
		gRepo.On("FindByStudent", ctx, "student-1").Return(gList, nil).Once()

		avg, err := svc.GetChildGradesAverage(ctx, "parent-1", "student-1")
		assert.NoError(t, err)
		// (7.0*1.0 + 9.0*1.0)/2 = 8.0
		assert.Equal(t, 8.0, avg)
	})

	t.Run("no published grades returns 0 and nil error", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, nil, nil)

		uRepo.On("IsGuardian", ctx, "parent-1", "student-1").Return(true, nil).Once()
		gRepo.On("FindByStudent", ctx, "student-1").Return([]grades.Grade{}, nil).Once()

		avg, err := svc.GetChildGradesAverage(ctx, "parent-1", "student-1")
		assert.NoError(t, err)
		assert.Equal(t, 0.0, avg)
	})
}

func TestGetDashboardStats(t *testing.T) {
	ctx := context.Background()

	t.Run("returns error when GetChildren fails", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		svc := NewService(nil, uRepo, nil, nil, nil)

		uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild(nil), errors.New("cannot fetch children")).Once()

		stats, err := svc.GetDashboardStats(ctx, "parent-1")
		assert.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("returns error when repo GetDashboardStats fails", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		pRepo := new(mockParentRepoForTest)
		svc := NewService(pRepo, uRepo, nil, nil, nil)

		uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild{
			{ID: "p1", UserID: "u1"},
			{ID: "p2", UserID: ""}, // child without userID
		}, nil).Once()

		pRepo.On("GetDashboardStats", ctx, "parent-1", []string{"u1"}).Return((*ParentDashboardStatsResponse)(nil), errors.New("stats query failed")).Once()

		stats, err := svc.GetDashboardStats(ctx, "parent-1")
		assert.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("success sets ChildrenCount and TotalChildren", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		pRepo := new(mockParentRepoForTest)
		svc := NewService(pRepo, uRepo, nil, nil, nil)

		children := []users.StudentChild{
			{ID: "p1", UserID: "u1"},
			{ID: "p2", UserID: "u2"},
		}
		uRepo.On("GetChildren", ctx, "parent-1").Return(children, nil).Once()

		mockStats := &ParentDashboardStatsResponse{
			UpcomingColloqui:     3,
			UnreadCommunications: 1,
		}
		pRepo.On("GetDashboardStats", ctx, "parent-1", []string{"u1", "u2"}).Return(mockStats, nil).Once()

		stats, err := svc.GetDashboardStats(ctx, "parent-1")
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 2, stats.ChildrenCount)
		assert.Equal(t, 2, stats.TotalChildren)
		assert.Equal(t, 3, stats.UpcomingColloqui)
	})
}

func TestGetDashboard_ErrorCases(t *testing.T) {
	ctx := context.Background()

	t.Run("GetChildren error bubbles up", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		svc := NewService(nil, uRepo, nil, nil, nil)

		uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild(nil), errors.New("fail")).Once()

		dash, err := svc.GetDashboard(ctx, "parent-1")
		assert.Error(t, err)
		assert.Nil(t, dash)
	})

	t.Run("Calculates child grade average correctly", func(t *testing.T) {
		uRepo := new(mockUserRepoForParentTest)
		gRepo := new(mockGradesRepoForParentTest)
		aRepo := new(mockAttRepoForParentTest)
		cRepo := new(mockCommsRepoForParentTest)
		svc := NewService(nil, uRepo, gRepo, aRepo, cRepo)

		uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild{
			{ID: "c1", UserID: "u1", FirstName: "Sara"},
		}, nil).Once()

		gRepo.On("FindByStudent", ctx, "u1").Return([]grades.Grade{
			{GradeValue: 7.0, IsPublished: true},
			{GradeValue: 9.0, IsPublished: true},
		}, nil).Once()

		uRepo.On("GetByID", ctx, "u1").Return((*users.User)(nil), errors.New("user not found")).Once()
		cRepo.On("ListBacheca", ctx, "", "u1").Return([]*communications.Message{}, nil).Once()

		dash, err := svc.GetDashboard(ctx, "parent-1")
		assert.NoError(t, err)
		assert.NotNil(t, dash)
		assert.Len(t, dash.Children, 1)
		assert.Equal(t, 8.0, dash.Children[0].AverageGrade)
		assert.Len(t, dash.Children[0].Grades, 2)
	})
}
