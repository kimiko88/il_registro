package it.scuola.registro.student

import it.scuola.registro.student.cache.OfflineCacheManager
import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.security.BiometricAuthCallback
import it.scuola.registro.student.security.BiometricAuthHelper
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

/**
 * Deterministic integration checks for local student security and offline cache.
 *
 * Live backend credentials and physical biometric hardware are intentionally not
 * prerequisites for JVM unit tests. Those checks belong in a provisioned device E2E job.
 */
class StudentOfflineIntegrationTest {

    @Test
    fun offlineCacheStoresAndClearsGrades() {
        val cacheManager = OfflineCacheManager()
        val grades = listOf(
            GradeEntry(
                id = "grade-1",
                subject = "Matematica",
                grade = 8.5,
                weight = 1.0,
                type = "Scritto",
                date = "2026-11-14",
                period = 1
            ),
            GradeEntry(
                id = "grade-2",
                subject = "Italiano",
                grade = 9.0,
                weight = 1.0,
                type = "Orale",
                date = "2026-11-13",
                period = 1
            )
        )

        assertFalse(cacheManager.isCacheValid())

        cacheManager.saveGrades(grades)

        assertTrue(cacheManager.isCacheValid())
        assertEquals(grades, cacheManager.getCachedGrades())

        cacheManager.clearCache()

        assertFalse(cacheManager.isCacheValid())
        assertTrue(cacheManager.getCachedGrades().isEmpty())
    }

    @Test
    fun biometricAuthenticationFailsClosedWhenPlatformPromptIsUnavailable() {
        val helper = BiometricAuthHelper()
        var succeeded = false
        var capturedErrorCode: Int? = null

        assertFalse(helper.isBiometricAvailable())
        helper.authenticate(object : BiometricAuthCallback {
            override fun onSuccess() {
                succeeded = true
            }

            override fun onError(errorCode: Int, errString: String) {
                capturedErrorCode = errorCode
            }

            override fun onFailed() = Unit
        })

        assertFalse("Unavailable biometrics must never authenticate", succeeded)
        assertEquals(BiometricAuthHelper.BIOMETRIC_UNAVAILABLE, capturedErrorCode)
    }
}
