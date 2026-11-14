package it.scuola.registro.secretary

import it.scuola.registro.secretary.cache.OfflineCacheManager
import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.security.BiometricAuthCallback
import it.scuola.registro.secretary.security.BiometricAuthHelper
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

/**
 * Deterministic integration checks for the secretary app's local security and offline cache.
 *
 * Live backend credentials and physical biometric hardware are intentionally not prerequisites
 * for JVM unit tests. Those checks belong in a separately provisioned device/backend E2E job.
 */
class SecretaryOfflineIntegrationTest {

    @Test
    fun offlineCacheStoresAndClearsUsers() {
        val cacheManager = OfflineCacheManager()
        val users = listOf(
            ManagedUser(
                id = "user-1",
                firstName = "Ada",
                lastName = "Lovelace",
                email = "ada@example.test",
                role = "teacher"
            ),
            ManagedUser(
                id = "user-2",
                firstName = "Alan",
                lastName = "Turing",
                email = "alan@example.test",
                role = "student"
            )
        )

        assertFalse(cacheManager.isCacheValid())

        cacheManager.saveUsers(users)

        assertTrue(cacheManager.isCacheValid())
        assertEquals(users, cacheManager.getCachedUsers())

        cacheManager.clearCache()

        assertFalse(cacheManager.isCacheValid())
        assertTrue(cacheManager.getCachedUsers().isEmpty())
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
