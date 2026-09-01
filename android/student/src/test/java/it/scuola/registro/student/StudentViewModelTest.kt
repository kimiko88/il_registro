package it.scuola.registro.student

import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.viewmodel.StudentViewModel
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

class StudentViewModelTest {

    private lateinit var viewModel: StudentViewModel

    @Before
    fun setUp() {
        viewModel = StudentViewModel()
        viewModel.loadSampleData()
    }

    @Test
    fun calculateGPA_correctlyAveragesGrades() {
        val gpa = viewModel.calculateGPA()
        // (8.5 + 7.0 + 8.0 + 9.0 + 7.0) / 5 = 39.5 / 5 = 7.9
        assertEquals(7.9, gpa, 0.01)
    }

    @Test
    fun getGradesForSubject_filtersCorrectly() {
        val mathGrades = viewModel.getGradesForSubject("Matematica")
        assertEquals(2, mathGrades.size)
        assertEquals("Matematica", mathGrades[0].subject)
    }

    @Test
    fun toggleHomeworkCompletion_updatesState() {
        val newState = viewModel.toggleHomeworkCompletion("1")
        assertTrue(newState)
    }

    @Test
    fun attendanceCounts_returnCorrectTotals() {
        assertEquals(2, viewModel.getTotalAbsences())
        assertEquals(1, viewModel.getTotalLates())
    }
}
