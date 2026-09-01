package it.scuola.registro.teacher

import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.viewmodel.TeacherViewModel
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class TeacherViewModelTest {

    private lateinit var viewModel: TeacherViewModel

    @Before
    fun setUp() {
        viewModel = TeacherViewModel()
        viewModel.loadSessionData()
    }

    @Test
    fun signLessonHour_signsSuccessfullyWithTopic() {
        val result = viewModel.signLessonHour("Equazioni esponenziali e logaritmiche")
        assertTrue(result)
        assertTrue(viewModel.activeSession?.isSigned == true)
        assertEquals("Equazioni esponenziali e logaritmiche", viewModel.activeSession?.lessonTopic)
    }

    @Test
    fun signLessonHour_failsWithEmptyTopic() {
        val result = viewModel.signLessonHour("")
        assertFalse(result)
    }

    @Test
    fun updateStudentStatus_updatesAttendanceCorrectly() {
        val updated = viewModel.updateStudentStatus("s3", "presente")
        assertTrue(updated)
        assertEquals("presente", viewModel.rollCallList.find { it.studentId == "s3" }?.status)
    }

    @Test
    fun addGrade_validatesRange() {
        val valid = viewModel.addGrade(GradeProposal("s1", "sub1", 8.5, 1.0, "Scritto", "Ottimo"))
        assertTrue(valid)

        val invalidHigh = viewModel.addGrade(GradeProposal("s1", "sub1", 11.0, 1.0, "Scritto", "Errore"))
        assertFalse(invalidHigh)

        val invalidLow = viewModel.addGrade(GradeProposal("s1", "sub1", 0.5, 1.0, "Scritto", "Errore"))
        assertFalse(invalidLow)
    }

    @Test
    fun deliberateDeferredScrutiny_recordsResolution() {
        val res = DeferredScrutinyResolution("s3", "sub1", 7.0, "recuperato", "promosso_con_debiti_saldati")
        val success = viewModel.deliberateDeferredScrutiny(res)
        assertTrue(success)
        assertEquals(1, viewModel.deferredResolutions.size)
    }
}
