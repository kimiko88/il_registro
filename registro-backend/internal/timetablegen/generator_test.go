package timetablegen

import (
	"context"
	"testing"
)

// TestGenerator_RoomRequirementsAndLabHours verifies that only the chosen number of lab hours (e.g. 2 out of 4)
// get allocated to the laboratory room, freeing the lab for other classes during remaining hours.
func TestGenerator_RoomRequirementsAndLabHours(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	buildingA := "building-centrale"

	assignments := []AssignmentData{
		{
			ClassID:       "class-sci",
			ClassName:     "1A",
			BuildingID:    &buildingA,
			SubjectID:     "sub-scienze",
			SubjectName:   "Scienze Naturali",
			TeacherID:     "teacher-sci",
			TeacherUserID: "user-sci",
			TeacherName:   "Prof Scienze",
			HoursPerWeek:  4, // 4 hours total per week
		},
	}

	rooms := []RoomData{
		{
			ID:         "room-lab-chimica",
			BuildingID: &buildingA,
			Name:       "Laboratorio di Chimica",
			RoomType:   "lab_chimica",
			Capacity:   25,
			IsActive:   true,
		},
	}

	// Only 2 of the 4 hours should be in the lab!
	roomReqs := map[string]SubjectRoomRequirement{
		"sub-scienze": {
			SubjectID:        "sub-scienze",
			RequiredRoomType: "lab_chimica",
			LabHours:         2,
			IsMandatory:      true,
		},
	}

	result, err := gen.Generate(ctx, assignments, rooms, roomReqs, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error during generation: %v", err)
	}

	if result.AssignedSlots != 4 {
		t.Fatalf("expected 4 assigned slots for class, got %d", result.AssignedSlots)
	}

	labCount := 0
	nonLabCount := 0
	for _, s := range result.Slots {
		if s.RoomID != nil && *s.RoomID == "room-lab-chimica" {
			labCount++
		} else {
			nonLabCount++
		}
	}

	if labCount != 2 {
		t.Errorf("expected exactly 2 lab hours in laboratory, got %d", labCount)
	}
	if nonLabCount != 2 {
		t.Errorf("expected exactly 2 non-lab hours in standard classroom, got %d", nonLabCount)
	}
}

