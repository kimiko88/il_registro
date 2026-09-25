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
		MaxIterations:    10000,
		TimeLimitSeconds: 60,
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
		cfg.MaxIterations = 10000
	}
	if cfg.TimeLimitSeconds <= 0 {
		cfg.TimeLimitSeconds = 60
	}
	return &TimetableGenerator{config: cfg}
}

// Generate runs constrained timetable allocation:
// 1. Respects laboratory constraints and allocates the chosen number of lab hours per subject.
// 2. Ignores teacher personal desiderata (neutral availability, structural rules only).
// 3. Fully supports cattedre non assegnate a docenti attualmente assunti (docenti da nominare/spezzoni).
// 4. Guarantees that each class receives all required hours for each subject.
// 5. Enforces that no teacher can be in multiple classes simultaneously EXCEPT when teaching an associated or linguistic group.
func (g *TimetableGenerator) Generate(
	ctx context.Context,
	assignments []AssignmentData,
	rooms []RoomData,
	roomReqs map[string]SubjectRoomRequirement,
	preferences []TeacherPreference,
	constraints []TimetableConstraint,
) (*TimetableGenerationResult, error) {
	startTime := time.Now()

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
	// teacherGroupBusy[teacherID][day][hour] = groupID (tracks which associated group is active)
	teacherGroupBusy := make(map[string]map[int]map[int]string)
	// classBusy[classID][day][hour] = bool
	classBusy := make(map[string]map[int]map[int]bool)
	// roomBusy[roomID][day][hour] = bool
	roomBusy := make(map[string]map[int]map[int]bool)
	// classSubjectDayCount[classID][subjectID][day] = count
	classSubjectDayCount := make(map[string]map[string]map[int]int)

	// Index teacher preferences and day-off requests
	teacherPrefMap := make(map[string]map[int]map[int]string)
	teacherDayOffMap := make(map[string]map[int]int)
	for _, p := range preferences {
		tID := p.TeacherID
		if teacherPrefMap[tID] == nil {
			teacherPrefMap[tID] = make(map[int]map[int]string)
		}
		if teacherPrefMap[tID][p.DayOfWeek] == nil {
			teacherPrefMap[tID][p.DayOfWeek] = make(map[int]string)
		}
		teacherPrefMap[tID][p.DayOfWeek][p.HourIndex] = p.PreferenceType
		if p.PreferenceType == PrefUnavailable {
			if teacherDayOffMap[tID] == nil {
				teacherDayOffMap[tID] = make(map[int]int)
			}
			teacherDayOffMap[tID][p.DayOfWeek]++
		}
	}

	var generatedSlots []GeneratedSlot
	var unassigned []UnassignedSlot
	var hardConflicts []HardConflict
	var softViolations []SoftViolation
	var warnings []string

	totalRequiredHours := 0
	for _, a := range assignments {
		totalRequiredHours += a.HoursPerWeek
	}

	// 1. Separate Associated Groups from Regular Assignments
	associatedGroupsMap := make(map[string][]AssignmentData)
	var regularAssignments []AssignmentData

	for _, a := range assignments {
		if a.IsAssociatedGroup && a.AssociatedGroupID != nil && *a.AssociatedGroupID != "" {
			gID := *a.AssociatedGroupID
			associatedGroupsMap[gID] = append(associatedGroupsMap[gID], a)
		} else {
			regularAssignments = append(regularAssignments, a)
		}
	}

	// ---------------- Phase 1A: Schedule Associated Groups (Gruppi Linguistici / Articolati) ----------------
	// Associated groups combine multiple classes simultaneously with the same teacher and subject.
	for groupID, groupAssigns := range associatedGroupsMap {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if len(groupAssigns) == 0 {
			continue
		}

		lead := groupAssigns[0]
		hoursToPlace := lead.HoursPerWeek
		reqRoom, hasRoomReq := roomReqs[lead.SubjectID]

		labHoursLimit := 0
		if hasRoomReq && reqRoom.RequiredRoomType != "" {
			labHoursLimit = reqRoom.LabHours
			if labHoursLimit <= 0 || labHoursLimit > hoursToPlace {
				labHoursLimit = hoursToPlace
			}
		}

		for h := 0; h < hoursToPlace; h++ {
			isLabHour := hasRoomReq && reqRoom.RequiredRoomType != "" && h < labHoursLimit

			type GroupCandidate struct {
				Day      int
				Hour     int
				RoomID   *string
				RoomName string
				Score    float64
			}

			var candidates []GroupCandidate

			for day := 1; day <= g.config.MaxDaysPerWeek; day++ {
				for hour := 1; hour <= g.config.MaxHoursPerDay; hour++ {
					// Check all participating classes are free
					classesFree := true
					for _, ga := range groupAssigns {
						if isBusy(classBusy, ga.ClassID, day, hour) {
							classesFree = false
							break
						}
					}
					if !classesFree {
						continue
					}

					// Check teacher is free or already allocated to this exact group
					if isBusy(teacherBusy, lead.TeacherID, day, hour) {
						if teacherGroupBusy[lead.TeacherID] == nil ||
							teacherGroupBusy[lead.TeacherID][day] == nil ||
							teacherGroupBusy[lead.TeacherID][day][hour] != groupID {
							continue
						}
					}

					// Room check if laboratory hour
					var chosenRoomID *string
					var chosenRoomName string
					if isLabHour {
						rt := reqRoom.RequiredRoomType
						bID := ""
						if lead.BuildingID != nil {
							bID = *lead.BuildingID
						}
						availableRoom := findAvailableRoom(rt, bID, day, hour, roomsByTypeAndBuilding, roomBusy)
						if availableRoom == nil && reqRoom.IsMandatory {
							continue
						}
						if availableRoom != nil {
							chosenRoomID = &availableRoom.ID
							chosenRoomName = availableRoom.Name
						}
					}

					// Score candidate
					score := 20.0 // Group bonus
					// Spread across days for classes
					curDayCount := 0
					for _, ga := range groupAssigns {
						if classSubjectDayCount[ga.ClassID] != nil && classSubjectDayCount[ga.ClassID][ga.SubjectID] != nil {
							if cnt := classSubjectDayCount[ga.ClassID][ga.SubjectID][day]; cnt > curDayCount {
								curDayCount = cnt
							}
						}
					}
					if curDayCount == 0 {
						score += 10.0
					} else if curDayCount >= 2 {
						score -= 15.0
					}

					if isHeavySubject(lead.SubjectName) && hour >= g.config.MaxHoursPerDay {
						score -= 12.0
					}
					if isLabHour && hour == 1 {
						score -= 3.0
					}

					// Teacher preference score
					pref := getTeacherPref(teacherPrefMap, lead.TeacherID, lead.TeacherUserID, day, hour)
					switch pref {
					case PrefPreferred:
						score += 15.0
					case PrefUnavailable:
						score -= 35.0
					}
					if getTeacherDayOffHours(teacherDayOffMap, lead.TeacherID, lead.TeacherUserID, day) >= 4 {
						score -= 40.0
					}

					candidates = append(candidates, GroupCandidate{
						Day:      day,
						Hour:     hour,
						RoomID:   chosenRoomID,
						RoomName: chosenRoomName,
						Score:    score,
					})
				}
			}

			if len(candidates) == 0 {
				for _, ga := range groupAssigns {
					unassigned = append(unassigned, UnassignedSlot{
						ClassID:     ga.ClassID,
						ClassName:   ga.ClassName,
						SubjectID:   ga.SubjectID,
						SubjectName: ga.SubjectName,
						TeacherID:   ga.TeacherID,
						TeacherName: ga.TeacherName,
						HoursNeeded: hoursToPlace - h,
						Reason:      "Nessuno slot comune disponibile per il gruppo linguistico/associato",
					})
				}
				break
			}

			// Pick best candidate
			best := candidates[0]
			for _, c := range candidates[1:] {
				if c.Score > best.Score {
					best = c
				}
			}

			// Mark busy matrices and record slots for all group classes
			setBusy(teacherBusy, lead.TeacherID, best.Day, best.Hour)
			setTeacherGroup(teacherGroupBusy, lead.TeacherID, best.Day, best.Hour, groupID)
			if best.RoomID != nil {
				setBusy(roomBusy, *best.RoomID, best.Day, best.Hour)
			}

			for _, ga := range groupAssigns {
				setBusy(classBusy, ga.ClassID, best.Day, best.Hour)

				if classSubjectDayCount[ga.ClassID] == nil {
					classSubjectDayCount[ga.ClassID] = make(map[string]map[int]int)
				}
				if classSubjectDayCount[ga.ClassID][ga.SubjectID] == nil {
					classSubjectDayCount[ga.ClassID][ga.SubjectID] = make(map[int]int)
				}
				classSubjectDayCount[ga.ClassID][ga.SubjectID][best.Day]++

				tID := ga.TeacherID
				tUID := ga.TeacherUserID

				generatedSlots = append(generatedSlots, GeneratedSlot{
					ClassID:       ga.ClassID,
					ClassName:     ga.ClassName,
					BuildingID:    ga.BuildingID,
					SubjectID:     ga.SubjectID,
					SubjectName:   ga.SubjectName,
					TeacherID:     &tID,
					TeacherUserID: &tUID,
					TeacherName:   ga.TeacherName,
					DayOfWeek:     best.Day,
					HourIndex:     best.Hour,
					RoomID:        best.RoomID,
					RoomName:      best.RoomName,
				})
			}
		}
	}

	// ---------------- Phase 1B: Schedule Regular Assignments ----------------
	// Sort regular assignments: prioritize subjects requiring special laboratories first, then by hours descending
	sort.Slice(regularAssignments, func(i, j int) bool {
		_, reqI := roomReqs[regularAssignments[i].SubjectID]
		_, reqJ := roomReqs[regularAssignments[j].SubjectID]
		if reqI != reqJ {
			return reqI // lab subjects scheduled first
		}
		return regularAssignments[i].HoursPerWeek > regularAssignments[j].HoursPerWeek
	})

	for _, a := range regularAssignments {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		hoursToPlace := a.HoursPerWeek
		reqRoom, hasRoomReq := roomReqs[a.SubjectID]

		labHoursLimit := 0
		if hasRoomReq && reqRoom.RequiredRoomType != "" {
			labHoursLimit = reqRoom.LabHours
			if labHoursLimit <= 0 || labHoursLimit > hoursToPlace {
				labHoursLimit = hoursToPlace
			}
		}

		isUnassignedTeacher := strings.HasPrefix(a.TeacherID, "unassigned-") || strings.HasPrefix(a.TeacherID, "spezzone-")

		for h := 0; h < hoursToPlace; h++ {
			isLabHour := hasRoomReq && reqRoom.RequiredRoomType != "" && h < labHoursLimit

			type Candidate struct {
				Day      int
				Hour     int
				RoomID   *string
				RoomName string
				Score    float64
			}

			var candidates []Candidate

			for day := 1; day <= g.config.MaxDaysPerWeek; day++ {
				for hour := 1; hour <= g.config.MaxHoursPerDay; hour++ {
					// Hard check 1: Class busy?
					if isBusy(classBusy, a.ClassID, day, hour) {
						continue
					}

					// Hard check 2: Teacher busy?
					// For assigned teachers: cannot be in another class at the same time
					if !isUnassignedTeacher && isBusy(teacherBusy, a.TeacherID, day, hour) {
						continue
					}
					// For specific unassigned chair (spezzone): check against its own chair id
					if isUnassignedTeacher && isBusy(teacherBusy, a.TeacherID, day, hour) {
						continue
					}

					// Hard check 3: Laboratory room availability (only during lab hours)
					var chosenRoomID *string
					var chosenRoomName string

					if isLabHour {
						rt := reqRoom.RequiredRoomType
						bID := ""
						if a.BuildingID != nil {
							bID = *a.BuildingID
						}
						availableRoom := findAvailableRoom(rt, bID, day, hour, roomsByTypeAndBuilding, roomBusy)
						if availableRoom == nil && reqRoom.IsMandatory {
							continue
						}
						if availableRoom != nil {
							chosenRoomID = &availableRoom.ID
							chosenRoomName = availableRoom.Name
						}
					}

					// Score candidate based on pedagogical rules (no teacher personal preferences)
					score := 10.0

					// 1. Spread across days (avoid piling up subject on same day)
					curDayCount := 0
					if classSubjectDayCount[a.ClassID] != nil && classSubjectDayCount[a.ClassID][a.SubjectID] != nil {
						curDayCount = classSubjectDayCount[a.ClassID][a.SubjectID][day]
					}
					if curDayCount == 0 {
						score += 10.0
					} else if curDayCount >= 2 {
						score -= 15.0
					}

					// 2. Heavy subjects (Matematica, Fisica, Latino) avoid last hour
					isHeavy := isHeavySubject(a.SubjectName)
					if isHeavy && hour >= g.config.MaxHoursPerDay {
						score -= 12.0
					}

					// 3. Lab / Gym avoid 1st hour if possible
					if isLabHour && hour == 1 {
						score -= 3.0
					}

					// 4. Teacher preference score (Desiderata)
					pref := getTeacherPref(teacherPrefMap, a.TeacherID, a.TeacherUserID, day, hour)
					switch pref {
					case PrefPreferred:
						score += 15.0
					case PrefUnavailable:
						score -= 35.0
					}
					if getTeacherDayOffHours(teacherDayOffMap, a.TeacherID, a.TeacherUserID, day) >= 4 {
						score -= 40.0
					}

					candidates = append(candidates, Candidate{
						Day:      day,
						Hour:     hour,
						RoomID:   chosenRoomID,
						RoomName: chosenRoomName,
						Score:    score,
					})
				}
			}

			// If no candidates found, attempt intelligent repair/swap within the class
			if len(candidates) == 0 {
				repaired := false
				// Try to find a slot in this class currently occupied by another subject that has no lab constraint,
				// where that other subject can be moved elsewhere and our teacher is free in that slot.
				for slotIdx, existing := range generatedSlots {
					if existing.ClassID != a.ClassID {
						continue
					}
					// Only swap with regular non-lab slots
					if existing.RoomID != nil {
						continue
					}
					// Teacher of assignment a must be free at existing slot
					day := existing.DayOfWeek
					hour := existing.HourIndex
					if !isUnassignedTeacher && isBusy(teacherBusy, a.TeacherID, day, hour) {
						continue
					}

					// Now see if existing subject's teacher can be moved to an alternate empty slot in this class
					existingTeacherID := ""
					if existing.TeacherID != nil {
						existingTeacherID = *existing.TeacherID
					}

					for altDay := 1; altDay <= g.config.MaxDaysPerWeek; altDay++ {
						for altHour := 1; altHour <= g.config.MaxHoursPerDay; altHour++ {
							if isBusy(classBusy, a.ClassID, altDay, altHour) {
								continue
							}
							if existingTeacherID != "" && isBusy(teacherBusy, existingTeacherID, altDay, altHour) {
								continue
							}

							// Valid swap found!
							// Move existing slot to (altDay, altHour)
							unsetBusy(classBusy, a.ClassID, day, hour)
							if existingTeacherID != "" {
								unsetBusy(teacherBusy, existingTeacherID, day, hour)
								setBusy(teacherBusy, existingTeacherID, altDay, altHour)
							}
							setBusy(classBusy, a.ClassID, altDay, altHour)

							generatedSlots[slotIdx].DayOfWeek = altDay
							generatedSlots[slotIdx].HourIndex = altHour

							// Now slot (day, hour) is free for assignment a!
							candidates = append(candidates, Candidate{
								Day:   day,
								Hour:  hour,
								Score: 5.0,
							})
							repaired = true
							break
						}
						if repaired {
							break
						}
					}
					if repaired {
						break
					}
				}
			}

			if len(candidates) == 0 {
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

	// ---------------- Phase 2: Local Search / Pedagogical & Preference Balancing ----------------
	timeBudget := time.Duration(g.config.TimeLimitSeconds) * time.Second
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if len(generatedSlots) > 1 {
		iterations := g.config.MaxIterations
		if g.config.TimeLimitSeconds >= 60 {
			iterations = g.config.TimeLimitSeconds * 2000
		}
		for iter := 0; iter < iterations; iter++ {
			if iter%50 == 0 {
				if time.Since(startTime) > timeBudget || ctx.Err() != nil {
					break
				}
			}

			idx1 := rng.Intn(len(generatedSlots))
			idx2 := rng.Intn(len(generatedSlots))
			if idx1 == idx2 {
				continue
			}

			s1 := generatedSlots[idx1]
			s2 := generatedSlots[idx2]

			// Don't swap associated group slots or slots in different classes
			if s1.ClassID != s2.ClassID {
				continue
			}
			if s1.DayOfWeek == s2.DayOfWeek && s1.HourIndex == s2.HourIndex {
				continue
			}
			// Don't swap if either requires a specific lab room
			if s1.RoomID != nil || s2.RoomID != nil {
				continue
			}

			t1 := ""
			if s1.TeacherID != nil {
				t1 = *s1.TeacherID
			}
			t2 := ""
			if s2.TeacherID != nil {
				t2 = *s2.TeacherID
			}

			t1Available := !isBusy(teacherBusy, t1, s2.DayOfWeek, s2.HourIndex) || (t1 == t2)
			t2Available := !isBusy(teacherBusy, t2, s1.DayOfWeek, s1.HourIndex) || (t1 == t2)

			if t1Available && t2Available {
				currentScore := evalTotalSlotScore(s1, teacherPrefMap, teacherDayOffMap) + evalTotalSlotScore(s2, teacherPrefMap, teacherDayOffMap)

				cand1 := s1
				cand1.DayOfWeek = s2.DayOfWeek
				cand1.HourIndex = s2.HourIndex

				cand2 := s2
				cand2.DayOfWeek = s1.DayOfWeek
				cand2.HourIndex = s1.HourIndex

				newScore := evalTotalSlotScore(cand1, teacherPrefMap, teacherDayOffMap) + evalTotalSlotScore(cand2, teacherPrefMap, teacherDayOffMap)

				if newScore > currentScore {
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

	// ---------------- Phase 3: Validation of Conflict Rules ----------------
	// Rule: Ogni docente non può essere in più classi contemporaneamente,
	// tranne se ha un gruppo linguistico o un gruppo associato per quella materia.
	teacherSlotsAtTime := make(map[string][]GeneratedSlot)
	for _, slot := range generatedSlots {
		if slot.TeacherID != nil && !strings.HasPrefix(*slot.TeacherID, "unassigned-") && !strings.HasPrefix(*slot.TeacherID, "spezzone-") {
			key := fmt.Sprintf("%s-%d-%d", *slot.TeacherID, slot.DayOfWeek, slot.HourIndex)
			teacherSlotsAtTime[key] = append(teacherSlotsAtTime[key], slot)
		}
	}

	for _, slotsAtTime := range teacherSlotsAtTime {
		if len(slotsAtTime) > 1 {
			// Check if all slots have the same subject (allowed for associated group / gruppo linguistico)
			sameSubject := true
			for i := 1; i < len(slotsAtTime); i++ {
				if slotsAtTime[i].SubjectID != slotsAtTime[0].SubjectID {
					sameSubject = false
					break
				}
			}
			if !sameSubject {
				hardConflicts = append(hardConflicts, HardConflict{
					Type: "teacher_double_booking",
					Description: fmt.Sprintf("Docente %s presente in più classi contemporaneamente in materie diverse (%s e %s)",
						slotsAtTime[0].TeacherName, slotsAtTime[0].SubjectName, slotsAtTime[1].SubjectName),
					TeacherID: *slotsAtTime[0].TeacherID,
					DayOfWeek: slotsAtTime[0].DayOfWeek,
					HourIndex: slotsAtTime[0].HourIndex,
				})
			}
		}
	}

	// Record any soft violations for teacher preferences
	for _, s := range generatedSlots {
		if s.TeacherID == nil {
			continue
		}
		tUID := ""
		if s.TeacherUserID != nil {
			tUID = *s.TeacherUserID
		}
		p := getTeacherPref(teacherPrefMap, *s.TeacherID, tUID, s.DayOfWeek, s.HourIndex)
		if p == PrefUnavailable {
			softViolations = append(softViolations, SoftViolation{
				ConstraintType: "teacher_unavailable",
				TeacherID:      s.TeacherID,
				ClassID:        &s.ClassID,
				Description:    fmt.Sprintf("Docente %s assegnato in orario non desiderato (%s %dª ora per %s)", s.TeacherName, dayName(s.DayOfWeek), s.HourIndex, s.ClassName),
				Penalty:        35.0,
			})
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
	if id == "" {
		return
	}
	if m[id] == nil {
		m[id] = make(map[int]map[int]bool)
	}
	if m[id][day] == nil {
		m[id][day] = make(map[int]bool)
	}
	m[id][day][hour] = true
}

func unsetBusy(m map[string]map[int]map[int]bool, id string, day, hour int) {
	if id == "" {
		return
	}
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

func setTeacherGroup(m map[string]map[int]map[int]string, teacherID string, day, hour int, groupID string) {
	if teacherID == "" {
		return
	}
	if m[teacherID] == nil {
		m[teacherID] = make(map[int]map[int]string)
	}
	if m[teacherID][day] == nil {
		m[teacherID][day] = make(map[int]string)
	}
	m[teacherID][day][hour] = groupID
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

func evalPedagogicalScore(s GeneratedSlot) float64 {
	score := 0.0
	if isHeavySubject(s.SubjectName) {
		if s.HourIndex <= 3 {
			score += 10.0 // Preferred morning hours
		} else if s.HourIndex >= 6 {
			score -= 10.0 // Avoid late hours
		}
	}
	return score
}

func getTeacherPref(prefMap map[string]map[int]map[int]string, tID, tUID string, day, hour int) string {
	if tID != "" && prefMap[tID] != nil && prefMap[tID][day] != nil {
		if p, ok := prefMap[tID][day][hour]; ok && p != "" {
			return p
		}
	}
	if tUID != "" && prefMap[tUID] != nil && prefMap[tUID][day] != nil {
		if p, ok := prefMap[tUID][day][hour]; ok && p != "" {
			return p
		}
	}
	return PrefNeutral
}

func getTeacherDayOffHours(dayOffMap map[string]map[int]int, tID, tUID string, day int) int {
	if tID != "" && dayOffMap[tID] != nil {
		if cnt, ok := dayOffMap[tID][day]; ok {
			return cnt
		}
	}
	if tUID != "" && dayOffMap[tUID] != nil {
		if cnt, ok := dayOffMap[tUID][day]; ok {
			return cnt
		}
	}
	return 0
}

func evalTotalSlotScore(s GeneratedSlot, prefMap map[string]map[int]map[int]string, dayOffMap map[string]map[int]int) float64 {
	score := evalPedagogicalScore(s)
	if s.TeacherID != nil {
		tUID := ""
		if s.TeacherUserID != nil {
			tUID = *s.TeacherUserID
		}
		p := getTeacherPref(prefMap, *s.TeacherID, tUID, s.DayOfWeek, s.HourIndex)
		switch p {
		case PrefPreferred:
			score += 15.0
		case PrefUnavailable:
			score -= 35.0
		}
		if getTeacherDayOffHours(dayOffMap, *s.TeacherID, tUID, s.DayOfWeek) >= 4 {
			score -= 40.0
		}
	}
	return score
}

func dayName(d int) string {
	switch d {
	case 1:
		return "Lunedì"
	case 2:
		return "Martedì"
	case 3:
		return "Mercoledì"
	case 4:
		return "Giovedì"
	case 5:
		return "Venerdì"
	case 6:
		return "Sabato"
	default:
		return fmt.Sprintf("Giorno %d", d)
	}
}
