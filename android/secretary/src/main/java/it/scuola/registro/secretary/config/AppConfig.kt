package it.scuola.registro.secretary.config

object AppConfig {
    var BASE_URL: String = "http://10.0.2.2:8080/api/v1"
    var WS_URL: String = "ws://10.0.2.2:8080/api/v1/ws"
    const val TIMEOUT_SECONDS: Long = 30
}
