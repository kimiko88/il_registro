package it.scuola.registro.secretary.network

import it.scuola.registro.secretary.data.CertificateRequest
import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.data.ScrutinyClassStatus
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
    suspend fun login(email: String, password: String): Result<String>
    suspend fun getUsers(token: String): Result<List<ManagedUser>>
    suspend fun createUser(token: String, user: ManagedUser): Result<ManagedUser>
    suspend fun getScrutinyStatus(token: String): Result<List<ScrutinyClassStatus>>
    suspend fun requestCertificatePdf(token: String, studentId: String, type: String): Result<CertificateRequest>
}

class HttpSecretaryApiService(
    private val baseUrl: String = "https://api.scuola.registro.it/api/v1"
) : SecretaryApiService {

    override suspend fun login(email: String, password: String): Result<String> =
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
                    val token = json.optString("token", json.optString("access_token"))
                    Result.success(token)
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    Result.failure(Exception("Login fallito (HTTP $responseCode): $err"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getUsers(token: String): Result<List<ManagedUser>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/users")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val list = mutableListOf<ManagedUser>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            ManagedUser(
                                id = obj.optString("id", i.toString()),
                                firstName = obj.optString("first_name", obj.optString("name", "")),
                                lastName = obj.optString("last_name", ""),
                                email = obj.optString("email", ""),
                                role = obj.optString("role", "staff")
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

    override suspend fun createUser(token: String, user: ManagedUser): Result<ManagedUser> =
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

                val payload = JSONObject().apply {
                    put("first_name", user.firstName)
                    put("last_name", user.lastName)
                    put("email", user.email)
                    put("role", user.role)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val obj = JSONObject(response)
                    Result.success(
                        ManagedUser(
                            id = obj.optString("id", user.id),
                            firstName = obj.optString("first_name", user.firstName),
                            lastName = obj.optString("last_name", user.lastName),
                            email = obj.optString("email", user.email),
                            role = obj.optString("role", user.role)
                        )
                    )
                } else {
                    Result.failure(Exception("Errore creazione utente: HTTP ${conn.responseCode}"))
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

    override suspend fun requestCertificatePdf(token: String, studentId: String, type: String): Result<CertificateRequest> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/certificates")
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
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val obj = JSONObject(response)
                    Result.success(
                        CertificateRequest(
                            id = obj.optString("id", "cert_1"),
                            studentId = studentId,
                            certificateType = type,
                            status = obj.optString("status", "completato"),
                            generatedPdfUrl = obj.optString("pdf_url", "/api/v1/certificates/${obj.optString("id")}.pdf")
                        )
                    )
                } else {
                    Result.failure(Exception("Errore generazione certificato: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
}
