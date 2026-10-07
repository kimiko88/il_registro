package privacy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateTrafficLight(t *testing.T) {
	// 1. Photo/video granted -> GREEN
	assert.Equal(t, TrafficLightGreen, EvaluateTrafficLight(true, true, true))
	assert.Equal(t, TrafficLightGreen, EvaluateTrafficLight(true, false, false))

	// 2. Educational granted, but no photo/video -> YELLOW
	assert.Equal(t, TrafficLightYellow, EvaluateTrafficLight(false, true, true))
	assert.Equal(t, TrafficLightYellow, EvaluateTrafficLight(false, true, false))
	assert.Equal(t, TrafficLightYellow, EvaluateTrafficLight(false, false, true))

	// 3. All denied -> RED
	assert.Equal(t, TrafficLightRed, EvaluateTrafficLight(false, false, false))
}

func TestPrivacyConsentsLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	// Save consent for student 1 (full green)
	c1, err := svc.SaveStudentConsent(ctx, "school-1", "parent-1", SaveConsentRequest{
		StudentID:               "std-1",
		SchoolYear:              "2025/2026",
		PhotoVideoSocialConsent: true,
		CloudWorkspaceConsent:   true,
		WalkingTripsConsent:     true,
	})
	require.NoError(t, err)
	assert.Equal(t, TrafficLightGreen, c1.TrafficLightBadge)

	// Save consent for student 2 (yellow: no photos)
	c2, err := svc.SaveStudentConsent(ctx, "school-1", "parent-2", SaveConsentRequest{
		StudentID:               "std-2",
		SchoolYear:              "2025/2026",
		PhotoVideoSocialConsent: false,
		CloudWorkspaceConsent:   true,
		WalkingTripsConsent:     true,
	})
	require.NoError(t, err)
	assert.Equal(t, TrafficLightYellow, c2.TrafficLightBadge)

	// Batch query for class roll-call (std-1, std-2, and std-3 unrecorded)
	badges, err := svc.GetClassBadges(ctx, []string{"std-1", "std-2", "std-3"}, "2025/2026")
	require.NoError(t, err)
	assert.Equal(t, TrafficLightGreen, badges["std-1"])
	assert.Equal(t, TrafficLightYellow, badges["std-2"])
	assert.Equal(t, TrafficLightRed, badges["std-3"]) // unrecorded defaults to RED
}
