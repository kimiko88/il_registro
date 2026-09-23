package timetablegen

import (
	"context"
	"testing"
	"time"
)

func TestGenerator_SeniorityPriority(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	// Senior Teacher: hired in 2010
	hiringSenior := time.Date(2010, 9, 1, 0, 0, 0, 0, time.UTC)
	// Junior Teacher: hired in 2022
	hiringJunior := time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC)

	assignments := []AssignmentData{
		{
			ClassID:       "class-1",
			ClassName:     "1A",
			SubjectID:     "sub-ita",
			SubjectName:   "Italiano",
			TeacherID:     "teacher-senior",
			TeacherUserID: "user-senior",
			TeacherName:   "Prof Senior",
			HiringDate:    &hiringSenior,
			HoursPerWeek:  4,
		},
		{
			ClassID:       "class-2",
			ClassName:     "2B",
			SubjectID:     "sub-mat",
			SubjectName:   "Matematica",
			TeacherID:     "teacher-junior",
			TeacherUserID: "user-junior",
			TeacherName:   "Prof Junior",
			HiringDate:    &hiringJunior,
			HoursPerWeek:  4,
		},
	}

	// Both teachers want Monday 1st hour as their top preference
	preferences := []TeacherPreference{
		{
			TeacherID:      "teacher-senior",
			DayOfWeek:      1,
			HourIndex:      1,
			PreferenceType: PrefPreferred,
		},
		{
			TeacherID:      "teacher-junior",
			DayOfWeek:      1,
			HourIndex:      1,
			PreferenceType: PrefPreferred,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, preferences, nil)
	if err != nil {
		t.Fatalf("unexpected error during generation: %v", err)
	}

	if result.AssignedSlots != 8 {
		t.Errorf("expected 8 assigned slots, got %d", result.AssignedSlots)
	}

	// Senior teacher should have gotten Monday hour 1
	var seniorMonH1, juniorMonH1 bool
	for _, s := range result.Slots {
		if s.DayOfWeek == 1 && s.HourIndex == 1 {
			if s.TeacherID != nil && *s.TeacherID == "teacher-senior" {
				seniorMonH1 = true
			}
			if s.TeacherID != nil && *s.TeacherID == "teacher-junior" {
				juniorMonH1 = true
			}
		}
	}

	if !seniorMonH1 {
		t.Errorf("expected senior teacher to be assigned Monday 1st hour due to seniority")
	}
	_ = juniorMonH1
}

func TestGenerator_RoomRequirementsAndBuilding(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	buildingA := "building-succursale"
	classBuilding := &buildingA

	assignments := []AssignmentData{
		{
			ClassID:       "class-info",
			ClassName:     "3C",
			BuildingID:    classBuilding,
			SubjectID:     "sub-info",
			SubjectName:   "Informatica",
			TeacherID:     "teacher-info",
			TeacherUserID: "user-info",
			TeacherName:   "Prof Byte",
			HoursPerWeek:  2,
		},
	}

	rooms := []RoomData{
		{
			ID:         "room-lab-succursale",
			BuildingID: &buildingA,
			Name:       "Lab Info Succursale",
			RoomType:   "lab_informatica",
			Capacity:   25,
			IsActive:   true,
		},
		{
			ID:         "room-lab-centrale",
			BuildingID: func(s string) *string { return &s }("building-centrale"),
			Name:       "Lab Info Centrale",
			RoomType:   "lab_informatica",
			Capacity:   30,
			IsActive:   true,
		},
	}

	roomReqs := map[string]SubjectRoomRequirement{
		"sub-info": {
			SubjectID:        "sub-info",
			RequiredRoomType: "lab_informatica",
			IsMandatory:      true,
		},
	}

	result, err := gen.Generate(ctx, assignments, rooms, roomReqs, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error during generation: %v", err)
	}

	if result.AssignedSlots != 2 {
		t.Fatalf("expected 2 assigned slots, got %d", result.AssignedSlots)
	}

	for _, s := range result.Slots {
		if s.RoomID == nil || *s.RoomID != "room-lab-succursale" {
			t.Errorf("expected room 'room-lab-succursale' in matching building, got %v", s.RoomID)
		}
		if s.RoomName != "Lab Info Succursale" {
			t.Errorf("expected RoomName 'Lab Info Succursale', got '%s'", s.RoomName)
		}
	}
}

func TestGenerator_NoDoubleBooking(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	// Single teacher teaching 2 different classes, total 8 hours
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

	// Verify no two slots have the same teacher at the same (day, hour)
	type slotKey struct {
		Day  int
		Hour int
	}
	seenSlots := make(map[slotKey]string)
	for _, s := range result.Slots {
		key := slotKey{Day: s.DayOfWeek, Hour: s.HourIndex}
		if existingClass, exists := seenSlots[key]; exists {
			t.Fatalf("collision detected for teacher: already assigned to class %s at day %d hour %d, trying to assign class %s",
				existingClass, key.Day, key.Hour, s.ClassName)
		}
		seenSlots[key] = s.ClassName
	}
}

func TestGenerator_UnavailableSlotsRespected(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	assignments := []AssignmentData{
		{
			ClassID:      "class-1",
			ClassName:    "1A",
			SubjectID:    "sub-1",
			SubjectName:  "Storia",
			TeacherID:    "teacher-unavail",
			HoursPerWeek: 3,
		},
	}

	// Teacher is explicitly unavailable on Friday hours 1 to 5
	var preferences []TeacherPreference
	for h := 1; h <= 5; h++ {
		preferences = append(preferences, TeacherPreference{
			TeacherID:      "teacher-unavail",
			DayOfWeek:      5,
			HourIndex:      h,
			PreferenceType: PrefUnavailable,
		})
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, preferences, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AssignedSlots != 3 {
		t.Fatalf("expected 3 assigned slots, got %d", result.AssignedSlots)
	}

	for _, s := range result.Slots {
		if s.DayOfWeek == 5 {
			t.Errorf("teacher was assigned to Friday hour %d despite being unavailable", s.HourIndex)
		}
	}
}

func TestGenerator_ClassNoDoubleBooking(t *testing.T) {
	gen := NewGenerator(DefaultConfig())
	ctx := context.Background()

	assignments := []AssignmentData{
		{
			ClassID:      "class-shared",
			ClassName:    "3A",
			SubjectID:    "sub-ita",
			SubjectName:  "Italiano",
			TeacherID:    "teacher-1",
			HoursPerWeek: 5,
		},
		{
			ClassID:      "class-shared",
			ClassName:    "3A",
			SubjectID:    "sub-mat",
			SubjectName:  "Matematica",
			TeacherID:    "teacher-2",
			HoursPerWeek: 4,
		},
		{
			ClassID:      "class-shared",
			ClassName:    "3A",
			SubjectID:    "sub-ing",
			SubjectName:  "Inglese",
			TeacherID:    "teacher-3",
			HoursPerWeek: 3,
		},
	}

	result, err := gen.Generate(ctx, assignments, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AssignedSlots != 12 {
		t.Fatalf("expected 12 assigned slots, got %d", result.AssignedSlots)
	}

	// Verify no two subjects in class 3A are in the same slot
	type slotKey struct {
		Day  int
		Hour int
	}
	classSlots := make(map[slotKey]string)
	for _, s := range result.Slots {
		key := slotKey{Day: s.DayOfWeek, Hour: s.HourIndex}
		if existingSub, exists := classSlots[key]; exists {
			t.Fatalf("class double-booking: %s and %s in slot day %d hour %d", existingSub, s.SubjectName, key.Day, key.Hour)
		}
		classSlots[key] = s.SubjectName
	}
}