// TestGenerator_UnassignedTeacherChairs verifies that cattedre non assegnate a docenti attualmente assunti
// (docenti da nominare / spezzoni) are fully scheduled for the class without conflicts.
func TestGenerator_UnassignedTeacherChairs(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	assignments := []AssignmentData{
		{
			ClassID:       "class-1",
			ClassName:     "1B",
			SubjectID:     "sub-mat",
			SubjectName:   "Matematica",
			TeacherID:     "unassigned-cs101",
			TeacherUserID: "unassigned-cs101",
			TeacherName:   "Docente da Nominare (Cattedra non assegnata)",
			HoursPerWeek:  4,
		},
		{
			ClassID:       "class-2",
			ClassName:     "2B",
			SubjectID:     "sub-mat",
			SubjectName:   "Matematica",
			TeacherID:     "unassigned-cs102",
			TeacherUserID: "unassigned-cs102",
			TeacherName:   "Docente da Nominare (Cattedra non assegnata)",
			HoursPerWeek:  4,
		},
		{
			ClassID:       "class-1",
			ClassName:     "1B",
			SubjectID:     "sub-spezzone",
			SubjectName:   "Diritto",
			TeacherID:     "spezzone-tc201",
			TeacherUserID: "spezzone-tc201",
			TeacherName:   "Docente da Nominare (Diritto)",
			HoursPerWeek:  2,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AssignedSlots != 10 {
		t.Fatalf("expected 10 assigned slots for unassigned chairs, got %d", result.AssignedSlots)
	}

	if len(result.Unassigned) > 0 {
		t.Errorf("expected 0 unassigned slots, got %d", len(result.Unassigned))
	}
}

// TestGenerator_AssociatedGroupCoTeaching verifies that a teacher cannot be in multiple classes at the same time,
// UNLESS they are teaching an associated group / linguistic group for that subject.
func TestGenerator_AssociatedGroupCoTeaching(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	groupID := "group-spagnolo-3ab"

	// Teacher teaches Spagnolo as an associated group for Class 3A and Class 3B
	assignments := []AssignmentData{
		{
			ClassID:           "class-3a",
			ClassName:         "3A",
			SubjectID:         "sub-spa",
			SubjectName:       "Spagnolo",
			TeacherID:         "teacher-lingue",
			TeacherUserID:     "user-lingue",
			TeacherName:       "Prof Lingua",
			HoursPerWeek:      3,
			AssociatedGroupID: &groupID,
			IsAssociatedGroup: true,
		},
		{
			ClassID:           "class-3b",
			ClassName:         "3B",
			SubjectID:         "sub-spa",
			SubjectName:       "Spagnolo",
			TeacherID:         "teacher-lingue",
			TeacherUserID:     "user-lingue",
			TeacherName:       "Prof Lingua",
			HoursPerWeek:      3,
			AssociatedGroupID: &groupID,
			IsAssociatedGroup: true,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AssignedSlots != 6 {
		t.Fatalf("expected 6 total class-slots assigned, got %d", result.AssignedSlots)
	}

	// Verify that each slot in 3A has an identical (Day, Hour) slot in 3B with the same teacher!
	slots3A := make(map[string]GeneratedSlot)
	slots3B := make(map[string]GeneratedSlot)

	for _, s := range result.Slots {
		slotKey := string(rune(s.DayOfWeek*10 + s.HourIndex))
		switch s.ClassID {
		case "class-3a":
			slots3A[slotKey] = s
		case "class-3b":
			slots3B[slotKey] = s
		}
	}

	if len(slots3A) != 3 || len(slots3B) != 3 {
		t.Fatalf("expected 3 slots each for 3A and 3B, got %d and %d", len(slots3A), len(slots3B))
	}

	for key := range slots3A {
		if _, exists := slots3B[key]; !exists {
			t.Errorf("associated group desynchronized: slot %v scheduled for 3A but not 3B", key)
		}
	}

	// Ensure NO hard conflicts were flagged for this simultaneous presence
	if len(result.HardConflicts) > 0 {
		t.Errorf("expected 0 hard conflicts for associated group co-teaching, got %d: %v", len(result.HardConflicts), result.HardConflicts)
	}
}

// TestGenerator_NoDoubleBookingForNonGroup verifies that normal (non-group) assignments strictly prevent
// a teacher from being scheduled in two classes at the same time.
func TestGenerator_NoDoubleBookingForNonGroup(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	// Teacher teaches Scienze to two different classes independently
	assignments := []AssignmentData{
		{
			ClassID:      "class-1",
			ClassName:    "1A",
			SubjectID:    "sub-sci",
			SubjectName:  "Scienze",
			TeacherID:    "teacher-single",
			HoursPerWeek: 4,
		},
		{
			ClassID:      "class-2",
			ClassName:    "2A",
			SubjectID:    "sub-sci",
			SubjectName:  "Scienze",
			TeacherID:    "teacher-single",
			HoursPerWeek: 4,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AssignedSlots != 8 {
		t.Fatalf("expected 8 assigned slots, got %d", result.AssignedSlots)
	}

	type slotKey struct {
		Day  int
		Hour int
	}
	seenSlots := make(map[slotKey]string)
	for _, s := range result.Slots {
		key := slotKey{Day: s.DayOfWeek, Hour: s.HourIndex}
		if existingClass, exists := seenSlots[key]; exists {
			t.Fatalf("collision detected for non-group teacher: already in class %s at day %d hour %d, trying to assign class %s",
				existingClass, key.Day, key.Hour, s.ClassName)
		}
		seenSlots[key] = s.ClassName
	}
}

// TestGenerator_ClassFullCoverage verifies that every class receives all required hours for each subject,
// with zero class collisions.
func TestGenerator_ClassFullCoverage(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	assignments := []AssignmentData{
		{
			ClassID:      "class-3a",
			ClassName:    "3A",
			SubjectID:    "sub-ita",
			SubjectName:  "Italiano",
			TeacherID:    "teacher-1",
			HoursPerWeek: 5,
		},
		{
			ClassID:      "class-3a",
			ClassName:    "3A",
			SubjectID:    "sub-mat",
			SubjectName:  "Matematica",
			TeacherID:    "teacher-2",
			HoursPerWeek: 4,
		},
		{
			ClassID:      "class-3a",
			ClassName:    "3A",
			SubjectID:    "sub-ing",
			SubjectName:  "Inglese",
			TeacherID:    "teacher-3",
			HoursPerWeek: 3,
		},
		{
			ClassID:      "class-3a",
			ClassName:    "3A",
			SubjectID:    "sub-fil",
			SubjectName:  "Filosofia",
			TeacherID:    "teacher-4",
			HoursPerWeek: 2,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTotal := 5 + 4 + 3 + 2
	if result.AssignedSlots != expectedTotal {
		t.Fatalf("expected %d assigned slots, got %d", expectedTotal, result.AssignedSlots)
	}
	if result.CoveragePct != 100.0 {
		t.Errorf("expected 100%% coverage, got %f%%", result.CoveragePct)
	}

	// Verify no double-booking for the class
	type slotKey struct {
		Day  int
		Hour int
	}
	classSlots := make(map[slotKey]string)
	for _, s := range result.Slots {
		key := slotKey{Day: s.DayOfWeek, Hour: s.HourIndex}
		if existingSub, exists := classSlots[key]; exists {
			t.Fatalf("class collision: %s and %s in same slot day %d hour %d", existingSub, s.SubjectName, key.Day, key.Hour)
		}
		classSlots[key] = s.SubjectName
	}
}

// TestGenerator_ThreeAlternatives_Comprehensive tests that 3 distinct alternatives are generated,
// with proper metadata, IDs, descriptions, valid slots, and composite scores.
func TestGenerator_ThreeAlternatives_Comprehensive(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxIterations = 200
	cfg.TimeLimitSeconds = 5
	gen := NewGenerator(cfg)
	ctx := context.Background()

	assignments := []AssignmentData{
		{
			ClassID:      "class-1a",
			ClassName:    "1A",
			SubjectID:    "sub-mat",
			SubjectName:  "Matematica",
			TeacherID:    "teacher-mat",
			TeacherName:  "Prof Matematica",
			HoursPerWeek: 4,
		},
		{
			ClassID:      "class-1a",
			ClassName:    "1A",
			SubjectID:    "sub-ita",
			SubjectName:  "Italiano",
			TeacherID:    "teacher-ita",
			TeacherName:  "Prof Italiano",
			HoursPerWeek: 5,
		},
		{
			ClassID:      "class-1a",
			ClassName:    "1A",
			SubjectID:    "sub-ing",
			SubjectName:  "Inglese",
			TeacherID:    "teacher-ing",
			TeacherName:  "Prof Inglese",
			HoursPerWeek: 3,
		},
		{
			ClassID:      "class-1a",
			ClassName:    "1A",
			SubjectID:    "sub-mot",
			SubjectName:  "Scienze Motorie",
			TeacherID:    "teacher-mot",
			TeacherName:  "Prof Motoria",
			HoursPerWeek: 2,
		},
		{
			ClassID:      "class-1b",
			ClassName:    "1B",
			SubjectID:    "sub-mat",
			SubjectName:  "Matematica",
			TeacherID:    "teacher-mat",
			TeacherName:  "Prof Matematica",
			HoursPerWeek: 4,
		},
		{
			ClassID:      "class-1b",
			ClassName:    "1B",
			SubjectID:    "sub-ita",
			SubjectName:  "Italiano",
			TeacherID:    "teacher-ita",
			TeacherName:  "Prof Italiano",
			HoursPerWeek: 5,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	// 1. Verify that exactly 3 alternatives are returned
	if len(result.Alternatives) != 3 {
		t.Fatalf("expected 3 alternatives, got %d", len(result.Alternatives))
	}

	expectedStrategies := map[int]string{
		1: "balanced",
		2: "didactic_first",
		3: "compact_teacher",
	}

	for _, alt := range result.Alternatives {
		expectedStrategy, ok := expectedStrategies[alt.ID]
		if !ok {
			t.Errorf("unexpected alternative ID: %d", alt.ID)
		}
		if alt.Strategy != expectedStrategy {
			t.Errorf("alt %d: expected strategy %s, got %s", alt.ID, expectedStrategy, alt.Strategy)
		}
		if alt.Description == "" {
			t.Errorf("alt %d: expected non-empty description", alt.ID)
		}
		if len(alt.Slots) != result.AssignedSlots {
			t.Errorf("alt %d: expected %d slots, got %d", alt.ID, result.AssignedSlots, len(alt.Slots))
		}
		if alt.Score <= 0 {
			t.Errorf("alt %d: expected positive pedagogical score, got %f", alt.ID, alt.Score)
		}
		if alt.CoveragePct <= 0 {
			t.Errorf("alt %d: expected positive coverage pct, got %f", alt.ID, alt.CoveragePct)
		}

		// Ensure no two assignments collide for the same class in this alternative
		type dayHourKey struct {
			class string
			day   int
			hour  int
		}
		seen := make(map[dayHourKey]bool)
		for _, s := range alt.Slots {
			k := dayHourKey{class: s.ClassID, day: s.DayOfWeek, hour: s.HourIndex}
			if seen[k] {
				t.Errorf("alt %d: collision for class %s at day %d hour %d", alt.ID, s.ClassID, s.DayOfWeek, s.HourIndex)
			}
			seen[k] = true
		}
	}
}

// TestGenerator_HelperFunctions tests helper logic such as isHeavySubject,
// isTeacherAdjacent, calculateTeacherGaps, and calculateAlternativeScore.
func TestGenerator_HelperFunctions(t *testing.T) {
	// 1. isHeavySubject
	heavyList := []string{
		"Matematica", "matematica", "MATEMATICA APPLICATA",
		"Italiano", "italiano e latino", "Lingua e Letteratura Italiana",
		"Fisica", "Fisica applicata",
		"Scienze", "Scienze Naturali", "Chimica", "Biologia",
		"Latino", "Greco", "Diritto", "Economia Politica",
	}
	for _, name := range heavyList {
		if !isHeavySubject(name) {
			t.Errorf("expected isHeavySubject(%q) to be true", name)
		}
	}

	lightList := []string{
		"Scienze Motorie", "Educazione Fisica", "Arte e Immagine",
		"Musica", "Religione Cattolica", "Laboratorio",
	}
	for _, name := range lightList {
		if isHeavySubject(name) {
			t.Errorf("expected isHeavySubject(%q) to be false", name)
		}
	}

	// 2. isTeacherAdjacent
	busyMap := make(map[string]map[int]map[int]bool)
	tID := "teacher-test-adj"
	setBusy(busyMap, tID, 1, 2) // Day 1, Hour 2

	// Adjacent to hour 2 are hour 1 and hour 3
	if !isTeacherAdjacent(busyMap, tID, 1, 1) {
		t.Errorf("expected hour 1 to be adjacent to hour 2 on day 1")
	}
	if !isTeacherAdjacent(busyMap, tID, 1, 3) {
		t.Errorf("expected hour 3 to be adjacent to hour 2 on day 1")
	}
	if isTeacherAdjacent(busyMap, tID, 1, 5) {
		t.Errorf("hour 5 should not be adjacent to hour 2")
	}
	if isTeacherAdjacent(busyMap, tID, 2, 2) {
		t.Errorf("different day should not be adjacent")
	}

	// 3. calculateTeacherGaps
	teacherT1 := "t1"
	slots := []GeneratedSlot{
		{TeacherID: &teacherT1, DayOfWeek: 1, HourIndex: 1},
		{TeacherID: &teacherT1, DayOfWeek: 1, HourIndex: 2},
		{TeacherID: &teacherT1, DayOfWeek: 1, HourIndex: 4}, // gap at hour 3!
		{TeacherID: &teacherT1, DayOfWeek: 1, HourIndex: 5},
	}
	gaps := calculateTeacherGaps(slots)
	if gaps != 1 {
		t.Errorf("expected exactly 1 gap, got %d", gaps)
	}

	// Consecutive slots: 0 gaps
	consecutiveSlots := []GeneratedSlot{
		{TeacherID: &teacherT1, DayOfWeek: 2, HourIndex: 1},
		{TeacherID: &teacherT1, DayOfWeek: 2, HourIndex: 2},
		{TeacherID: &teacherT1, DayOfWeek: 2, HourIndex: 3},
	}
	if g := calculateTeacherGaps(consecutiveSlots); g != 0 {
		t.Errorf("expected 0 gaps for consecutive slots, got %d", g)
	}

	// 4. calculateAlternativeScore
	res := TimetableGenerationResult{
		AssignedSlots: 20,
		CoveragePct:   100.0,
		Slots:         consecutiveSlots,
	}
	score := calculateAlternativeScore(&res)
	if score <= 0 || score > 1000 {
		t.Errorf("expected score between 0 and 1000, got %f", score)
	}
}

// TestGenerator_DidacticFirstStrategy verifies that didactic_first favors morning hours for heavy subjects.
func TestGenerator_DidacticFirstStrategy(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxIterations = 300
	cfg.TimeLimitSeconds = 5
	gen := NewGenerator(cfg)
	ctx := context.Background()

	assignments := []AssignmentData{
		{
			ClassID:      "class-1c",
			ClassName:    "1C",
			SubjectID:    "sub-mat",
			SubjectName:  "Matematica",
			TeacherID:    "teacher-mat",
			TeacherName:  "Prof Matematica",
			HoursPerWeek: 4,
		},
		{
			ClassID:      "class-1c",
			ClassName:    "1C",
			SubjectID:    "sub-mot",
			SubjectName:  "Scienze Motorie",
			TeacherID:    "teacher-mot",
			TeacherName:  "Prof Motoria",
			HoursPerWeek: 2,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	var didacticAlt *TimetableAlternative
	for i := range result.Alternatives {
		if result.Alternatives[i].Strategy == "didactic_first" {
			didacticAlt = &result.Alternatives[i]
			break
		}
	}

	if didacticAlt == nil {
		t.Fatalf("didactic_first alternative not found")
	}

	// Check that Matematica slots exist and are scheduled
	matCount := 0
	for _, s := range didacticAlt.Slots {
		if s.SubjectName == "Matematica" {
			matCount++
		}
	}
	if matCount != 4 {
		t.Errorf("expected 4 Matematica slots, got %d", matCount)
	}
}
