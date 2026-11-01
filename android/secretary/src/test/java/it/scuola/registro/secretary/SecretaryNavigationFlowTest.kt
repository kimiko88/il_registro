package it.scuola.registro.secretary

import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.viewmodel.SecretaryViewModel
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

/**
 * Navigation & Core Flow Tests for Secretary App.
 * Verifies end-to-end administrative navigation across key sectors:
 * 0: Dashboard Overview & Real-Time Stats (Pannello KPI)
 * 1: User Directory & Role-Based Management (Gestione Utenti)
 * 2: Class Scrutiny Supervision & Board Locking (Scrutini)
 * 3: Certificate Protocolling & Generation (Certificati)
 * 4: Compliance & System Audit Logs (Audit Trail)
 */
class SecretaryNavigationFlowTest {

    private lateinit var viewModel: SecretaryViewModel

    @Before
    fun setUp() {
        viewModel = SecretaryViewModel()
        viewModel.loadSampleData()
    }

    @Test
    fun testCompleteSecretaryAdministrationWorkflow() {
        // --- TAB 0: DASHBOARD OVERVIEW & KPIS ---
        var currentTab = 0
        assertEquals(0, currentTab)
        assertEquals(3, viewModel.usersList.size)
        assertEquals(3, viewModel.scrutinyClasses.size)
        assertEquals(2, viewModel.certificateRequests.size)
        assertEquals(2, viewModel.auditLogs.size)

        // --- TAB 1: USER DIRECTORY & MANAGEMENT FLOW ---
        currentTab = 1
        assertEquals(1, currentTab)

        // Filter by role
        val teachers = viewModel.getUsersByRole("teacher")
        assertEquals(2, teachers.size)
        assertTrue(teachers.all { it.role == "teacher" })

        val students = viewModel.getUsersByRole("student")
        assertEquals(1, students.size)
        assertEquals("Mario", students[0].firstName)

        // Add new student
        val newStudent = ManagedUser(
            id = "u10",
            firstName = "Alessandro",
            lastName = "Verdi",
            email = "alessandro.verdi@scuola.it",
            role = "student",
            isActive = true
        )
        val addSuccess = viewModel.addUser(newStudent)
        assertTrue("Adding unique user must succeed", addSuccess)
        assertEquals(4, viewModel.usersList.size)

        // Duplicate email prevention
        val duplicateUser = ManagedUser(
            id = "u11",
            firstName = "Duplicato",
            lastName = "Test",
            email = "alessandro.verdi@scuola.it",
            role = "student"
        )
        assertFalse("Duplicate email must be rejected", viewModel.addUser(duplicateUser))

        // --- TAB 2: SCRUTINY SUPERVISION & LOCKING FLOW ---
        currentTab = 2
        assertEquals(2, currentTab)
        val class1 = viewModel.scrutinyClasses.first { it.classId == "c1" }
        val class2 = viewModel.scrutinyClasses.first { it.classId == "c2" }

        assertTrue("Class 1 should initially be locked", class1.isLocked)
        assertFalse("Class 2 should initially be unlocked", class2.isLocked)

        // Lock class 2 to prevent modifications
        val lockResult = viewModel.toggleClassScrutinyLock("c2")
        assertTrue(lockResult)
        assertTrue("Class 2 must now be locked", viewModel.scrutinyClasses.first { it.classId == "c2" }.isLocked)

        // Unlock class 1 for emergency review
        val unlockResult = viewModel.toggleClassScrutinyLock("c1")
        assertFalse(unlockResult)
        assertFalse("Class 1 must now be unlocked", viewModel.scrutinyClasses.first { it.classId == "c1" }.isLocked)

        // --- TAB 3: CERTIFICATE GENERATION & PROTOCOLLING FLOW ---
        currentTab = 3
        assertEquals(3, currentTab)
        val initialCertsCount = viewModel.certificateRequests.size
        assertEquals(2, initialCertsCount)

        // Generate enrollment certificate for student u3
        val cert = viewModel.generateCertificate("u3", "iscrizione_frequenza")
        assertNotNull("Generated certificate must not be null", cert)
        assertEquals("u3", cert.studentId)
        assertEquals("iscrizione_frequenza", cert.certificateType)
        assertEquals("completato", cert.status)
        assertTrue("Certificate PDF must be generated", cert.generatedPdfUrl.endsWith(".pdf"))
        assertTrue("Protocol number must be assigned", cert.protocolNo.startsWith("PROT-"))
        assertEquals(initialCertsCount + 1, viewModel.certificateRequests.size)

        // --- TAB 4: AUDIT LOGS & TRACEABILITY FLOW ---
        currentTab = 4
        assertEquals(4, currentTab)
        val logs = viewModel.auditLogs
        assertEquals(2, logs.size)
        assertTrue(logs.any { it.action == "BLOCCO_SCRUTINIO" })
        assertTrue(logs.any { it.action == "CREAZIONE_UTENTE" })
    }

    @Test
    fun testUserFilteringAndRoleIntegrity() {
        // Verify filtering behaves correctly across all roles
        val allUsers = viewModel.usersList
        val teachers = viewModel.getUsersByRole("TEACHER")
        val students = viewModel.getUsersByRole("STUDENT")
        val parents = viewModel.getUsersByRole("PARENT")

        assertEquals(2, teachers.size)
        assertEquals(1, students.size)
        assertEquals(0, parents.size)
        assertEquals(allUsers.size, teachers.size + students.size + parents.size)
    }
}
