package it.scuola.registro.student.network

interface WebSocketListener {
    fun onConnected()
    fun onMessageReceived(message: String)
    fun onDisconnected()
}

/**
 * Placeholder until a real TLS WebSocket transport with authenticated lifecycle
 * handling is implemented. It must never claim a connection or synthesize events.
 */
class StudentWebSocketClient {
    fun connect(wsUrl: String, token: String, listener: WebSocketListener) {
        listener.onDisconnected()
    }

    fun simulateIncomingMessage(message: String) = Unit

    fun disconnect() = Unit

    fun isConnectionActive(): Boolean = false
}
