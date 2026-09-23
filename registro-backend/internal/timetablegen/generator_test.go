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
