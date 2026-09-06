package it.scuola.registro.student

import it.scuola.registro.student.network.HttpStudentApiService
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

class StudentRoleLoginSecurityTest {

    private lateinit var serverSocket: ServerSocket
    private lateinit var apiService: HttpStudentApiService
    @Volatile private var mockRoleResponse: String = "student"
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

                            val responseJson = """{"token":"valid_jwt_token_sample","user":{"id":"usr_student_1","role":"$mockRoleResponse","first_name":"Mario","last_name":"Rossi","email":"studente@scuola.it","class_name":"3A"}}"""
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
        apiService = HttpStudentApiService("http://127.0.0.1:$port/api/v1")
    }

    @After
    fun tearDown() {
        running = false
        try { serverSocket.close() } catch (_: Exception) {}
    }

    @Test
    fun login_rejectsParentCredentialsInStudentApp() = runBlocking {
        mockRoleResponse = "parent"

        val result = apiService.login("parent@scuola.it", "Password123!")

        assertTrue("Login con ruolo parent nell'app studente deve fallire", result.isFailure)
        val errorMessage = result.exceptionOrNull()?.message ?: ""
        assertTrue(
            "Messaggio di errore deve indicare che l'account non è uno studente",
            errorMessage.contains("Accesso non consentito: questo account non appartiene a uno studente.")
        )
    }

    @Test
    fun login_rejectsTeacherCredentialsInStudentApp() = runBlocking {
        mockRoleResponse = "teacher"

        val result = apiService.login("docente@scuola.it", "Password123!")

        assertTrue("Login con ruolo teacher nell'app studente deve fallire", result.isFailure)
        val errorMessage = result.exceptionOrNull()?.message ?: ""
        assertTrue(errorMessage.contains("Accesso non consentito"))
    }

    @Test
    fun login_allowsStudentCredentialsInStudentApp() = runBlocking {
        mockRoleResponse = "student"

        val result = apiService.login("studente@scuola.it", "Password123!")

        assertTrue("Login con ruolo student deve riuscire", result.isSuccess)
        assertEquals("valid_jwt_token_sample", result.getOrNull()?.first)
        assertEquals("Mario", result.getOrNull()?.second?.firstName)
    }
}
