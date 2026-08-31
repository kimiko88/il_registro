package it.scuola.registro.student.network

interface WebSocketListener {
    fun onConnected()
    fun onMessageReceived(message: String)
    fun onDisconnected()
}

class StudentWebSocketClient {
    private var isConnected: Boolean = false
    private var listener: WebSocketListener? = null

    fun connect(wsUrl: String, token: String, listener: WebSocketListener) {
        this.listener = listener
        if (wsUrl.isNotBlank() && token.isNotBlank()) {
            isConnected = true
            listener.onConnected()
        }
    }

    fun simulateIncomingMessage(message: String) {
        if (isConnected) {
            listener?.onMessageReceived(message)
        }
    }

    fun disconnect() {
        isConnected = false
        listener?.onDisconnected()
        listener = null
    }

    fun isConnectionActive(): Boolean = isConnected
}
