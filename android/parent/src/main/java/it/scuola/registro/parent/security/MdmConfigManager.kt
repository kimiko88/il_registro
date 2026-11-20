package it.scuola.registro.parent.security

import android.content.Context
import android.content.RestrictionsManager
import android.os.Bundle
import java.net.URI
import it.scuola.registro.parent.config.AppConfig

class MdmConfigManager(private val context: Context?) {

    private fun isSecureServerUrl(value: String): Boolean = try {
        val uri = URI(value.trim())
        uri.scheme.equals("https", ignoreCase = true) &&
            !uri.host.isNullOrBlank() &&
            uri.userInfo == null &&
            uri.query == null &&
            uri.fragment == null
    } catch (_: Exception) {
        false
    }

    fun applyMdmPolicies(): Bundle {
        val bundle = Bundle()
        if (context == null) return bundle

        val restrictionsManager = context.getSystemService(Context.RESTRICTIONS_SERVICE) as? RestrictionsManager
        val appRestrictions = restrictionsManager?.applicationRestrictions

        if (appRestrictions != null) {
            if (appRestrictions.containsKey("server_url")) {
                val serverUrl = appRestrictions.getString("server_url")
                if (!serverUrl.isNullOrBlank() && isSecureServerUrl(serverUrl)) {
                    AppConfig.BASE_URL = serverUrl.trim().trimEnd('/')
                    bundle.putString("server_url", AppConfig.BASE_URL)
                } else {
                    bundle.putBoolean("server_url_rejected", true)
                }
            }

            if (appRestrictions.containsKey("kiosk_mode_enabled")) {
                val kiosk = appRestrictions.getBoolean("kiosk_mode_enabled", false)
                bundle.putBoolean("kiosk_mode_enabled", kiosk)
            }
        }
        return bundle
    }
}
