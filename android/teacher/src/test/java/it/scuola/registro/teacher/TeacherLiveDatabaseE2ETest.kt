package it.scuola.registro.teacher

import it.scuola.registro.teacher.cache.OfflineCacheManager
import it.scuola.registro.teacher.data.ClassSession
import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall
import it.scuola.registro.teacher.security.BiometricAuthCallback
import it.scuola.registro.teacher.security.BiometricAuthHelper
import org.junit.Assert.*
import org.junit.Test

/**
 * Deterministic integration checks for the teacher offline-data and authentication
 * boundaries. Tests must not depend on a deployed school backend, shared credentials,
 * enrolled biometric hardware, or mutable production data.
 */
class TeacherOfflineIntegrationTest {

    @Test
    fun testTeacherOfflineCacheFullFlow() {
        val cacheManager = OfflineCacheManager()
        val classes = listOf(
            ClassSession(
                classId = "class-1",
                className = "3A",
                subject = "Matematica",
                hourSlot = "08:00-09:00"
            )
        )
        val rollCalls = listOf(
            StudentRollCall(studentId = "student-1", fullName = "Studente Test")
        )
        val grades = listOf(
            GradeProposal(studentId = "student-1", subjectId = "math", grade = 8.0)
        )
        val deferred = listOf(
            DeferredScrutinyResolution(
                studentId = "student-1",
                subjectId = "math",
                recoveryGrade = 7.0,
                status = "recuperato",
                finalDecision = "ammesso"
            )
        )

        assertFalse("A new cache must not be considered valid", cacheManager.isCacheValid())

        cacheManager.saveClasses(classes)
        cacheManager.saveRollCalls(rollCalls)
        cacheManager.saveGrades(grades)
        cacheManager.saveDeferredScrutiny(deferred)

        assertTrue("Cache must become valid after classes are saved", cacheManager.isCacheValid())
        assertEquals(classes, cacheManager.getCachedClasses())
        assertEquals(rollCalls, cacheManager.getCachedRollCalls())
        assertEquals(grades, cacheManager.getCachedGrades())
        assertEquals(deferred, cacheManager.getCachedDeferredScrutiny())

        cacheManager.clearCache()

        assertFalse("Cleared cache must be invalid", cacheManager.isCacheValid())
        assertTrue(cacheManager.getCachedClasses().isEmpty())
        assertTrue(cacheManager.getCachedRollCalls().isEmpty())
        assertTrue(cacheManager.getCachedGrades().isEmpty())
        assertTrue(cacheManager.getCachedDeferredScrutiny().isEmpty())
    }

    @Test
    fun testBiometricAuthenticationFailsClosedWhenNotConfigured() {
        val biometricHelper = BiometricAuthHelper()
        var successCalled = false
        var failedCalled = false
        var reportedError: Pair<Int, String>? = null

        assertFalse("Biometrics must remain unavailable until platform authentication is wired",
            biometricHelper.isBiometricAvailable())

        biometricHelper.authenticate(object : BiometricAuthCallback {
            override fun onSuccess() {
                successCalled = true
            }

            override fun onError(errorCode: Int, errString: String) {
                reportedError = errorCode to errString
            }

            override fun onFailed() {
                failedCalled = true
            }
        })

        assertFalse("An unavailable biometric implementation must never authenticate",
            successCalled)
        assertFalse("Unavailable authentication is an error, not a failed biometric match",
            failedCalled)
        assertNotNull("The caller must be informed that authentication is unavailable",
            reportedError)
        assertEquals(BiometricAuthHelper.BIOMETRIC_UNAVAILABLE, reportedError?.first)
        assertTrue(reportedError?.second?.isNotBlank() == true)
    }
}
