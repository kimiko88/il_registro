package it.scuola.registro.teacher.security

/**
 * Authentication is intentionally unavailable until a real platform BiometricPrompt
 * is wired to a cryptographically protected session. Never report a biometric success
 * without a system-authenticated prompt.
 */
interface BiometricAuthCallback {
    fun onSuccess()
    fun onError(errorCode: Int, errString: String)
    fun onFailed()
}

class BiometricAuthHelper {
    fun isBiometricAvailable(): Boolean = false

    fun authenticate(callback: BiometricAuthCallback) {
        callback.onError(BIOMETRIC_UNAVAILABLE, "Biometric authentication is not configured")
    }

    companion object {
        const val BIOMETRIC_UNAVAILABLE = -1
    }
}
