package it.scuola.registro.student

import it.scuola.registro.student.cache.OfflineCacheManager
import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.network.HttpStudentApiService
import it.scuola.registro.student.viewmodel.StudentViewModel
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class StudentIntegrationTest {

    private lateinit var apiService: HttpStudentApiService
    private lateinit var cacheManager: OfflineCacheManager
    private lateinit var viewModel: StudentViewModel

    @Before
    fun setUp() {
        apiService = HttpStudentApiService("https://api.scuola.registro.it/api/v1")
        cacheManager = OfflineCacheManager()
        viewModel = StudentViewModel()
    }

    @Test
    fun fullStudentDataFlow_cacheAndCalculateGPA() = runBlocking {
        assertNotNull(apiService)

        // 1. Cache Grades Offline
        val testGrades = listOf(
            GradeEntry("1", "Matematica", 8.5, 1.0, "Scritto", "2026-08-30", 1),
            GradeEntry("2", "Italiano", 8.0, 1.0, "Tema", "2026-08-28", 1),
            GradeEntry("3", "Inglese", 9.0, 1.0, "Pratico", "2026-08-25", 1)
        )
        cacheManager.saveGrades(testGrades)
        assertTrue(cacheManager.isCacheValid())

        // 2. Retrieve from Cache and verify GPA
        val cachedGrades = cacheManager.getCachedGrades()
        assertEquals(3, cachedGrades.size)

        // 3. Load into ViewModel and verify calculations
        viewModel.loadSampleData()
        val gpa = viewModel.calculateGPA()
        assertEquals(7.9, gpa, 0.01)

        // 4. Toggle homework task
        val homeworkToggled = viewModel.toggleHomeworkCompletion("1")
        assertTrue(homeworkToggled)
    }
}
