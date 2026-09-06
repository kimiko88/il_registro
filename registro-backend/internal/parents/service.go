package parents

import (
	"context"
	"errors"
	"math"
	"registro-backend/internal/attendance"
	"registro-backend/internal/communications"
	"registro-backend/internal/grades"
	"registro-backend/internal/users"
)

var (
	ErrUnauthorized = errors.New("unauthorized: target student is not linked to this parent")
)

type Service struct {
	repo       Repository
	usersRepo  users.Repository
	gradesRepo grades.Repository
	attRepo    attendance.Repository
	commsRepo  communications.Repository
}

func NewService(repo Repository, ur users.Repository, gr grades.Repository, ar attendance.Repository, cr communications.Repository) *Service {
	return &Service{
		repo:       repo,
		usersRepo:  ur,
		gradesRepo: gr,
		attRepo:    ar,
		commsRepo:  cr,
	}
}

func (s *Service) GetDashboard(ctx context.Context, parentUserID string) (*ParentDashboardResponse, error) {
	children, err := s.usersRepo.GetChildren(ctx, parentUserID)
	if err != nil {
		return nil, err
	}

	resp := &ParentDashboardResponse{
		ParentID: parentUserID,
		Children: []ChildOverview{},
	}

	for _, child := range children {
		overview := ChildOverview{
			Student: child,
		}

		// Fetch grades by student user_id
		if gradesList, err := s.gradesRepo.FindByStudent(ctx, child.UserID); err == nil {
			var sum float64
			var count int
			for _, g := range gradesList {
				if g.IsPublished && g.DeletedAt == nil {
					overview.Grades = append(overview.Grades, g)
					sum += g.GradeValue
					count++
				}
			}
			if count > 0 {
				overview.AverageGrade = sum / float64(count)
			}
		}

		// Fetch attendance stats by student profile id
		if stats, err := s.attRepo.GetStats(child.ID); err == nil {
			overview.AttendanceStats = stats
		}

		// Fetch pending circulars
		childSchoolID := ""
		if childUser, err := s.usersRepo.GetByID(ctx, child.UserID); err == nil && childUser != nil && childUser.SchoolID != nil {
			childSchoolID = *childUser.SchoolID
		}
		if msgs, err := s.commsRepo.ListBacheca(ctx, childSchoolID, child.UserID); err == nil {
			for _, m := range msgs {
				if m.RequiresSignature && !m.IsSigned {
					overview.PendingCirculars = append(overview.PendingCirculars, m)
				}
			}
		}

		resp.Children = append(resp.Children, overview)
	}

	return resp, nil
}

func (s *Service) GetChildGradesAverage(ctx context.Context, parentUserID, studentID string) (float64, error) {
	isGuard, err := s.usersRepo.IsGuardian(ctx, parentUserID, studentID)
	if err != nil || !isGuard {
		return 0, ErrUnauthorized
	}

	gradesList, err := s.gradesRepo.FindByStudent(ctx, studentID)
	if err != nil {
		return 0, err
	}

	var weightedSum float64
	var totalWeight float64
	var unweightedSum float64
	var count int

	for _, g := range gradesList {
		if !g.IsPublished || g.DeletedAt != nil {
			continue
		}
		// Skip non-summative grades if summative category is specified
		if g.GradeCategory != "" && g.GradeCategory != grades.GradeCategorySummative {
			continue
		}

		w := g.Weight
		if w <= 0 {
			w = 1.0
		}
		weightedSum += g.GradeValue * w
		totalWeight += w
		unweightedSum += g.GradeValue
		count++
	}

	// Fallback to all published grades if no summative grades found
	if count == 0 {
		for _, g := range gradesList {
			if g.IsPublished && g.DeletedAt == nil {
				w := g.Weight
				if w <= 0 {
					w = 1.0
				}
				weightedSum += g.GradeValue * w
				totalWeight += w
				count++
			}
		}
	}

	if count == 0 || totalWeight == 0 {
		return 0, nil
	}

	avg := weightedSum / totalWeight
	return math.Round(avg*100) / 100, nil
}

func (s *Service) GetDashboardStats(ctx context.Context, parentUserID string) (*ParentDashboardStatsResponse, error) {
	children, err := s.usersRepo.GetChildren(ctx, parentUserID)
	if err != nil {
		return nil, err
	}

	var childUserIDs []string
	for _, child := range children {
		if child.UserID != "" {
			childUserIDs = append(childUserIDs, child.UserID)
		}
	}

	stats, err := s.repo.GetDashboardStats(ctx, parentUserID, childUserIDs)
	if err != nil {
		return nil, err
	}

	stats.ChildrenCount = len(children)
	stats.TotalChildren = len(children)

	return stats, nil
}
