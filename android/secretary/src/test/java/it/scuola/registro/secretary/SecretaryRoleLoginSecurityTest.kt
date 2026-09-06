package it.scuola.registro.secretary

import it.scuola.registro.secretary.network.HttpSecretaryApiService
import kotlinx.coroutines.runBlocking
import org.junit.After
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.PrintWriter
import java.net.ServerSocket
import java.net.Socket
import kotlin.concurrent.thread

class SecretaryRoleLoginSecurityTest {

    private lateinit var serverSocket: ServerSocket
    private lateinit var apiService: HttpSecretaryApiService
    @Volatile private var mockRoleResponse: String = "secretary"
    @Volatile private var running = true

    @Before
    fun setUp() {
        serverSocket = ServerSocket(0)
        running = true
        thread(isDaemon = true) {
            while (running) {
                try {
                    val socket = serverSocket.accept()
                    thread(isDaemon = true) {
                        socket.use { s ->
                            val reader = BufferedReader(InputStreamReader(s.getInputStream()))
                            var line: String?
                            var contentLength = 0
                            while (reader.readLine().also { line = it } != null) {
                                if (line!!.isEmpty()) break
                                if (line!!.lowercase().startsWith("content-length:")) {
                                    contentLength = line!!.substringAfter(":").trim().toIntOrNull() ?: 0
                                }
                            }
                            if (contentLength > 0) {
                                val buf = CharArray(contentLength)
                                reader.read(buf, 0, contentLength)
                            }

                            val responseJson = """{"token":"valid_jwt_token_sample","user":{"id":"usr_sec_1","role":"$mockRoleResponse"}}"""
                            val bytes = responseJson.toByteArray(Charsets.UTF_8)
                            val out = s.getOutputStream()
                            val writer = PrintWriter(out)
                            writer.print("HTTP/1.1 200 OK\r\n")
                            writer.print("Content-Type: application/json\r\n")
                            writer.print("Content-Length: ${bytes.size}\r\n")
                            writer.print("Connection: close\r\n\r\n")
                            writer.flush()
                            out.write(bytes)
                            out.flush()
                        }
                    }
                } catch (_: Exception) {}
            }
        }
        val port = serverSocket.localPort
        apiService = HttpSecretaryApiService("http://127.0.0.1:$port/api/v1")
    }

    @After
    fun tearDown() {
        running = false
        try { serverSocket.close() } catch (_: Exception) {}
    }

    @Test
    fun login_rejectsParentCredentialsInSecretaryApp() = runBlocking {
        mockRoleResponse = "parent"

        val result = apiService.login("parent@scuola.it", "Password123!")

        assertTrue("Login con ruolo parent nell'app segreteria deve fallire", result.isFailure)
        val errorMessage = result.exceptionOrNull()?.message ?: ""
        assertTrue(
            "Messaggio di errore deve indicare mancanza di permessi",
            errorMessage.contains("Accesso non consentito: questo account non ha i permessi di segreteria o amministrazione.")
        )
    }

    @Test
    fun login_rejectsStudentCredentialsInSecretaryApp() = runBlocking {
        mockRoleResponse = "student"

        val result = apiService.login("student@scuola.it", "Password123!")

        assertTrue("Login con ruolo student nell'app segreteria deve fallire", result.isFailure)
        val errorMessage = result.exceptionOrNull()?.message ?: ""
        assertTrue(errorMessage.contains("Accesso non consentito"))
    }

    @Test
    fun login_allowsSecretaryCredentialsInSecretaryApp() = runBlocking {
        mockRoleResponse = "secretary"

        val result = apiService.login("segreteria@scuola.it", "Password123!")

        assertTrue("Login con ruolo secretary deve riuscire", result.isSuccess)
        assertEquals("valid_jwt_token_sample", result.getOrNull()?.first)
    }

    @Test
    fun login_allowsAdminCredentialsInSecretaryApp() = runBlocking {
        mockRoleResponse = "admin"

        val result = apiService.login("admin@scuola.it", "Password123!")

        assertTrue("Login con ruolo admin deve riuscire", result.isSuccess)
        assertEquals("valid_jwt_token_sample", result.getOrNull()?.first)
    }
}
