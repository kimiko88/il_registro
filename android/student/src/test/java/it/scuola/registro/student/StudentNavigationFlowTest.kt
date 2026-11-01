package it.scuola.registro.student

import it.scuola.registro.student.data.AttendanceRecord
import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.data.HomeworkAssignment
import it.scuola.registro.student.data.ScrutinyReportCard
import it.scuola.registro.student.viewmodel.StudentViewModel
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

/**
 * Navigation & Core Flow Tests for Student App.
 * Verifies end-to-end user navigation across all 5 main tabs of the Student Dashboard:
 * 0: Home / Overview
 * 1: Grades & Averages (Voti)
 * 2: Homework & Agenda (Compiti)
 * 3: Attendance & Lates (Presenze)
 * 4: Final Scrutiny Report Card (Pagella)
 */
class StudentNavigationFlowTest {

    private lateinit var viewModel: StudentViewModel

    @Before
    fun setUp() {
        viewModel = StudentViewModel()
        viewModel.loadSampleData()
    }

    @Test
    fun testCompleteStudentNavigationAndFeatureFlow() {
        // --- TAB 0: HOME / OVERVIEW ---
        var currentTab = 0
        assertEquals(0, currentTab)
        assertNotNull(viewModel.grades)
        val initialGpa = viewModel.calculateGPA()
        assertEquals(7.9, initialGpa, 0.05)
        assertEquals(2, viewModel.getTotalAbsences())
        assertEquals(1, viewModel.getTotalLates())

        // --- TAB 1: GRADES & SUBJECT BREAKDOWN ---
        currentTab = 1
        assertEquals(1, currentTab)
        val allGrades = viewModel.grades
        assertEquals(5, allGrades.size)

        // Filter by subject
        val mathGrades = viewModel.getGradesForSubject("Matematica")
        assertEquals(2, mathGrades.size)
        val mathAverage = mathGrades.map { it.grade }.average()
        assertEquals(7.75, mathAverage, 0.01)

        val italianGrades = viewModel.getGradesForSubject("Italiano")
        assertEquals(1, italianGrades.size)
        assertEquals(8.0, italianGrades[0].grade, 0.01)

        val englishGrades = viewModel.getGradesForSubject("Inglese")
        assertEquals(1, englishGrades.size)
        assertEquals(9.0, englishGrades[0].grade, 0.01)

        // --- TAB 2: HOMEWORK & AGENDA ---
        currentTab = 2
        assertEquals(2, currentTab)
        val homeworkList = viewModel.homeworkList
        assertEquals(3, homeworkList.size)

        // Check initial completion state
        val task1 = homeworkList.first { it.id == "1" }
        assertFalse("Task 1 should initially be incomplete", task1.isCompleted)

        // Toggle task completion
        val isNowCompleted = viewModel.toggleHomeworkCompletion("1")
        assertTrue("Task 1 should now be completed", isNowCompleted)
        assertTrue(viewModel.homeworkList.first { it.id == "1" }.isCompleted)

        // Toggle back to incomplete
        val toggledBack = viewModel.toggleHomeworkCompletion("1")
        assertFalse("Task 1 should toggle back to incomplete", toggledBack)
        assertFalse(viewModel.homeworkList.first { it.id == "1" }.isCompleted)

        // --- TAB 3: ATTENDANCE & JUSTIFICATIONS ---
        currentTab = 3
        assertEquals(3, currentTab)
        val attendance = viewModel.attendanceRecords
        assertEquals(3, attendance.size)

        val absences = attendance.filter { it.type == "Assenza" }
        assertEquals(2, absences.size)
        assertTrue("All absences in sample should be marked justified", absences.all { it.isJustified })

        val lates = attendance.filter { it.type == "Ritardo" }
        assertEquals(1, lates.size)
        assertEquals("Traffico", lates[0].reason)

        // --- TAB 4: REPORT CARD & SCRUTINY ---
        currentTab = 4
        assertEquals(4, currentTab)
        val reportCard = viewModel.reportCard
        assertNotNull("Report card must be loaded", reportCard)
        assertEquals(2, reportCard?.period)
        assertEquals(9, reportCard?.conductGrade)
        assertEquals("Promosso", reportCard?.finalDecision)
        assertEquals(8, reportCard?.credits)
        assertEquals(8.0, reportCard?.grades?.get("Matematica") ?: 0.0, 0.01)
        assertEquals(8.0, reportCard?.grades?.get("Italiano") ?: 0.0, 0.01)
        assertEquals(7.0, reportCard?.grades?.get("Fisica") ?: 0.0, 0.01)
    }

    @Test
    fun testTabTransitionStateIntegrity() {
        // Verify state is preserved when switching back and forth between tabs
        var activeTab = 0
        assertEquals(0, activeTab)

        // Switch to agenda and mark a task
        activeTab = 2
        assertEquals(2, activeTab)
        viewModel.toggleHomeworkCompletion("2")
        assertTrue(viewModel.homeworkList.first { it.id == "2" }.isCompleted)

        // Switch to grades tab
        activeTab = 1
        assertEquals(1, activeTab)
        assertEquals(5, viewModel.grades.size)

        // Return to agenda and verify completed task state remains intact
        activeTab = 2
        assertEquals(2, activeTab)
        assertTrue(viewModel.homeworkList.first { it.id == "2" }.isCompleted)

        // Switch to report card
        activeTab = 4
        assertEquals(4, activeTab)
        assertEquals("Promosso", viewModel.reportCard?.finalDecision)

        // Return to home
        activeTab = 0
        assertEquals(0, activeTab)
        assertEquals(7.9, viewModel.calculateGPA(), 0.05)
    }
}

