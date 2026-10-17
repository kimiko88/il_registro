package it.scuola.registro.teacher

import it.scuola.registro.teacher.cache.OfflineCacheManager
import it.scuola.registro.teacher.network.HttpTeacherApiService
import it.scuola.registro.teacher.security.BiometricAuthCallback
import it.scuola.registro.teacher.security.BiometricAuthHelper
import it.scuola.registro.teacher.security.SecurityCheckHelper
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

/**
 * End-to-End & Integration Test for Teacher App against the REAL Live Backend and Database.
 * Zero mocks / Zero stubs: queries live classes and students directly from the school server.
 */
class TeacherLiveDatabaseE2ETest {

    private lateinit var apiService: HttpTeacherApiService
    private lateinit var cacheManager: OfflineCacheManager
    private lateinit var securityHelper: SecurityCheckHelper
    private lateinit var biometricHelper: BiometricAuthHelper

    @Before
    fun setUp() {
        apiService = HttpTeacherApiService("https://registro-backend-fdu2.onrender.com/api/v1")
        cacheManager = OfflineCacheManager()
        securityHelper = SecurityCheckHelper()
        biometricHelper = BiometricAuthHelper()
    }

    @Test
    fun testLiveTeacherFullFlow_RealDatabaseIntegration() = runBlocking {
        // 1. Device Security Integrity Checks
        val securityReport = securityHelper.getSecurityReport()
        assertNotNull(securityReport)
        assertTrue(securityReport.containsKey("isRooted"))
        assertTrue(securityReport.containsKey("isEmulator"))

        // 2. Biometric Authentication Helper
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

        // 3. Real Server Login with verified live teacher account
        val loginResult = apiService.login("docente.incarichi@scuola.it", "password")
        assertTrue("Real database login must succeed: ${loginResult.exceptionOrNull()?.message}", loginResult.isSuccess)

        val token = loginResult.getOrThrow()
        assertTrue("JWT token must be populated", token.isNotEmpty())

        // 4. Retrieve REAL live classes from database
        val classesResult = apiService.getClasses(token)
        assertTrue("Classes fetch must succeed: ${classesResult.exceptionOrNull()?.message}", classesResult.isSuccess)
        val classes = classesResult.getOrThrow()
        assertTrue("Teacher must have classes assigned in the real database", classes.isNotEmpty())

        val firstClass = classes.first()
        assertTrue("Class ID must be valid", firstClass.classId.isNotBlank())
        assertTrue("Class name must be valid", firstClass.className.isNotBlank())

        // 5. Retrieve REAL students of that class from the database
        val studentsResult = apiService.getClassStudents(token, firstClass.classId)
        assertTrue("Students fetch must succeed: ${studentsResult.exceptionOrNull()?.message}", studentsResult.isSuccess)
        val students = studentsResult.getOrThrow()
        assertTrue("Class must contain students enrolled in database", students.isNotEmpty())

        val firstStudent = students.first()
        assertTrue("Student ID must be valid", firstStudent.studentId.isNotBlank())
        assertTrue("Student name must be populated", firstStudent.fullName.isNotBlank())

        // 6. Store real classes and roll calls into OfflineCacheManager
        cacheManager.saveClasses(classes)
        cacheManager.saveRollCalls(students)
        assertTrue("Cache must be valid after saving live data", cacheManager.isCacheValid())
        assertEquals("Cached classes count must match", classes.size, cacheManager.getCachedClasses().size)
        assertEquals("Cached students count must match", students.size, cacheManager.getCachedRollCalls().size)

        // 7. Clear cache
        cacheManager.clearCache()
        assertFalse("Cache must be invalidated after clear", cacheManager.isCacheValid())
    }
}
