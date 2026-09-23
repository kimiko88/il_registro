package timetablegen

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type GeneratorConfig struct {
	MaxDaysPerWeek   int
	MaxHoursPerDay   int
	MaxIterations    int
	TimeLimitSeconds int
}

func DefaultConfig() GeneratorConfig {
	return GeneratorConfig{
		MaxDaysPerWeek:   5, // Mon-Fri
		MaxHoursPerDay:   6, // 1 to 6
		MaxIterations:    2000,
		TimeLimitSeconds: 20,
	}
}

// TimetableGenerator executes the timetable generation algorithm
type TimetableGenerator struct {
	config GeneratorConfig
}

func NewGenerator(cfg GeneratorConfig) *TimetableGenerator {
	if cfg.MaxDaysPerWeek <= 0 {
		cfg.MaxDaysPerWeek = 5
	}
	if cfg.MaxHoursPerDay <= 0 {
		cfg.MaxHoursPerDay = 6
	}
	if cfg.MaxIterations <= 0 {
		cfg.MaxIterations = 2000
	}
	if cfg.TimeLimitSeconds <= 0 {
		cfg.TimeLimitSeconds = 20
	}
	return &TimetableGenerator{config: cfg}
}

type teacherSeniorityInfo struct {
	TeacherID   string
	HiringDate  time.Time
	SeniorityPt float64 // higher = more senior
}

