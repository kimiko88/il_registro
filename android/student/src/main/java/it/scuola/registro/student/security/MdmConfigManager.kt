package it.scuola.registro.student.security

import android.content.Context
import android.content.RestrictionsManager
import android.os.Bundle
import it.scuola.registro.student.config.AppConfig

class MdmConfigManager(private val context: Context?) {

    fun applyMdmPolicies(): Bundle {
        val bundle = Bundle()
        if (context == null) return bundle

        val restrictionsManager = context.getSystemService(Context.RESTRICTIONS_SERVICE) as? RestrictionsManager
        val appRestrictions = restrictionsManager?.applicationRestrictions

        if (appRestrictions != null) {
            if (appRestrictions.containsKey("server_url")) {
                val serverUrl = appRestrictions.getString("server_url")
                if (!serverUrl.isNullOrBlank()) {
                    AppConfig.BASE_URL = serverUrl
                    bundle.putString("server_url", serverUrl)
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
