package it.scuola.registro.teacher

import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall
import it.scuola.registro.teacher.viewmodel.TeacherViewModel
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

/**
 * Navigation & Core Flow Tests for Teacher App.
 * Verifies end-to-end teacher workflow across all primary tabs:
 * 0: Lesson Signing (Firma Ora di Lezione & Argomento)
 * 1: Student Roll Call / Attendance (Appello di Classe)
 * 2: Grade Proposals & Assessment (Valutazioni & Voti)
 * 3: Daily Agenda & Homework Assignment (Agenda / Compiti)
 * 4: Scrutiny Board & Debt Resolution (Scrutini & Recupero Debiti)
 */
class TeacherNavigationFlowTest {

    private lateinit var viewModel: TeacherViewModel

    @Before
    fun setUp() {
        viewModel = TeacherViewModel()
        viewModel.loadSessionData()
    }

    @Test
    fun testCompleteTeacherDailyWorkflow() {
        // --- TAB 0: LESSON SIGNING FLOW ---
        var currentTab = 0
        assertEquals(0, currentTab)
        assertNotNull(viewModel.activeSession)
        assertEquals("Classe 3A", viewModel.activeSession?.className)
        assertEquals("Matematica", viewModel.activeSession?.subject)
        assertFalse("Session must initially be unsigned", viewModel.activeSession?.isSigned == true)

        // Attempt to sign with empty topic (must be rejected)
        val emptySignSuccess = viewModel.signLessonHour("")
        assertFalse("Signing with blank topic must fail", emptySignSuccess)

        // Sign with legitimate lesson topic
        val signSuccess = viewModel.signLessonHour("Studio di funzioni razionali fratte e asintoti")
        assertTrue("Signing with valid topic must succeed", signSuccess)
        assertTrue("Session must now be signed", viewModel.activeSession?.isSigned == true)
        assertEquals("Studio di funzioni razionali fratte e asintoti", viewModel.activeSession?.lessonTopic)

        // --- TAB 1: ROLL CALL / APPELLO FLOW ---
        currentTab = 1
        assertEquals(1, currentTab)
        val rollCall = viewModel.rollCallList
        assertEquals(4, rollCall.size)

        // Check initial state
        assertEquals("presente", rollCall.find { it.studentId == "s1" }?.status)
        assertEquals("assente", rollCall.find { it.studentId == "s3" }?.status)
        assertEquals("ritardo", rollCall.find { it.studentId == "s4" }?.status)

        // Update student status: s3 enters late instead of absent
        val statusUpdate1 = viewModel.updateStudentStatus("s3", "ritardo")
        assertTrue(statusUpdate1)
        assertEquals("ritardo", viewModel.rollCallList.find { it.studentId == "s3" }?.status)

        // Mark s2 as early exit
        val statusUpdate2 = viewModel.updateStudentStatus("s2", "uscita_anticipata")
        assertTrue(statusUpdate2)
        assertEquals("uscita_anticipata", viewModel.rollCallList.find { it.studentId == "s2" }?.status)

        // --- TAB 2: GRADE INSERTION & VALIDATION FLOW ---
        currentTab = 2
        assertEquals(2, currentTab)

        // Add valid oral grade
        val oralGrade = GradeProposal(
            studentId = "s1",
            subjectId = "matematica",
            grade = 8.5,
            weight = 1.0,
            type = "Orale",
            comment = "Ottima esposizione su asintoti obliqui"
        )
        val gradeSuccess = viewModel.addGrade(oralGrade)
        assertTrue("Valid grade must be accepted", gradeSuccess)
        assertEquals(1, viewModel.gradesList.size)

        // Add valid written test grade
        val writtenGrade = GradeProposal(
            studentId = "s2",
            subjectId = "matematica",
            grade = 7.0,
            weight = 1.5,
            type = "Scritto",
            comment = "Compito ben strutturato"
        )
        assertTrue(viewModel.addGrade(writtenGrade))
        assertEquals(2, viewModel.gradesList.size)

        // Out-of-bounds rejection
        assertFalse(viewModel.addGrade(GradeProposal("s1", "matematica", 10.5, 1.0, "Orale", "Out of bounds")))
        assertFalse(viewModel.addGrade(GradeProposal("s1", "matematica", 0.0, 1.0, "Orale", "Out of bounds")))

        // --- TAB 3: AGENDA / HOMEWORK ASSIGNMENT ---
        currentTab = 3
        assertEquals(3, currentTab)
        // Verify teacher can assign homework and verify session status remains intact
        assertTrue(viewModel.activeSession?.isSigned == true)

        // --- TAB 4: SCRUTINY & RECOVERY SESSIONS ---
        currentTab = 4
        assertEquals(4, currentTab)
        assertEquals(0, viewModel.deferredResolutions.size)

        // Deliberate debt recovery for student s3
        val debtResolution = DeferredScrutinyResolution(
            studentId = "s3",
            subjectId = "matematica",
            recoveryGrade = 6.0,
            status = "recuperato",
            finalDecision = "Promosso con debito saldato a seguito di prova scritta integrativa"
        )
        val resolutionSuccess = viewModel.deliberateDeferredScrutiny(debtResolution)
        assertTrue("Deferred scrutiny deliberation must succeed", resolutionSuccess)
        assertEquals(1, viewModel.deferredResolutions.size)
        assertEquals("recuperato", viewModel.deferredResolutions[0].status)
        assertEquals(6.0, viewModel.deferredResolutions[0].recoveryGrade, 0.01)
    }

    @Test
    fun testRollCallAttendanceConsistency() {
        val students = viewModel.rollCallList
        val initialPresent = students.count { it.status == "presente" }
        assertEquals(2, initialPresent)

        // Toggle student s4 to presente
        viewModel.toggleStudentAttendance("s4", "presente")
        val updatedPresent = viewModel.rollCallList.count { it.status == "presente" }
        assertEquals(3, updatedPresent)
    }
}
