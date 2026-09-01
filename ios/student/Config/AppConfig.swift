import Foundation

public struct AppConfig {
    /**
     * CONFIGURAZIONE URL SERVER API iOS
     *
     * 1. Per il SIMULATORE iOS (Default in sviluppo locale):
     *    Usa "http://localhost:8080/api/v1" (il simulatore condivide la rete dell'host Mac).
     *
     * 2. Per DISPOSITIVO FISICO iOS (iPhone/iPad via Wi-Fi):
     *    Usa l'indirizzo IP locale del tuo computer/server Mac nella rete locale, ad esempio:
     *    "http://192.168.1.100:8080/api/v1"
     *
     * 3. Per AMBIENTE DI PRODUZIONE / SERVER REMOTO:
     *    Usa il dominio pubblico HTTPS, ad esempio:
     *    "https://registro.tuascuola.it/api/v1"
     */
    public static var baseURL: String = "http://localhost:8080/api/v1"

    /**
     * URL WEBSOCKET PER AGGIORNAMENTI LIVE & NOTIFICHE IN TEMPO REALE
     */
    public static var wsURL: String = "ws://localhost:8080/api/v1/ws"

    public static let timeoutInterval: TimeInterval = 30.0
}
