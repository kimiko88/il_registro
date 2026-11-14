package it.scuola.registro.teacher

import it.scuola.registro.teacher.security.BiometricAuthCallback
import it.scuola.registro.teacher.security.BiometricAuthHelper
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Test

class BiometricAuthHelperSecurityTest {
    @Test
    fun authenticationFailsClosedUntilPlatformPromptIsImplemented() {
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

        assertFalse(succeeded)
        assertEquals(BiometricAuthHelper.BIOMETRIC_UNAVAILABLE, errorCode)
    }
}
