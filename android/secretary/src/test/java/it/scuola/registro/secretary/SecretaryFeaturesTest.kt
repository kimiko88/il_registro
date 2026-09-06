package it.scuola.registro.secretary

import it.scuola.registro.secretary.data.CreateUserPayload
import it.scuola.registro.secretary.data.UpdateUserPayload
import it.scuola.registro.secretary.network.HttpSecretaryApiService
import it.scuola.registro.secretary.viewmodel.SecretaryViewModel
import kotlinx.coroutines.runBlocking
import org.junit.After
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.PrintWriter
import java.net.ServerSocket
import kotlin.concurrent.thread

class SecretaryFeaturesTest {

    private lateinit var serverSocket: ServerSocket
    private lateinit var apiService: HttpSecretaryApiService
    private lateinit var viewModel: SecretaryViewModel
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
                            val firstLine = reader.readLine() ?: return@use
                            val parts = firstLine.split(" ")
                            val method = parts.getOrNull(0) ?: "GET"
                            val path = parts.getOrNull(1) ?: "/"

                            var line: String?
                            var contentLength = 0
                            while (reader.readLine().also { line = it } != null) {
                                if (line!!.isEmpty()) break
                                if (line!!.lowercase().startsWith("content-length:")) {
                                    contentLength = line!!.substringAfter(":").trim().toIntOrNull() ?: 0
                                }
                            }
                            var requestBody = ""
                            if (contentLength > 0) {
                                val buf = CharArray(contentLength)
                                reader.read(buf, 0, contentLength)
                                requestBody = String(buf)
                            }

                            val responseJson = when {
                                path.startsWith("/api/v1/admin/dashboard/stats") ->
                                    """{"total_students":12,"total_teachers":5,"total_users":25,"total_documents":8}"""

                                path.startsWith("/api/v1/users") && method == "GET" ->
                                    """{"users":[{"id":"u1","first_name":"Giovanni","last_name":"Rossi","email":"giovanni@scuola.it","role":"teacher","fiscal_code":"RSSGNN80A01H501Z"},{"id":"u2","first_name":"Luca","last_name":"Verdi","email":"luca@scuola.it","role":"student","class_name":"3A"}],"total_count":2}"""

                                path.startsWith("/api/v1/users") && method == "POST" ->
                                    """{"id":"u_new_99","first_name":"Nuovo","last_name":"Docente","email":"nuovo@scuola.it","role":"teacher"}"""

                                path.startsWith("/api/v1/users/u1") && method == "PATCH" ->
                                    """{"id":"u1","first_name":"Giovanni","last_name":"Rossi Modificato"}"""

                                path.startsWith("/api/v1/users/u1") && method == "DELETE" ->
                                    """{"message":"user deleted"}"""

                                path.startsWith("/api/v1/certificates/generate") && method == "POST" ->
                                    """{"id":"cert_new_1","pdf_url":"/api/v1/certificates/cert_new_1.pdf","cert":{"id":"cert_new_1","protocol_no":"PROT-2026-999","student_name":"Luca Verdi","type":"iscrizione","issued_at":"2026-09-03T18:00:00Z"}}"""

                                path.startsWith("/api/v1/certificates") && method == "GET" ->
                                    """[{"id":"cert_101","student_id":"u2","student_name":"Luca Verdi","class_name":"3A","type":"frequenza","protocol_no":"PROT-2026-001","issued_at":"2026-09-01T10:00:00Z","academic_year":"2025/2026","pdf_url":"/api/v1/certificates/cert_101.pdf"}]"""

                                path.startsWith("/api/v1/audit-log") ->
                                    """{"data":[{"id":"aud_1","actor_name":"Segreteria 1","actor_role":"secretary","action":"USER_CREATE","entity_type":"user","details":"Creato utente docente","ip_address":"192.168.1.50","created_at":"2026-09-03T17:30:00Z"}],"total":1}"""

                                path.startsWith("/api/v1/auth/me") ->
                                    """{"id":"usr_sec_1","first_name":"Anna","last_name":"Segretaria","email":"segreteria@scuola.it","role":"secretary","school_id":"sch_1"}"""

                                path.contains("/change-password") ->
                                    """{"message":"Password aggiornata con successo"}"""

                                else -> """{"message":"ok"}"""
                            }

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
        viewModel = SecretaryViewModel(apiService)
    }

    @After
    fun tearDown() {
        running = false
        try { serverSocket.close() } catch (_: Exception) {}
    }

    @Test
    fun getUsers_parsesWrappedObjectSuccessfully() = runBlocking {
        val result = apiService.getUsers("mock_token")
        assertTrue("getUsers deve avere successo", result.isSuccess)
        val users = result.getOrNull()
        assertNotNull(users)
        assertEquals(2, users!!.size)
        assertEquals("Giovanni", users[0].firstName)
        assertEquals("teacher", users[0].role)
        assertEquals("Luca", users[1].firstName)
        assertEquals("student", users[1].role)
    }

    @Test
    fun getDashboardStats_parsesAccurately() = runBlocking {
        val result = apiService.getDashboardStats("mock_token")
        assertTrue("getDashboardStats deve avere successo", result.isSuccess)
        val stats = result.getOrNull()
        assertNotNull(stats)
        assertEquals(12, stats!!.totalStudents)
        assertEquals(5, stats.totalTeachers)
        assertEquals(25, stats.totalUsers)
    }

    @Test
    fun getCertificates_parsesCertificatesList() = runBlocking {
        val result = apiService.getCertificates("mock_token")
        assertTrue(result.isSuccess)
        val certs = result.getOrNull()
        assertNotNull(certs)
        assertEquals(1, certs!!.size)
        assertEquals("PROT-2026-001", certs[0].protocolNo)
        assertEquals("Luca Verdi", certs[0].studentName)
    }

    @Test
    fun getAuditLogs_parsesAuditList() = runBlocking {
        val result = apiService.getAuditLogs("mock_token")
        assertTrue(result.isSuccess)
        val logs = result.getOrNull()
        assertNotNull(logs)
        assertEquals(1, logs!!.size)
        assertEquals("USER_CREATE", logs[0].action)
        assertEquals("Segreteria 1", logs[0].actorName)
    }

    @Test
    fun viewModel_loadFromDatabase_populatesAllState() = runBlocking {
        val success = viewModel.loadFromDatabase("mock_token")
        assertTrue(success)
        assertEquals(2, viewModel.usersList.size)
        assertEquals(12, viewModel.stats.totalStudents)
        assertEquals(5, viewModel.stats.totalTeachers)
        assertEquals(1, viewModel.certificateRequests.size)
        assertEquals(1, viewModel.auditLogs.size)
        assertNotNull(viewModel.currentUserProfile)
        assertEquals("Anna", viewModel.currentUserProfile?.firstName)
    }

    @Test
    fun viewModel_createUserOnline_addsToState() = runBlocking {
        viewModel.loadFromDatabase("mock_token")
        val initialSize = viewModel.usersList.size
        val created = viewModel.createUserOnline(
            "mock_token",
            CreateUserPayload("Nuovo", "Docente", "nuovo@scuola.it", "Password123!", "teacher")
        )
        assertTrue(created)
        assertEquals(initialSize + 1, viewModel.usersList.size)
        assertEquals("nuovo@scuola.it", viewModel.usersList[0].email)
    }

    @Test
    fun viewModel_changePasswordOnline_succeeds() = runBlocking {
        val res = viewModel.changePasswordOnline("mock_token", "usr_1", "OldPassword123!", "NewPassword123!")
        assertTrue(res.isSuccess)
        assertNotNull(viewModel.successMessage)
    }
}
