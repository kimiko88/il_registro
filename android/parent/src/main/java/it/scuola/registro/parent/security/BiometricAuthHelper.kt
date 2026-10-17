package it.scuola.registro.parent.security

interface BiometricAuthCallback {
    fun onSuccess()
    fun onError(errorCode: Int, errString: String)
    fun onFailed()
}

class BiometricAuthHelper {
    fun isBiometricAvailable(): Boolean {
        return true
    }

    fun authenticate(callback: BiometricAuthCallback) {
        if (isBiometricAvailable()) {
            callback.onSuccess()
        } else {
            callback.onError(-1, "Biometric authentication not supported")
        }
    }
}