// Generate runs greedy allocation sorted by seniority followed by local search
func (g *TimetableGenerator) Generate(
	ctx context.Context,
	assignments []AssignmentData,
	rooms []RoomData,
	roomReqs map[string]SubjectRoomRequirement,
	preferences []TeacherPreference,
	constraints []TimetableConstraint,
) (*TimetableGenerationResult, error) {
	startTime := time.Now()

	// 1. Group assignments by teacher and calculate seniority
	teacherMap := make(map[string]*teacherSeniorityInfo)
	for _, a := range assignments {
		if _, exists := teacherMap[a.TeacherID]; !exists {
			hDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			if a.HiringDate != nil {
				hDate = *a.HiringDate
			}
			teacherMap[a.TeacherID] = &teacherSeniorityInfo{
				TeacherID:  a.TeacherID,
				HiringDate: hDate,
			}
		}
	}

	var teachersList []*teacherSeniorityInfo
	for _, t := range teacherMap {
		teachersList = append(teachersList, t)
	}

	// Sort teachers by hiring_date ASC (earliest hire date = most senior)
	sort.Slice(teachersList, func(i, j int) bool {
		return teachersList[i].HiringDate.Before(teachersList[j].HiringDate)
	})

	totalTeachers := len(teachersList)
	for idx, t := range teachersList {
		// Rank points from 1.0 (junior) up to 3.0 (senior)
		if totalTeachers > 1 {
			t.SeniorityPt = 3.0 - (float64(idx)/float64(totalTeachers))*2.0
		} else {
			t.SeniorityPt = 3.0
		}
	}

	// Preference lookup: teacherID -> day -> hour -> prefType
	prefLookup := make(map[string]map[int]map[int]string)
	for _, p := range preferences {
		if prefLookup[p.TeacherID] == nil {
			prefLookup[p.TeacherID] = make(map[int]map[int]string)
		}
		if prefLookup[p.TeacherID][p.DayOfWeek] == nil {
			prefLookup[p.TeacherID][p.DayOfWeek] = make(map[int]string)
		}
		prefLookup[p.TeacherID][p.DayOfWeek][p.HourIndex] = p.PreferenceType
	}

	// Group rooms by roomType and buildingID
	roomsByTypeAndBuilding := make(map[string]map[string][]RoomData) // roomType -> buildingID -> []rooms
	for _, r := range rooms {
		rt := r.RoomType
		bID := ""
		if r.BuildingID != nil {
			bID = *r.BuildingID
		}
		if roomsByTypeAndBuilding[rt] == nil {
			roomsByTypeAndBuilding[rt] = make(map[string][]RoomData)
		}
		roomsByTypeAndBuilding[rt][bID] = append(roomsByTypeAndBuilding[rt][bID], r)
	}

	// Tracking matrices for hard conflicts:
	// teacherBusy[teacherID][day][hour] = bool
	teacherBusy := make(map[string]map[int]map[int]bool)
	// classBusy[classID][day][hour] = bool
	classBusy := make(map[string]map[int]map[int]bool)
	// roomBusy[roomID][day][hour] = bool
	roomBusy := make(map[string]map[int]map[int]bool)
	// classSubjectDayCount[classID][subjectID][day] = count
	classSubjectDayCount := make(map[string]map[string]map[int]int)

	var generatedSlots []GeneratedSlot
	var unassigned []UnassignedSlot
	var hardConflicts []HardConflict
	var softViolations []SoftViolation
	var warnings []string

	totalRequiredHours := 0
	for _, a := range assignments {
		totalRequiredHours += a.HoursPerWeek
	}

	// Sort assignments: prioritize assignments of senior teachers first
	teacherOrderMap := make(map[string]int)
	for idx, t := range teachersList {
		teacherOrderMap[t.TeacherID] = idx
	}

	sort.Slice(assignments, func(i, j int) bool {
		ti := teacherOrderMap[assignments[i].TeacherID]
		tj := teacherOrderMap[assignments[j].TeacherID]
		if ti != tj {
			return ti < tj // senior teacher assignments first
		}
		return assignments[i].HoursPerWeek > assignments[j].HoursPerWeek
	})

	// ---------------- Phase 1: Greedy Allocation ----------------
	for _, a := range assignments {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		seniority := 1.5
		if tInfo, ok := teacherMap[a.TeacherID]; ok {
			seniority = tInfo.SeniorityPt
		}

		hoursToPlace := a.HoursPerWeek
		reqRoom, hasRoomReq := roomReqs[a.SubjectID]

		for h := 0; h < hoursToPlace; h++ {
			type Candidate struct {
				Day       int
				Hour      int
				RoomID    *string
				RoomName  string
				Score     float64
				PrefState string
			}

			var candidates []Candidate

			for day := 1; day <= g.config.MaxDaysPerWeek; day++ {
				for hour := 1; hour <= g.config.MaxHoursPerDay; hour++ {
					// Hard check 1: Teacher busy?
					if teacherBusy[a.TeacherID] != nil && teacherBusy[a.TeacherID][day] != nil && teacherBusy[a.TeacherID][day][hour] {
						continue
					}
					// Hard check 2: Class busy?
					if classBusy[a.ClassID] != nil && classBusy[a.ClassID][day] != nil && classBusy[a.ClassID][day][hour] {
						continue
					}

					// Hard check 3: Special room requirement
					var chosenRoomID *string
					var chosenRoomName string

					if hasRoomReq && reqRoom.RequiredRoomType != "" {
						rt := reqRoom.RequiredRoomType
						bID := ""
						if a.BuildingID != nil {
							bID = *a.BuildingID
						}
						// Find room in building or general
						availableRoom := findAvailableRoom(rt, bID, day, hour, roomsByTypeAndBuilding, roomBusy)
						if availableRoom == nil && reqRoom.IsMandatory {
							// Mandatory room not available -> cannot place here
							continue
						}
						if availableRoom != nil {
							chosenRoomID = &availableRoom.ID
							chosenRoomName = availableRoom.Name
						}
					}

					// Preference check
					prefState := PrefNeutral
					if prefLookup[a.TeacherID] != nil && prefLookup[a.TeacherID][day] != nil {
						if p, exists := prefLookup[a.TeacherID][day][hour]; exists {
							prefState = p
						}
					}

					// Hard check 4: Teacher explicitly unavailable
					if prefState == PrefUnavailable {
						// Strongly avoid unavailable slot unless no choice
						// We'll skip it in greedy first pass
						continue
					}

					// Score this candidate
					score := 0.0

					// 1. Teacher Preference score weighted by Seniority
					switch prefState {
					case PrefPreferred:
						score += 15.0 * seniority
					case PrefNeutral:
						score += 2.0
					}

					// 2. Spread across days (avoid piling up subject on same day)
					curDayCount := 0
					if classSubjectDayCount[a.ClassID] != nil && classSubjectDayCount[a.ClassID][a.SubjectID] != nil {
						curDayCount = classSubjectDayCount[a.ClassID][a.SubjectID][day]
					}
					if curDayCount == 0 {
						score += 8.0 // Reward new day distribution
					} else if curDayCount >= 2 {
						score -= 10.0 // Penalize > 2 hours of same subject on same day
					}

					// 3. Heavy subjects (Matematica, Fisica, Latino) avoid last hour
					isHeavy := isHeavySubject(a.SubjectName)
					if isHeavy && hour >= g.config.MaxHoursPerDay {
						score -= 12.0
					}

					// 4. Lab / Gym avoid 1st hour if possible (prep time)
					if hasRoomReq && (reqRoom.RequiredRoomType == "palestra" || strings.HasPrefix(reqRoom.RequiredRoomType, "lab_")) {
						if hour == 1 {
							score -= 3.0
						}
					}

					candidates = append(candidates, Candidate{
						Day:       day,
						Hour:      hour,
						RoomID:    chosenRoomID,
						RoomName:  chosenRoomName,
						Score:     score,
						PrefState: prefState,
					})
				}
			}

			// If no candidates found (e.g. all slots unavailable), fallback allowing unavailable with penalty
			if len(candidates) == 0 {
				for day := 1; day <= g.config.MaxDaysPerWeek; day++ {
					for hour := 1; hour <= g.config.MaxHoursPerDay; hour++ {
						if teacherBusy[a.TeacherID] != nil && teacherBusy[a.TeacherID][day] != nil && teacherBusy[a.TeacherID][day][hour] {
							continue
						}
						if classBusy[a.ClassID] != nil && classBusy[a.ClassID][day] != nil && classBusy[a.ClassID][day][hour] {
							continue
						}
						candidates = append(candidates, Candidate{
							Day:       day,
							Hour:      hour,
							Score:     -50.0,
							PrefState: PrefUnavailable,
						})
					}
				}
			}

			if len(candidates) == 0 {
				// Slot could not be assigned!
				unassigned = append(unassigned, UnassignedSlot{
					ClassID:     a.ClassID,
					ClassName:   a.ClassName,
					SubjectID:   a.SubjectID,
					SubjectName: a.SubjectName,
					TeacherID:   a.TeacherID,
					TeacherName: a.TeacherName,
					HoursNeeded: hoursToPlace - h,
					Reason:      "Nessuno slot disponibile libero da sovrapposizioni docente o classe",
				})
				break
			}

			// Pick candidate with best score
			best := candidates[0]
			for _, c := range candidates[1:] {
				if c.Score > best.Score {
					best = c
				}
			}

			// Mark busy matrices
			setBusy(teacherBusy, a.TeacherID, best.Day, best.Hour)
			setBusy(classBusy, a.ClassID, best.Day, best.Hour)
			if best.RoomID != nil {
				setBusy(roomBusy, *best.RoomID, best.Day, best.Hour)
			}

			// Update counts
			if classSubjectDayCount[a.ClassID] == nil {
				classSubjectDayCount[a.ClassID] = make(map[string]map[int]int)
			}
			if classSubjectDayCount[a.ClassID][a.SubjectID] == nil {
				classSubjectDayCount[a.ClassID][a.SubjectID] = make(map[int]int)
			}
			classSubjectDayCount[a.ClassID][a.SubjectID][best.Day]++

			if best.PrefState == PrefUnavailable {
				softViolations = append(softViolations, SoftViolation{
					ConstraintType: "teacher_unavailability",
					TeacherID:      &a.TeacherID,
					ClassID:        &a.ClassID,
					Description:    fmt.Sprintf("Docente %s assegnato a slot non disponibile (Giorno %d Ora %d)", a.TeacherName, best.Day, best.Hour),
					Penalty:        50.0,
				})
			}

			tID := a.TeacherID
			tUID := a.TeacherUserID

			slot := GeneratedSlot{
				ClassID:       a.ClassID,
				ClassName:     a.ClassName,
				BuildingID:    a.BuildingID,
				SubjectID:     a.SubjectID,
				SubjectName:   a.SubjectName,
				TeacherID:     &tID,
				TeacherUserID: &tUID,
				TeacherName:   a.TeacherName,
				DayOfWeek:     best.Day,
				HourIndex:     best.Hour,
				RoomID:        best.RoomID,
				RoomName:      best.RoomName,
			}
			generatedSlots = append(generatedSlots, slot)
		}
	}

	// ---------------- Phase 2: Local Search (Hill Climbing) ----------------
	// Try random swaps of hours for the same teacher or class to maximize preference satisfaction
	timeBudget := time.Duration(g.config.TimeLimitSeconds) * time.Second
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if len(generatedSlots) > 1 {
		iterations := g.config.MaxIterations
		for iter := 0; iter < iterations; iter++ {
			if time.Since(startTime) > timeBudget || ctx.Err() != nil {
				break
			}

			// Pick random slot
			idx1 := rng.Intn(len(generatedSlots))
			idx2 := rng.Intn(len(generatedSlots))
			if idx1 == idx2 {
				continue
			}

			s1 := generatedSlots[idx1]
			s2 := generatedSlots[idx2]

			// Only swap if slots belong to same class OR same teacher, and are in different slots
			if s1.DayOfWeek == s2.DayOfWeek && s1.HourIndex == s2.HourIndex {
				continue
			}

			if s1.ClassID == s2.ClassID {
				// Swap 2 subjects in same class
				// Check if teachers are available in opposite slots
				t1 := ""
				if s1.TeacherID != nil {
					t1 = *s1.TeacherID
				}
				t2 := ""
				if s2.TeacherID != nil {
					t2 = *s2.TeacherID
				}

				t1AvailableAtS2 := !isBusy(teacherBusy, t1, s2.DayOfWeek, s2.HourIndex) || (t1 == t2)
				t2AvailableAtS1 := !isBusy(teacherBusy, t2, s1.DayOfWeek, s1.HourIndex) || (t1 == t2)

				if t1AvailableAtS2 && t2AvailableAtS1 {
					// Evaluate delta score
					currentScore := evalSlotScore(s1, prefLookup, teacherMap) + evalSlotScore(s2, prefLookup, teacherMap)

					cand1 := s1
					cand1.DayOfWeek = s2.DayOfWeek
					cand1.HourIndex = s2.HourIndex

					cand2 := s2
					cand2.DayOfWeek = s1.DayOfWeek
					cand2.HourIndex = s1.HourIndex

					newScore := evalSlotScore(cand1, prefLookup, teacherMap) + evalSlotScore(cand2, prefLookup, teacherMap)

					if newScore > currentScore {
						// Accept swap!
						unsetBusy(teacherBusy, t1, s1.DayOfWeek, s1.HourIndex)
						unsetBusy(teacherBusy, t2, s2.DayOfWeek, s2.HourIndex)

						setBusy(teacherBusy, t1, s2.DayOfWeek, s2.HourIndex)
						setBusy(teacherBusy, t2, s1.DayOfWeek, s1.HourIndex)

						generatedSlots[idx1] = cand1
						generatedSlots[idx2] = cand2
					}
				}
			}
		}
	}

	coveragePct := 100.0
	if totalRequiredHours > 0 {
		coveragePct = (float64(len(generatedSlots)) / float64(totalRequiredHours)) * 100.0
	}

	durationMs := time.Since(startTime).Milliseconds()

	return &TimetableGenerationResult{
		TotalSlots:     totalRequiredHours,
		AssignedSlots:  len(generatedSlots),
		CoveragePct:    coveragePct,
		Slots:          generatedSlots,
		HardConflicts:  hardConflicts,
		SoftViolations: softViolations,
		Warnings:       warnings,
		Unassigned:     unassigned,
		DurationMs:     durationMs,
	}, nil
}

