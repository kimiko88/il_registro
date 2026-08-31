package it.scuola.registro.student

import it.scuola.registro.student.network.StudentWebSocketClient
import it.scuola.registro.student.network.WebSocketListener
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class StudentWebSocketClientTest {

    private lateinit var wsClient: StudentWebSocketClient

    @Before
    fun setUp() {
        wsClient = StudentWebSocketClient()
    }

    @Test
    fun connectAndReceiveMessage_triggersCallbacks() {
        var connected = false
        var receivedMsg = ""
        var disconnected = false

        wsClient.connect(
            wsUrl = "wss://api.scuola.it/ws",
            token = "jwt_student_token",
            listener = object : WebSocketListener {
                override fun onConnected() {
                    connected = true
                }

                override fun onMessageReceived(message: String) {
                    receivedMsg = message
                }

                override fun onDisconnected() {
                    disconnected = true
                }
            }
        )

        assertTrue(connected)
        assertTrue(wsClient.isConnectionActive())

        wsClient.simulateIncomingMessage("{\"event\":\"NEW_GRADE\",\"subject\":\"Matematica\",\"grade\":9.0}")
        assertEquals("{\"event\":\"NEW_GRADE\",\"subject\":\"Matematica\",\"grade\":9.0}", receivedMsg)

        wsClient.disconnect()
        assertFalse(wsClient.isConnectionActive())
        assertTrue(disconnected)
    }
}
