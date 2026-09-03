package it.scuola.registro.parent.network

import it.scuola.registro.parent.data.ColloquioBooking
import it.scuola.registro.parent.data.ParentChild
import it.scuola.registro.parent.data.PendingAbsence
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONArray
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

interface ParentApiService {
    suspend fun login(email: String, password: String): Result<String>
    suspend fun getChildren(token: String): Result<List<ParentChild>>
    suspend fun getAbsences(token: String, childId: String): Result<List<PendingAbsence>>
    suspend fun submitJustification(token: String, absenceId: String, note: String): Result<Boolean>
    suspend fun bookColloquio(token: String, colloquioId: String): Result<Boolean>
}

class HttpParentApiService(
    private val baseUrl: String = "https://registro-backend-fdu2.onrender.com/api/v1"
) : ParentApiService {

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
                    val userObj = json.optJSONObject("user")
                    val role = userObj?.optString("role")
                    if (role != "parent") {
                        Result.failure(Exception("Accesso non consentito: questo account non appartiene a un genitore."))
                    } else {
                        val token = json.optString("token", json.optString("access_token"))
                        Result.success(token)
                    }
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    Result.failure(Exception("Login fallito (HTTP $responseCode): $err"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getChildren(token: String): Result<List<ParentChild>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/parents/my-children")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val list = mutableListOf<ParentChild>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            ParentChild(
                                id = obj.optString("id", i.toString()),
                                firstName = obj.optString("first_name", obj.optString("name", "Figlio")),
                                lastName = obj.optString("last_name", ""),
                                className = obj.optString("class_name", "Classe")
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero figli: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getAbsences(token: String, childId: String): Result<List<PendingAbsence>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/attendance/student/$childId")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val list = mutableListOf<PendingAbsence>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            PendingAbsence(
                                id = obj.optString("id", i.toString()),
                                childId = childId,
                                date = obj.optString("date", ""),
                                type = obj.optString("type", "Assenza"),
                                reason = obj.optString("reason", ""),
                                isJustified = obj.optBoolean("is_justified", false)
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero assenze: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun submitJustification(token: String, absenceId: String, note: String): Result<Boolean> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/attendance/$absenceId/justify")
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
                    put("reason", note)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    Result.success(true)
                } else {
                    Result.failure(Exception("Errore invio giustifica: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun bookColloquio(token: String, colloquioId: String): Result<Boolean> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/colloqui/$colloquioId/book")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "POST"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Content-Type", "application/json; charset=UTF-8")
                    doOutput = true
                    doInput = true
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    Result.success(true)
                } else {
                    Result.failure(Exception("Errore prenotazione colloquio: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
}
