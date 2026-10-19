package it.scuola.registro.teacher.security

import android.os.Build
import java.io.File

class SecurityCheckHelper {

    private val suBinaryPaths = listOf(
        "/system/app/Superuser.apk",
        "/sbin/su",
        "/system/bin/su",
        "/system/xbin/su",
        "/data/local/xbin/su",
        "/data/local/bin/su",
        "/system/sd/xbin/su",
        "/system/bin/failsafe/su",
        "/data/local/su",
        "/su/bin/su",
        "/system/xbin/daemonsu"
    )

    private val magiskPaths = listOf(
        "/sbin/.magisk",
        "/cache/.disable_magisk",
        "/data/adb/magisk"
    )

    fun isDeviceRooted(): Boolean {
        return checkRootBinaries() || checkTestKeys() || checkSuExecution() || checkMagisk()
    }

    private fun checkRootBinaries(): Boolean {
        return suBinaryPaths.any { path ->
            try {
                File(path).exists()
            } catch (_: Exception) {
                false
            }
        }
    }

    private fun checkMagisk(): Boolean {
        return magiskPaths.any { path ->
            try {
                File(path).exists()
            } catch (_: Exception) {
                false
            }
        }
    }

    private fun checkTestKeys(): Boolean {
        val buildTags = Build.TAGS
        return buildTags != null && buildTags.contains("test-keys")
    }

    private fun checkSuExecution(): Boolean {
        var process: Process? = null
        return try {
            process = Runtime.getRuntime().exec(arrayOf("/system/xbin/which", "su"))
            val exitCode = process.waitFor()
            exitCode == 0
        } catch (_: Exception) {
            false
        } finally {
            process?.destroy()
        }
    }

    fun isEmulator(): Boolean {
        val fingerprint = Build.FINGERPRINT ?: ""
        val model = Build.MODEL ?: ""
        val manufacturer = Build.MANUFACTURER ?: ""
        val brand = Build.BRAND ?: ""
        val device = Build.DEVICE ?: ""
        val product = Build.PRODUCT ?: ""
        val hardware = Build.HARDWARE ?: ""

        return fingerprint.startsWith("generic") ||
                fingerprint.startsWith("unknown") ||
                model.contains("google_sdk") ||
                model.contains("Emulator") ||
                model.contains("Android SDK built for x86") ||
                manufacturer.contains("Genymotion") ||
                hardware.contains("goldfish") ||
                hardware.contains("ranchu") ||
                product.contains("sdk_google") ||
                product.contains("google_sdk") ||
                product.contains("sdk") ||
                brand.startsWith("generic") && device.startsWith("generic")
    }

    fun getSecurityReport(): Map<String, Boolean> {
        return mapOf(
            "isRooted" to isDeviceRooted(),
            "isEmulator" to isEmulator(),
            "hasRootBinaries" to checkRootBinaries(),
            "hasTestKeys" to checkTestKeys(),
            "hasMagisk" to checkMagisk()
        )
    }
}
