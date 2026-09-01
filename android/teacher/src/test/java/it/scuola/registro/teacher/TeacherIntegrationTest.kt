package it.scuola.registro.teacher

import it.scuola.registro.teacher.network.HttpTeacherApiService
import it.scuola.registro.teacher.viewmodel.TeacherViewModel
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class TeacherIntegrationTest {

    private lateinit var apiService: HttpTeacherApiService
    private lateinit var viewModel: TeacherViewModel

    @Before
    fun setUp() {
        apiService = HttpTeacherApiService("https://api.scuola.registro.it/api/v1")
        viewModel = TeacherViewModel()
    }

    @Test
    fun fullTeacherFlow_signLessonRollCallGradeAndDeferredScrutiny() = runBlocking {
        assertNotNull(apiService)

        // 1. Load initial classroom session
        viewModel.loadSampleSession()
        assertNotNull(viewModel.currentSession)
        assertEquals("Classe 3A", viewModel.currentSession?.className)

        // 2. Sign Lesson
        val signed = viewModel.signCurrentLesson("Studio delle funzioni esponenziali")
        assertTrue(signed)
        assertTrue(viewModel.currentSession?.isSigned == true)

        // 3. Mark Roll Call
        val rollCallSaved = viewModel.toggleStudentAttendance("st1", "Assente")
        assertTrue(rollCallSaved)

        // 4. Submit Grade
        val gradeSubmitted = viewModel.submitGradeForStudent("st1", 8.5, 1.0, "Scritto", "Ottima verifica")
        assertTrue(gradeSubmitted)

        // 5. Save Deferred Scrutiny Resolution
        val scrutinySaved = viewModel.resolveDeferredStudent("st2", 6.0, "Ammesso", "Debito formativo saldato con prova scritta positiva")
        assertTrue(scrutinySaved)
    }
}
