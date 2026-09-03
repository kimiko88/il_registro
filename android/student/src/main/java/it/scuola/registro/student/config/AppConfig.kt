package it.scuola.registro.student.config

object AppConfig {
    /**
     * CONFIGURAZIONE URL SERVER API
     *
     * 1. Per l'EMULATORE ANDROID (Default):
     *    Usa "http://10.0.2.2:8080/api/v1" (10.0.2.2 è l'alias dell'host localhost sulla macchina di sviluppo).
     *
     * 2. Per DISPOSITIVO FISICO ANDROID collegato via Wi-Fi:
     *    Usa l'indirizzo IP locale del tuo computer, ad esempio:
     *    "http://192.168.1.100:8080/api/v1"
     *
     * 3. Per AMBIENTE DI PRODUZIONE / SERVER REMOTO:
     *    Usa il dominio pubblico protetto da HTTPS, ad esempio:
     *    "https://registro.tuascuola.it/api/v1"
     */
    var BASE_URL: String = "https://registro-backend-fdu2.onrender.com/api/v1"

    /**
     * URL WEBSOCKET PER AGGIORNAMENTI LIVE & NOTIFICHE IN TEMPO REALE
     */
    var WS_URL: String = "wss://registro-backend-fdu2.onrender.com/api/v1/ws"
    
    const val TIMEOUT_SECONDS: Long = 30
}
