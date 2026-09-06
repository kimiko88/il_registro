package it.scuola.registro.secretary

import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.viewmodel.SecretaryViewModel
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class SecretaryViewModelTest {

    private lateinit var viewModel: SecretaryViewModel

    @Before
    fun setUp() {
        viewModel = SecretaryViewModel()
        viewModel.loadSampleData()
    }

    @Test
    fun getUsersByRole_filtersAccurately() {
        val teachers = viewModel.getUsersByRole("teacher")
        assertEquals(2, teachers.size)

        val students = viewModel.getUsersByRole("student")
        assertEquals(1, students.size)
    }

    @Test
    fun addUser_preventsDuplicateEmail() {
        val duplicate = ManagedUser("u5", "Test", "User", "maria.rossi@scuola.it", "teacher")
        val result = viewModel.addUser(duplicate)
        assertFalse(result)

        val newValid = ManagedUser("u6", "Nuovo", "Docente", "nuovo.docente@scuola.it", "teacher")
        val success = viewModel.addUser(newValid)
        assertTrue(success)
        assertEquals(4, viewModel.usersList.size)
    }

    @Test
    fun toggleClassScrutinyLock_invertsLockStatus() {
        val initialState = viewModel.scrutinyClasses.first { it.classId == "c2" }.isLocked
        assertFalse(initialState)

        val locked = viewModel.toggleClassScrutinyLock("c2")
        assertTrue(locked)
        assertTrue(viewModel.scrutinyClasses.first { it.classId == "c2" }.isLocked)
    }

    @Test
    fun generateCertificate_createsCompletedRequestWithPdfUrl() {
        val cert = viewModel.generateCertificate("u3", "iscrizione")
        assertNotNull(cert)
        assertEquals("completato", cert.status)
        assertTrue(cert.generatedPdfUrl.endsWith(".pdf"))
    }
}
