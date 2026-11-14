import Foundation

public struct AppConfig {
    /**
     * Endpoint API iOS.
     * Per sicurezza sono richiesti HTTPS e WSS anche in rete locale.
     * La configurazione MDM rifiuta endpoint HTTP/WS o URL con credenziali/query.
     */
    public static var baseURL: String = "https://registro-backend-fdu2.onrender.com/api/v1"
    /**
     * URL WEBSOCKET PER AGGIORNAMENTI LIVE & NOTIFICHE IN TEMPO REALE
     */
    public static var wsURL: String = "wss://registro-backend-fdu2.onrender.com/api/v1/ws"

    public static let timeoutInterval: TimeInterval = 30.0
}
