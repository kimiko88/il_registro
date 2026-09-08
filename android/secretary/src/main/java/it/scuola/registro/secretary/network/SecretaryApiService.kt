package it.scuola.registro.secretary.network

import it.scuola.registro.secretary.data.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONArray
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

interface SecretaryApiService {
    suspend fun login(email: String, password: String): Result<Pair<String, UserProfile?>>
    suspend fun getDashboardStats(token: String): Result<DashboardStats>
    suspend fun getUsers(token: String): Result<List<ManagedUser>>
    suspend fun createUser(token: String, payload: CreateUserPayload): Result<ManagedUser>
    suspend fun updateUser(token: String, id: String, payload: UpdateUserPayload): Result<Boolean>
    suspend fun deleteUser(token: String, id: String): Result<Boolean>
    suspend fun getScrutinyStatus(token: String): Result<List<ScrutinyClassStatus>>
    suspend fun getCertificates(token: String): Result<List<CertificateRequest>>
    suspend fun generateCertificate(token: String, studentId: String, type: String, academicYear: String, notes: String): Result<CertificateRequest>
    suspend fun getAuditLogs(token: String): Result<List<AuditLogEntry>>
    suspend fun getCurrentUser(token: String): Result<UserProfile>
    suspend fun changePassword(token: String, userId: String, currentPassword: String, newPassword: String): Result<String>
}

