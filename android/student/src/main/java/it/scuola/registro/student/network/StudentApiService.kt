package it.scuola.registro.student.network

import it.scuola.registro.student.data.AttendanceRecord
import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.data.HomeworkAssignment
import it.scuola.registro.student.data.ScrutinyReportCard
import it.scuola.registro.student.data.StudentUser
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONArray
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

interface StudentApiService {
    suspend fun login(email: String, password: String): Result<Pair<String, StudentUser>>
    suspend fun getGrades(token: String): Result<List<GradeEntry>>
    suspend fun getAttendance(token: String): Result<List<AttendanceRecord>>
    suspend fun getHomework(token: String): Result<List<HomeworkAssignment>>
    suspend fun getReportCard(token: String): Result<ScrutinyReportCard>
}

class HttpStudentApiService(
    private val baseUrl: String = "https://registro-backend-fdu2.onrender.com/api/v1"
) : StudentApiService {

    override suspend fun login(email: String, password: String): Result<Pair<String, StudentUser>> =
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
                    val userObj = json.optJSONObject("user") ?: JSONObject()
                    val role = userObj.optString("role")
                    if (role != "student") {
                        Result.failure(Exception("Accesso non consentito: questo account non appartiene a uno studente."))
                    } else {
                        val token = json.optString("token", json.optString("access_token"))
                        val user = StudentUser(
                            id = userObj.optString("id", "s1"),
                            firstName = userObj.optString("first_name", userObj.optString("name", "Studente")),
                            lastName = userObj.optString("last_name", ""),
                            email = userObj.optString("email", email),
                            className = userObj.optString("class_name", "Classe")
                        )
                        Result.success(Pair(token, user))
                    }
                } else {
                    val err = BufferedReader(InputStreamReader(conn.errorStream ?: conn.inputStream)).use { it.readText() }
                    Result.failure(Exception("Login fallito (HTTP $responseCode): $err"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getGrades(token: String): Result<List<GradeEntry>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/grades/my-grades")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val list = mutableListOf<GradeEntry>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            GradeEntry(
                                id = obj.optString("id", i.toString()),
                                subject = obj.optString("subject_name", obj.optString("subject", "Materia")),
                                grade = obj.optDouble("grade", 0.0),
                                weight = obj.optDouble("weight", 1.0),
                                type = obj.optString("type", "Orale"),
                                date = obj.optString("date", ""),
                                period = obj.optInt("period", 1)
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero voti: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getAttendance(token: String): Result<List<AttendanceRecord>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/attendance/my-attendance")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val list = mutableListOf<AttendanceRecord>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            AttendanceRecord(
                                id = obj.optString("id", i.toString()),
                                date = obj.optString("date", ""),
                                type = obj.optString("status", obj.optString("type", "Presenza")),
                                isJustified = obj.optBoolean("is_justified", false),
                                reason = obj.optString("justification_reason", obj.optString("reason", ""))
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero presenze: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getHomework(token: String): Result<List<HomeworkAssignment>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/agenda/my-homework")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val list = mutableListOf<HomeworkAssignment>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            HomeworkAssignment(
                                id = obj.optString("id", i.toString()),
                                subject = obj.optString("subject_name", obj.optString("subject", "")),
                                title = obj.optString("title", ""),
                                description = obj.optString("description", ""),
                                dueDate = obj.optString("due_date", ""),
                                isCompleted = obj.optBoolean("is_completed", false)
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero compiti: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getReportCard(token: String): Result<ScrutinyReportCard> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/scrutiny/my-report-card")
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
                    val gradesMap = mutableMapOf<String, Double>()
                    val gradesJson = obj.optJSONObject("grades")
                    gradesJson?.keys()?.forEach { key ->
                        gradesMap[key] = gradesJson.optDouble(key, 0.0)
                    }

                    Result.success(
                        ScrutinyReportCard(
                            studentId = obj.optString("student_id", ""),
                            period = obj.optInt("period", 2),
                            grades = gradesMap,
                            conductGrade = obj.optInt("conduct_grade", 8),
                            finalDecision = obj.optString("final_decision", "Ammesso"),
                            credits = obj.optInt("credits", 0)
                        )
                    )
                } else {
                    Result.failure(Exception("Errore recupero pagella: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
}
