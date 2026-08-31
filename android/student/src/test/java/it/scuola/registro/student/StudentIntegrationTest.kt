package it.scuola.registro.student

import it.scuola.registro.student.cache.OfflineCacheManager
import it.scuola.registro.student.network.MockStudentApiService
import it.scuola.registro.student.viewmodel.StudentViewModel
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class StudentIntegrationTest {

    private lateinit var apiService: MockStudentApiService
    private lateinit var cacheManager: OfflineCacheManager
    private lateinit var viewModel: StudentViewModel

    @Before
    fun setUp() {
        apiService = MockStudentApiService()
        cacheManager = OfflineCacheManager()
        viewModel = StudentViewModel()
    }

    @Test
    fun fullStudentDataFlow_loginFetchCacheAndCalculateGPA() = runBlocking {
        // 1. Authenticate
        val loginResult = apiService.login("mario.rossi@studenti.it", "password123")
        assertTrue(loginResult.isSuccess)
        val (token, user) = loginResult.getOrThrow()
        assertEquals("Mario", user.firstName)

        // 2. Fetch Grades from API
        val gradesResult = apiService.getGrades(token)
        assertTrue(gradesResult.isSuccess)
        val grades = gradesResult.getOrThrow()
        assertEquals(3, grades.size)

        // 3. Cache Grades Offline
        cacheManager.saveGrades(grades)
        assertTrue(cacheManager.isCacheValid())

        // 4. Retrieve from Cache and verify GPA
        val cachedGrades = cacheManager.getCachedGrades()
        assertEquals(3, cachedGrades.size)

        // Load into ViewModel and verify calculations
        viewModel.loadSampleData()
        val gpa = viewModel.calculateGPA()
        assertEquals(7.9, gpa, 0.01)

        // 5. Toggle homework task
        val homeworkToggled = viewModel.toggleHomeworkCompletion("1")
        assertTrue(homeworkToggled)
    }
}