class HttpSecretaryApiService(
    private val baseUrl: String = "https://registro-backend-fdu2.onrender.com/api/v1"
) : SecretaryApiService {

    override suspend fun login(email: String, password: String): Result<Pair<String, UserProfile?>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/auth/login")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "POST"
                    setRequestProperty("Content-Type", "application/json; charset=UTF-8")
                    doOutput = true
                    doInput = true
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                val payload = JSONObject().apply {
                    put("email", email)
                    put("password", password)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                val responseCode = conn.responseCode
                if (responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val json = JSONObject(response)
                    val userObj = json.optJSONObject("user")
                    val role = userObj?.optString("role") ?: ""
                    val allowedRoles = listOf("secretary", "admin", "superadmin", "principal", "vice_principal")
                    if (role !in allowedRoles) {
                        Result.failure(Exception("Accesso non consentito: questo account non ha i permessi di segreteria o amministrazione."))
                    } else {
                        val token = json.optString("token", json.optString("access_token"))
                        val profile = userObj?.let {
                            UserProfile(
                                id = it.optString("id"),
                                email = it.optString("email"),
                                firstName = it.optString("first_name"),
                                lastName = it.optString("last_name"),
                                role = it.optString("role"),
                                schoolId = it.optString("school_id")
                            )
                        }
                        Result.success(Pair(token, profile))
                    }
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    Result.failure(Exception("Login fallito (HTTP $responseCode): $err"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getDashboardStats(token: String): Result<DashboardStats> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/admin/dashboard/stats")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val obj = JSONObject(response)
                    Result.success(
                        DashboardStats(
                            totalStudents = obj.optInt("total_students", 0),
                            totalTeachers = obj.optInt("total_teachers", 0),
                            totalUsers = obj.optInt("total_users", 0),
                            totalDocuments = obj.optInt("total_documents", 0)
                        )
                    )
                } else {
                    Result.failure(Exception("Errore statistiche dashboard: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getUsers(token: String): Result<List<ManagedUser>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/users?page_size=200")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }.trim()
                    val list = mutableListOf<ManagedUser>()
                    val arr = if (response.startsWith("{")) {
                        val rootObj = JSONObject(response)
                        rootObj.optJSONArray("users") ?: JSONArray()
                    } else {
                        JSONArray(response)
                    }

                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            ManagedUser(
                                id = obj.optString("id", i.toString()),
                                firstName = obj.optString("first_name", obj.optString("name", "")),
                                lastName = obj.optString("last_name", ""),
                                email = obj.optString("email", ""),
                                role = obj.optString("role", "staff"),
                                isActive = obj.optBoolean("is_active", true),
                                fiscalCode = obj.optString("fiscal_code", ""),
                                phoneNumber = obj.optString("phone_number", ""),
                                className = obj.optString("class_name", ""),
                                classId = obj.optString("class_id", "")
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero utenti: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun createUser(token: String, payload: CreateUserPayload): Result<ManagedUser> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/users")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "POST"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Content-Type", "application/json; charset=UTF-8")
                    doOutput = true
                    doInput = true
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                val body = JSONObject().apply {
                    put("first_name", payload.firstName)
                    put("last_name", payload.lastName)
                    put("email", payload.email)
                    put("password", payload.password)
                    put("role", payload.role)
                    if (payload.fiscalCode.isNotBlank()) put("fiscal_code", payload.fiscalCode.trim().uppercase())
                    if (payload.phoneNumber.isNotBlank()) put("phone_number", payload.phoneNumber.trim())
                    if (!payload.classId.isNullOrBlank()) put("class_id", payload.classId)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(body.toString()) }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val obj = JSONObject(response)
                    Result.success(
                        ManagedUser(
                            id = obj.optString("id"),
                            firstName = obj.optString("first_name", payload.firstName),
                            lastName = obj.optString("last_name", payload.lastName),
                            email = obj.optString("email", payload.email),
                            role = obj.optString("role", payload.role),
                            fiscalCode = obj.optString("fiscal_code", payload.fiscalCode),
                            phoneNumber = obj.optString("phone_number", payload.phoneNumber)
                        )
                    )
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    Result.failure(Exception("Errore creazione utente: $err"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun updateUser(token: String, id: String, payload: UpdateUserPayload): Result<Boolean> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/users/$id")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "PATCH"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Content-Type", "application/json; charset=UTF-8")
                    doOutput = true
                    doInput = true
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                val body = JSONObject().apply {
                    if (payload.firstName.isNotBlank()) put("first_name", payload.firstName)
                    if (payload.lastName.isNotBlank()) put("last_name", payload.lastName)
                    if (payload.fiscalCode.isNotBlank()) put("fiscal_code", payload.fiscalCode.trim().uppercase())
                    if (payload.phoneNumber.isNotBlank()) put("phone_number", payload.phoneNumber.trim())
                    if (!payload.role.isNullOrBlank()) put("role", payload.role)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(body.toString()) }

                if (conn.responseCode in 200..299) {
                    Result.success(true)
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    Result.failure(Exception("Errore modifica utente: $err"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun deleteUser(token: String, id: String): Result<Boolean> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/users/$id")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "DELETE"
                    setRequestProperty("Authorization", "Bearer $token")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    Result.success(true)
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    Result.failure(Exception("Errore eliminazione utente: $err"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getScrutinyStatus(token: String): Result<List<ScrutinyClassStatus>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/scrutiny/overview")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val list = mutableListOf<ScrutinyClassStatus>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            ScrutinyClassStatus(
                                classId = obj.optString("class_id", i.toString()),
                                className = obj.optString("class_name", "Classe"),
                                period = obj.optInt("period", 2),
                                isLocked = obj.optBoolean("is_locked", false),
                                status = obj.optString("status", "in_corso")
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero scrutinio: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getCertificates(token: String): Result<List<CertificateRequest>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/certificates")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val list = mutableListOf<CertificateRequest>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        val rawType = obj.optString("type", "iscrizione")
                        val displayType = when (rawType) {
                            "iscrizione" -> "Iscrizione"
                            "frequenza" -> "Frequenza"
                            "promozione" -> "Promozione"
                            "condotta" -> "Buona Condotta"
                            else -> rawType.replaceFirstChar { it.uppercase() }
                        }
                        list.add(
                            CertificateRequest(
                                id = obj.optString("id", i.toString()),
                                studentId = obj.optString("student_id"),
                                studentName = obj.optString("student_name", "Studente"),
                                className = obj.optString("class_name", ""),
                                certificateType = displayType,
                                protocolNo = obj.optString("protocol_no", "PROT-${i + 1}"),
                                issuedAt = obj.optString("issued_at", "").take(10),
                                academicYear = obj.optString("academic_year", "2025/2026"),
                                generatedPdfUrl = obj.optString("pdf_url")
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero certificati: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun generateCertificate(
        token: String,
        studentId: String,
        type: String,
        academicYear: String,
        notes: String
    ): Result<CertificateRequest> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/certificates/generate")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "POST"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Content-Type", "application/json; charset=UTF-8")
                    doOutput = true
                    doInput = true
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                val payload = JSONObject().apply {
                    put("student_id", studentId)
                    put("type", type)
                    if (academicYear.isNotBlank()) put("academic_year", academicYear)
                    if (notes.isNotBlank()) put("notes", notes)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val root = JSONObject(response)
                    val certObj = root.optJSONObject("cert") ?: root
                    val displayType = when (type) {
                        "iscrizione" -> "Iscrizione"
                        "frequenza" -> "Frequenza"
                        "promozione" -> "Promozione"
                        "condotta" -> "Buona Condotta"
                        else -> type.replaceFirstChar { it.uppercase() }
                    }
                    Result.success(
                        CertificateRequest(
                            id = certObj.optString("id", root.optString("id")),
                            studentId = studentId,
                            studentName = certObj.optString("student_name", "Studente"),
                            className = certObj.optString("class_name", ""),
                            certificateType = displayType,
                            protocolNo = certObj.optString("protocol_no", "PROT-REG"),
                            issuedAt = certObj.optString("issued_at", "").take(10),
                            academicYear = academicYear,
                            generatedPdfUrl = root.optString("pdf_url")
                        )
                    )
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    Result.failure(Exception("Errore generazione certificato: $err"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getAuditLogs(token: String): Result<List<AuditLogEntry>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/audit-log?limit=50")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val root = JSONObject(response)
                    val arr = root.optJSONArray("data") ?: JSONArray()
                    val list = mutableListOf<AuditLogEntry>()
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        val rawDate = obj.optString("created_at")
                        val formattedDate = if (rawDate.length >= 16) {
                            rawDate.substring(0, 10) + " " + rawDate.substring(11, 16)
                        } else rawDate

                        list.add(
                            AuditLogEntry(
                                id = obj.optString("id", i.toString()),
                                actorName = obj.optString("actor_name", "Sistema"),
                                actorRole = obj.optString("actor_role", "system"),
                                action = obj.optString("action", "OP"),
                                entityType = obj.optString("entity_type", ""),
                                details = obj.optString("details", ""),
                                ipAddress = obj.optString("ip_address", "127.0.0.1"),
                                createdAt = formattedDate
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero log di sistema: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getCurrentUser(token: String): Result<UserProfile> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/auth/me")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val obj = JSONObject(response)
                    Result.success(
                        UserProfile(
                            id = obj.optString("id"),
                            email = obj.optString("email"),
                            firstName = obj.optString("first_name"),
                            lastName = obj.optString("last_name"),
                            role = obj.optString("role"),
                            schoolId = obj.optString("school_id")
                        )
                    )
                } else {
                    Result.failure(Exception("Errore profilo utente: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun changePassword(
        token: String,
        userId: String,
        currentPassword: String,
        newPassword: String
    ): Result<String> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/users/$userId/change-password")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "POST"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Content-Type", "application/json; charset=UTF-8")
                    doOutput = true
                    doInput = true
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                val payload = JSONObject().apply {
                    put("current_password", currentPassword)
                    put("new_password", newPassword)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    Result.success("Password aggiornata con successo")
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    val errMsg = try {
                        JSONObject(err).optString("message", JSONObject(err).optString("error", err))
                    } catch (_: Exception) {
                        err
                    }
                    Result.failure(Exception("Errore aggiornamento password: $errMsg"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
}
