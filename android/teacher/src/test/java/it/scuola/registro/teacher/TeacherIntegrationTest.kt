package it.scuola.registro.teacher

import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall
import it.scuola.registro.teacher.network.MockTeacherApiService
import it.scuola.registro.teacher.viewmodel.TeacherViewModel
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class TeacherIntegrationTest {

    private lateinit var apiService: MockTeacherApiService
    private lateinit var viewModel: TeacherViewModel

    @Before
    fun setUp() {
        apiService = MockTeacherApiService()
        viewModel = TeacherViewModel()
    }

    @Test
    fun fullTeacherFlow_signHourTakeRollCallInsertGradeAndDeferredScrutiny() = runBlocking {
        // 1. Authenticate Teacher
        val loginResult = apiService.login("docente.bianchi@scuola.it", "password123")
        assertTrue(loginResult.isSuccess)
        val token = loginResult.getOrThrow()

        // 2. Sign Lesson Hour via API & ViewModel
        val topic = "Calcolo infinitesimale ed integrali"
        val signResult = apiService.signLesson(token, "3A", topic)
        assertTrue(signResult.isSuccess)
        assertTrue(signResult.getOrThrow().isSigned)

        viewModel.loadSessionData()
        val vmSign = viewModel.signLessonHour(topic)
        assertTrue(vmSign)

        // 3. Roll Call Submission
        val rollCall = listOf(
            StudentRollCall("s1", "Banchi Andrea", "presente"),
            StudentRollCall("s2", "Ferrari Matteo", "assente")
        )
        val rollCallResult = apiService.submitRollCall(token, "3A", rollCall)
        assertTrue(rollCallResult.isSuccess)

        // 4. Grade Input
        val gradeProposal = GradeProposal("s1", "sub1", 8.5, 1.0, "Scritto", "Ottima prova")
        val gradeApiResult = apiService.submitGrade(token, gradeProposal)
        assertTrue(gradeApiResult.isSuccess)

        val gradeVmResult = viewModel.addGrade(gradeProposal)
        assertTrue(gradeVmResult)

        // 5. Deferred Scrutiny Deliberation
        val deferredResolution = DeferredScrutinyResolution("s2", "sub1", 7.0, "recuperato", "promosso_con_debiti_saldati")
        val defApiResult = apiService.saveDeferredScrutiny(token, deferredResolution)
        assertTrue(defApiResult.isSuccess)

        val defVmResult = viewModel.deliberateDeferredScrutiny(deferredResolution)
        assertTrue(defVmResult)
        assertEquals(1, viewModel.deferredResolutions.size)
    }
}
