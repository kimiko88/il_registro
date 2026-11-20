package it.scuola.registro.parent

import it.scuola.registro.parent.cache.OfflineCacheManager
import it.scuola.registro.parent.data.ParentChild
import it.scuola.registro.parent.security.BiometricAuthCallback
import it.scuola.registro.parent.security.BiometricAuthHelper
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

/**
 * Deterministic integration checks for the parent app's local security and offline cache.
 *
 * Live backend credentials and device biometric hardware are intentionally not prerequisites
 * for JVM unit tests. Those checks belong in an explicitly provisioned device/backend E2E job.
 */
class ParentOfflineIntegrationTest {

    @Test
    fun offlineCacheStoresAndClearsChildren() {
        val cacheManager = OfflineCacheManager()
        val children = listOf(
            ParentChild("child-1", "Marco", "Rossi", "2A"),
            ParentChild("child-2", "Giulia", "Bianchi", "3B")
        )

        assertFalse(cacheManager.isCacheValid())

        cacheManager.saveChildren(children)

        assertTrue(cacheManager.isCacheValid())
        assertEquals(children, cacheManager.getCachedChildren())

        cacheManager.clearCache()

        assertFalse(cacheManager.isCacheValid())
        assertTrue(cacheManager.getCachedChildren().isEmpty())
    }

    @Test
    fun biometricAuthenticationFailsClosedWhenPlatformPromptIsUnavailable() {
        val helper = BiometricAuthHelper()
        var succeeded = false
        var errorCode: Int? = null

        assertFalse(helper.isBiometricAvailable())
        helper.authenticate(object : BiometricAuthCallback {
            override fun onSuccess() {
                succeeded = true
            }

            override fun onError(errorCodeValue: Int, errString: String) {
                errorCode = errorCodeValue
            }

            override fun onFailed() = Unit
        })

        assertFalse("Unavailable biometrics must never authenticate", succeeded)
        assertEquals(BiometricAuthHelper.BIOMETRIC_UNAVAILABLE, errorCode)
    }
}
