package scheduling

import "context"

type AnalyticsService struct {
	repo Repository
}

func NewAnalyticsService(repo Repository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (a *AnalyticsService) GetStats(ctx context.Context, schoolID string) *AnalyticsResponse {
	return &AnalyticsResponse{
		TotalSlots:       100, // Mock
		UtilizationRate:  0.75,
		CancellationRate: 0.05,
	}
}
