package timetablegen

import (
	"context"
	"testing"
)

// TestGenerator_CompactTeacherStrategy_GapMinimization verifies that the compact_teacher
// strategy optimizes slot placement to minimize gaps (buchi orario) between teacher lessons.
func TestGenerator_CompactTeacherStrategy_GapMinimization(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxIterations = 400
	cfg.TimeLimitSeconds = 5
	gen := NewGenerator(cfg)
	ctx := context.Background()

	teacherA := "teacher-compact-1"
	teacherName := "Prof Compatto"

	// Teacher has 6 hours spread across 2 classes (1A and 1B)
	assignments := []AssignmentData{
		{
			ClassID:      "class-1a",
			ClassName:    "1A",
			SubjectID:    "sub-ita",
			SubjectName:  "Italiano",
			TeacherID:    teacherA,
			TeacherName:  teacherName,
			HoursPerWeek: 3,
		},
		{
			ClassID:      "class-1b",
			ClassName:    "1B",
			SubjectID:    "sub-ita",
			SubjectName:  "Italiano",
			TeacherID:    teacherA,
			TeacherName:  teacherName,
			HoursPerWeek: 3,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	var compactAlt *TimetableAlternative
	for i := range result.Alternatives {
		if result.Alternatives[i].Strategy == "compact_teacher" {
			compactAlt = &result.Alternatives[i]
			break
		}
	}

	if compactAlt == nil {
		t.Fatalf("expected compact_teacher alternative to be generated")
	}

	// Verify all 6 hours are assigned
	if len(compactAlt.Slots) != 6 {
		t.Errorf("expected 6 slots in compact alternative, got %d", len(compactAlt.Slots))
	}

	// Check teacher gaps in compact alternative
	gaps := calculateTeacherGaps(compactAlt.Slots)
	if gaps > 2 {
		t.Errorf("expected minimal gaps (<=2) for compact_teacher strategy, got %d", gaps)
	}
}

// TestGenerator_TeacherDayOff_HardCompliance verifies that when a teacher sets a day-off
// preference (e.g. Wednesday = day 3), the generator assigns 0 hours on that day whenever possible.
func TestGenerator_TeacherDayOff_HardCompliance(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxIterations = 500
	cfg.TimeLimitSeconds = 5
	gen := NewGenerator(cfg)
	ctx := context.Background()

	teacherID := "teacher-dayoff-wed"
	teacherName := "Prof MercolediLibero"

	assignments := []AssignmentData{
		{
			ClassID:      "class-2a",
			ClassName:    "2A",
			SubjectID:    "sub-filo",
			SubjectName:  "Filosofia",
			TeacherID:    teacherID,
			TeacherName:  teacherName,
			HoursPerWeek: 4, // 4 hours in a 5-day week
		},
	}

	// Preference: day 3 (Wednesday) is day-off
	prefs := []TeacherPreference{
		{
			TeacherID:      teacherID,
			DayOfWeek:      3,
			HourIndex:      1,
			PreferenceType: PrefUnavailable,
		},
		{
			TeacherID:      teacherID,
			DayOfWeek:      3,
			HourIndex:      2,
			PreferenceType: PrefUnavailable,
		},
		{
			TeacherID:      teacherID,
			DayOfWeek:      3,
			HourIndex:      3,
			PreferenceType: PrefUnavailable,
		},
		{
			TeacherID:      teacherID,
			DayOfWeek:      3,
			HourIndex:      4,
			PreferenceType: PrefUnavailable,
		},
		{
			TeacherID:      teacherID,
			DayOfWeek:      3,
			HourIndex:      5,
			PreferenceType: PrefUnavailable,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, prefs, nil)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	if result.AssignedSlots != 4 {
		t.Fatalf("expected 4 assigned slots, got %d", result.AssignedSlots)
	}

	// Verify no slots were scheduled on Wednesday (day 3)
	for _, s := range result.Slots {
		if s.TeacherID != nil && *s.TeacherID == teacherID {
			if s.DayOfWeek == 3 {
				t.Errorf("expected 0 slots on Wednesday (day 3), but found slot at hour %d", s.HourIndex)
			}
		}
	}
}

// TestGenerator_ParallelClasses_LabContention verifies that when two separate classes
// require the same unique laboratory, the generator never double-books the laboratory.
func TestGenerator_ParallelClasses_LabContention(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxIterations = 400
	cfg.TimeLimitSeconds = 5
	gen := NewGenerator(cfg)
	ctx := context.Background()

	buildingID := "building-succursale"
	labRoomID := "lab-informatica-1"

	rooms := []RoomData{
		{
			ID:         labRoomID,
			BuildingID: &buildingID,
			Name:       "Laboratorio Informatica",
			RoomType:   "lab_info",
			Capacity:   30,
			IsActive:   true,
		},
	}

	roomReqs := map[string]SubjectRoomRequirement{
		"sub-info": {
			SubjectID:        "sub-info",
			RequiredRoomType: "lab_info",
			LabHours:         2,
			IsMandatory:      true,
		},
	}

	assignments := []AssignmentData{
		{
			ClassID:       "class-3a",
			ClassName:     "3A",
			BuildingID:    &buildingID,
			SubjectID:     "sub-info",
			SubjectName:   "Informatica",
			TeacherID:     "teacher-info-1",
			TeacherUserID: "user-info-1",
			TeacherName:   "Prof Info 1",
			HoursPerWeek:  2,
		},
		{
			ClassID:       "class-3b",
			ClassName:     "3B",
			BuildingID:    &buildingID,
			SubjectID:     "sub-info",
			SubjectName:   "Informatica",
			TeacherID:     "teacher-info-2",
			TeacherUserID: "user-info-2",
			TeacherName:   "Prof Info 2",
			HoursPerWeek:  2,
		},
	}

	result, err := gen.Generate(ctx, assignments, rooms, roomReqs, nil, nil)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	if result.AssignedSlots != 4 {
		t.Fatalf("expected 4 assigned slots, got %d", result.AssignedSlots)
	}

	// Check room occupation map (day, hour) -> room
	roomUsage := make(map[[2]int]string)
	for _, s := range result.Slots {
		if s.RoomID != nil && *s.RoomID == labRoomID {
			key := [2]int{s.DayOfWeek, s.HourIndex}
			if existingClass, found := roomUsage[key]; found {
				t.Fatalf("Room collision detected in lab %s on day %d hour %d between %s and %s",
					labRoomID, s.DayOfWeek, s.HourIndex, existingClass, s.ClassName)
			}
			roomUsage[key] = s.ClassName
		}
	}

	if len(roomUsage) != 4 {
		t.Errorf("expected 4 distinct lab slot times, got %d", len(roomUsage))
	}
}

// TestGenerator_ImpossibleConstraint_GracefulBestEffort verifies that when an impossible constraint
// is presented, the generator gracefully terminates without panic or hanging, reporting a best-effort schedule.
func TestGenerator_ImpossibleConstraint_GracefulBestEffort(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxIterations = 100
	cfg.TimeLimitSeconds = 2
	gen := NewGenerator(cfg)
	ctx := context.Background()

	// 50 hours requested in a 30-hour week (5 days x 6 hours max)
	assignments := []AssignmentData{
		{
			ClassID:      "class-overloaded",
			ClassName:    "Overloaded",
			SubjectID:    "sub-overload",
			SubjectName:  "Materia Impossibile",
			TeacherID:    "teacher-overload",
			TeacherName:  "Prof Overload",
			HoursPerWeek: 50,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("generation should not return error on impossible schedule, got: %v", err)
	}

	// Should cap at 30 slots maximum without crashing
	if result.AssignedSlots > 30 {
		t.Errorf("assigned slots exceeded theoretical week capacity of 30, got %d", result.AssignedSlots)
	}

	// Alternatives should still be created
	if len(result.Alternatives) != 3 {
		t.Errorf("expected 3 alternatives even for impossible constraints, got %d", len(result.Alternatives))
	}
}

// TestGenerator_MultiAlternative_Differentiation ensures that the 3 generated alternatives
// have proper metadata, scores, and non-nil slot collections.
func TestGenerator_MultiAlternative_Differentiation(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxIterations = 300
	cfg.TimeLimitSeconds = 4
	gen := NewGenerator(cfg)
	ctx := context.Background()

	assignments := []AssignmentData{
		{
			ClassID:      "class-diff",
			ClassName:    "1D",
			SubjectID:    "sub-mat",
			SubjectName:  "Matematica",
			TeacherID:    "teacher-diff-1",
			TeacherName:  "Prof Diff 1",
			HoursPerWeek: 4,
		},
		{
			ClassID:      "class-diff",
			ClassName:    "1D",
			SubjectID:    "sub-art",
			SubjectName:  "Arte",
			TeacherID:    "teacher-diff-2",
			TeacherName:  "Prof Diff 2",
			HoursPerWeek: 2,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	if len(result.Alternatives) != 3 {
		t.Fatalf("expected 3 alternatives, got %d", len(result.Alternatives))
	}

	strategies := make(map[string]bool)
	for _, alt := range result.Alternatives {
		strategies[alt.Strategy] = true
		if alt.Score <= 0 {
			t.Errorf("expected positive score for alternative %s, got %f", alt.Strategy, alt.Score)
		}
		if len(alt.Slots) != 6 {
			t.Errorf("expected 6 slots in alternative %s, got %d", alt.Strategy, len(alt.Slots))
		}
		if alt.ID < 1 || alt.ID > 3 {
			t.Errorf("invalid alternative ID %d", alt.ID)
		}
	}

	if !strategies["balanced"] || !strategies["didactic_first"] || !strategies["compact_teacher"] {
		t.Errorf("missing expected strategy from set: %v", strategies)
	}
}