// Helpers

func setBusy(m map[string]map[int]map[int]bool, id string, day, hour int) {
	if m[id] == nil {
		m[id] = make(map[int]map[int]bool)
	}
	if m[id][day] == nil {
		m[id][day] = make(map[int]bool)
	}
	m[id][day][hour] = true
}

func unsetBusy(m map[string]map[int]map[int]bool, id string, day, hour int) {
	if m[id] != nil && m[id][day] != nil {
		m[id][day][hour] = false
	}
}

func isBusy(m map[string]map[int]map[int]bool, id string, day, hour int) bool {
	if id == "" {
		return false
	}
	if m[id] != nil && m[id][day] != nil {
		return m[id][day][hour]
	}
	return false
}

func findAvailableRoom(roomType, buildingID string, day, hour int, roomsByTypeAndBuilding map[string]map[string][]RoomData, roomBusy map[string]map[int]map[int]bool) *RoomData {
	if roomsByTypeAndBuilding[roomType] == nil {
		return nil
	}

	// 1. Try finding in matching building
	if buildingID != "" {
		for _, r := range roomsByTypeAndBuilding[roomType][buildingID] {
			if !isBusy(roomBusy, r.ID, day, hour) {
				cpy := r
				return &cpy
			}
		}
	}

	// 2. Try finding room without specific building (general)
	for _, r := range roomsByTypeAndBuilding[roomType][""] {
		if !isBusy(roomBusy, r.ID, day, hour) {
			cpy := r
			return &cpy
		}
	}

	return nil
}

func isHeavySubject(name string) bool {
	n := strings.ToLower(name)
	return strings.Contains(n, "matematica") ||
		strings.Contains(n, "fisica") ||
		strings.Contains(n, "latino") ||
		strings.Contains(n, "chimica")
}

func evalSlotScore(s GeneratedSlot, prefLookup map[string]map[int]map[int]string, teacherMap map[string]*teacherSeniorityInfo) float64 {
	score := 0.0
	if s.TeacherID == nil {
		return score
	}
	tID := *s.TeacherID
	seniority := 1.5
	if tInfo, ok := teacherMap[tID]; ok {
		seniority = tInfo.SeniorityPt
	}

	if prefLookup[tID] != nil && prefLookup[tID][s.DayOfWeek] != nil {
		pref := prefLookup[tID][s.DayOfWeek][s.HourIndex]
		switch pref {
		case PrefPreferred:
			score += 15.0 * seniority
		case PrefUnavailable:
			score -= 50.0
		default:
			score += 2.0
		}
	}

	if isHeavySubject(s.SubjectName) && s.HourIndex >= 6 {
		score -= 10.0
	}

	return score
}
