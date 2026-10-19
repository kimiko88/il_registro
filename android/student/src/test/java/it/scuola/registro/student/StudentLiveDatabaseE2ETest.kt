package it.scuola.registro.student

import it.scuola.registro.student.cache.OfflineCacheManager
import it.scuola.registro.student.network.HttpStudentApiService
import it.scuola.registro.student.security.BiometricAuthCallback
import it.scuola.registro.student.security.BiometricAuthHelper
import it.scuola.registro.student.security.SecurityCheckHelper
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

/**
 * End-to-End & Integration Test for Student App against the REAL Live Backend and Database.
 * Strictly adheres to the zero-mock / zero-stub rule: all data originates from the live server.
 */
class StudentLiveDatabaseE2ETest {

    private lateinit var apiService: HttpStudentApiService
    private lateinit var cacheManager: OfflineCacheManager
    private lateinit var securityHelper: SecurityCheckHelper
    private lateinit var biometricHelper: BiometricAuthHelper

    @Before
    fun setUp() {
        apiService = HttpStudentApiService("https://registro-backend-fdu2.onrender.com/api/v1")
        cacheManager = OfflineCacheManager()
        securityHelper = SecurityCheckHelper()
        biometricHelper = BiometricAuthHelper()
    }

    @Test
    fun testLiveStudentFullFlow_RealDatabaseIntegration() = runBlocking {
        // 1. Device Security Integrity Checks
        val securityReport = securityHelper.getSecurityReport()
        assertNotNull("Security report must be generated", securityReport)
        assertTrue("Security report must contain isRooted check", securityReport.containsKey("isRooted"))
        assertTrue("Security report must contain isEmulator check", securityReport.containsKey("isEmulator"))

        // 2. Biometric Authentication Hardware / Enrolment verification
        var biometricSuccess = false
        biometricHelper.authenticate(object : BiometricAuthCallback {
            override fun onSuccess() {
                biometricSuccess = true
            }
            override fun onError(errorCode: Int, errString: String) {
                biometricSuccess = false
            }
            override fun onFailed() {
                biometricSuccess = false
            }
        })
        assertTrue("Biometric authentication helper must function", biometricSuccess)

        // 3. Real Server Login with verified live database account
        val loginResult = apiService.login("studente2b_10@scuola.it", "password")
        assertTrue("Real database login must succeed: ${loginResult.exceptionOrNull()?.message}", loginResult.isSuccess)

        val (token, studentUser) = loginResult.getOrThrow()
        assertTrue("Live JWT token must not be empty", token.isNotEmpty())
        assertEquals("studente2b_10@scuola.it", studentUser.email.lowercase())

        // 4. Retrieve REAL live student grades from database
        val gradesResult = apiService.getGrades(token)
        assertTrue("Live grades fetch must succeed: ${gradesResult.exceptionOrNull()?.message}", gradesResult.isSuccess)
        val realGrades = gradesResult.getOrThrow()
        assertTrue("Student must have real recorded grades in database", realGrades.isNotEmpty())

        val firstGrade = realGrades.first()
        assertTrue("Grade value must be greater than 0", firstGrade.grade > 0.0)
        assertTrue("Subject name must be populated", firstGrade.subject.isNotBlank())

        // 5. Store real live grades into OfflineCacheManager and verify resilience
        cacheManager.saveGrades(realGrades)
        assertTrue("Cache must be valid after saving live data", cacheManager.isCacheValid())
        assertEquals("Cached grade count must match live database count", realGrades.size, cacheManager.getCachedGrades().size)

        // 6. Retrieve attendance records from live database
        val attendanceResult = apiService.getAttendance(token)
        assertTrue("Attendance request must complete successfully", attendanceResult.isSuccess)

        // 7. Clear cache upon session termination
        cacheManager.clearCache()
        assertFalse("Cache must be invalidated after clear", cacheManager.isCacheValid())
    }
}
