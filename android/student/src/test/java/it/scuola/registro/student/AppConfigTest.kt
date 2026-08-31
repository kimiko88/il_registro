package it.scuola.registro.student

import it.scuola.registro.student.config.AppConfig
import org.junit.Assert.*
import org.junit.Test
import java.net.URI

class AppConfigTest {

    @Test
    fun appConfigBaseUrl_isValidUriAndHttpScheme() {
        val uri = URI.create(AppConfig.BASE_URL)
        assertNotNull(uri)
        assertTrue(uri.scheme == "http" || uri.scheme == "https")
        assertTrue(AppConfig.BASE_URL.endsWith("/api/v1"))
    }

    @Test
    fun appConfigWsUrl_isValidUriAndWsScheme() {
        val uri = URI.create(AppConfig.WS_URL)
        assertNotNull(uri)
        assertTrue(uri.scheme == "ws" || uri.scheme == "wss")
        assertTrue(AppConfig.WS_URL.endsWith("/ws"))
    }

    @Test
    fun timeoutSeconds_isReasonable() {
        assertTrue(AppConfig.TIMEOUT_SECONDS in 10..120)
    }
}
