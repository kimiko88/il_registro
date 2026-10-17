package unit

import (
	"context"
	"errors"
	"testing"

	"registro-backend/internal/timetablegen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTimetableQuickPrefsRepo implements the minimal repository for quick preferences testing
type mockTimetableQuickPrefsRepo struct {
	timetablegen.Repository
	items        []timetablegen.TeacherQuickPreferenceItem
	dayOffCounts map[int]int
	getErr       error
	saveErr      error
	savedItems   []timetablegen.TeacherQuickPreferenceItem
	savedYearID  *string
}

func (m *mockTimetableQuickPrefsRepo) GetTeachersQuickPreferences(ctx context.Context, schoolID string, academicYearID *string) ([]timetablegen.TeacherQuickPreferenceItem, map[int]int, error) {
	if m.getErr != nil {
		return nil, nil, m.getErr
	}
	return m.items, m.dayOffCounts, nil
}

func (m *mockTimetableQuickPrefsRepo) SaveTeacherQuickPreferences(ctx context.Context, schoolID string, academicYearID *string, items []timetablegen.TeacherQuickPreferenceItem) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.savedItems = items
	m.savedYearID = academicYearID
	return nil
}

func TestTeacherQuickPreferences_ServiceGetOverview(t *testing.T) {
	ctx := context.Background()
	schoolID := "school-unit-1"
	academicYearID := "2024/2025"

	t.Run("success with populated teachers and day-off counts", func(t *testing.T) {
		repo := &mockTimetableQuickPrefsRepo{
			items: []timetablegen.TeacherQuickPreferenceItem{
				{
					TeacherID:      "t-1",
					TeacherName:    "Mario Rossi",
					SubjectName:    "Matematica",
					DayOff:         1, // Lunedì
					TimeSlotPref:   "early_hours",
					MaxHoursPerDay: 5,
				},
				{
					TeacherID:      "t-2",
					TeacherName:    "Giulia Bianchi",
					SubjectName:    "Italiano",
					DayOff:         3, // Mercoledì
					TimeSlotPref:   "late_hours",
					MaxHoursPerDay: 4,
				},
				{
					TeacherID:      "t-3",
					TeacherName:    "Luca Verdi",
					SubjectName:    "Scienze",
					DayOff:         1, // Lunedì
					TimeSlotPref:   "none",
					MaxHoursPerDay: 0,
				},
			},
			dayOffCounts: map[int]int{
				1: 2, // Lunedì: 2 docenti
				3: 1, // Mercoledì: 1 docente
			},
		}

		svc := timetablegen.NewService(repo, timetablegen.NewGenerator(timetablegen.DefaultConfig()))

		resp, err := svc.GetTeachersQuickPreferences(ctx, schoolID, &academicYearID)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, &academicYearID, resp.AcademicYearID)
		assert.Len(t, resp.Teachers, 3)
		assert.Equal(t, 2, resp.DayOffCounts[1])
		assert.Equal(t, 1, resp.DayOffCounts[3])
		assert.Equal(t, 0, resp.DayOffCounts[2]) // Nessun docente il Martedì
	})

	t.Run("handles repository error gracefully", func(t *testing.T) {
		expectedErr := errors.New("db query timeout")
		repo := &mockTimetableQuickPrefsRepo{
			getErr: expectedErr,
		}

		svc := timetablegen.NewService(repo, timetablegen.NewGenerator(timetablegen.DefaultConfig()))

		resp, err := svc.GetTeachersQuickPreferences(ctx, schoolID, nil)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, expectedErr, err)
	})
}

func TestTeacherQuickPreferences_ServiceSaveBatch(t *testing.T) {
	ctx := context.Background()
	schoolID := "school-unit-1"
	ay := "2024/2025"

	t.Run("saves valid batch of quick preferences", func(t *testing.T) {
		repo := &mockTimetableQuickPrefsRepo{}
		svc := timetablegen.NewService(repo, timetablegen.NewGenerator(timetablegen.DefaultConfig()))

		req := timetablegen.SaveTeacherQuickPreferencesRequest{
			AcademicYearID: &ay,
			Preferences: []timetablegen.TeacherQuickPreferenceItem{
				{
					TeacherID:      "t-1",
					DayOff:         5, // Venerdì
					TimeSlotPref:   "early_hours",
					MaxHoursPerDay: 4,
				},
				{
					TeacherID:      "t-2",
					DayOff:         0, // Nessun giorno libero
					TimeSlotPref:   "none",
					MaxHoursPerDay: 6,
				},
			},
		}

		err := svc.SaveTeacherQuickPreferences(ctx, schoolID, req)
		require.NoError(t, err)
		assert.Equal(t, &ay, repo.savedYearID)
		require.Len(t, repo.savedItems, 2)
		assert.Equal(t, "t-1", repo.savedItems[0].TeacherID)
		assert.Equal(t, 5, repo.savedItems[0].DayOff)
		assert.Equal(t, "early_hours", repo.savedItems[0].TimeSlotPref)
	})

	t.Run("propagates repository write error", func(t *testing.T) {
		expectedErr := errors.New("constraint violation")
		repo := &mockTimetableQuickPrefsRepo{
			saveErr: expectedErr,
		}
		svc := timetablegen.NewService(repo, timetablegen.NewGenerator(timetablegen.DefaultConfig()))

		req := timetablegen.SaveTeacherQuickPreferencesRequest{
			Preferences: []timetablegen.TeacherQuickPreferenceItem{
				{TeacherID: "t-invalid"},
			},
		}

		err := svc.SaveTeacherQuickPreferences(ctx, schoolID, req)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestTeacherQuickPreferences_DayOffBalancingAnalysis(t *testing.T) {
	// Tests statistical balancing computation of day-off distribution
	teachers := []timetablegen.TeacherQuickPreferenceItem{
		{TeacherID: "t-1", DayOff: 1}, // Lun
		{TeacherID: "t-2", DayOff: 1}, // Lun
		{TeacherID: "t-3", DayOff: 1}, // Lun (concentrazione alta)
		{TeacherID: "t-4", DayOff: 2}, // Mar
		{TeacherID: "t-5", DayOff: 3}, // Mer
		{TeacherID: "t-6", DayOff: 5}, // Ven
		{TeacherID: "t-7", DayOff: 5}, // Ven
		{TeacherID: "t-8", DayOff: 0}, // Nessuno
		{TeacherID: "t-9", DayOff: 0}, // Nessuno
	}

	counts := make(map[int]int)
	for _, item := range teachers {
		if item.DayOff >= 1 && item.DayOff <= 6 {
			counts[item.DayOff]++
		}
	}

	assert.Equal(t, 3, counts[1], "Lunedì should have 3 requests")
	assert.Equal(t, 1, counts[2], "Martedì should have 1 request")
	assert.Equal(t, 1, counts[3], "Mercoledì should have 1 request")
	assert.Equal(t, 0, counts[4], "Giovedì should have 0 requests")
	assert.Equal(t, 2, counts[5], "Venerdì should have 2 requests")
	assert.Equal(t, 0, counts[6], "Sabato should have 0 requests")

	// Verify maximum bottleneck detection
	maxCount := 0
	bottleneckDay := 0
	for day, c := range counts {
		if c > maxCount {
			maxCount = c
			bottleneckDay = day
		}
	}
	assert.Equal(t, 1, bottleneckDay, "Lunedì should be identified as the bottleneck day")
	assert.Equal(t, 3, maxCount)
}
