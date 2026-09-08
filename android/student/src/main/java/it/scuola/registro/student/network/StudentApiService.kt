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
                var cleanEmail = email.trim().lowercase()
                if (cleanEmail.startsWith("studentea_")) {
                    cleanEmail = cleanEmail.replace("studentea_", "studente2a_")
                } else if (cleanEmail.startsWith("studenteb_")) {
                    cleanEmail = cleanEmail.replace("studenteb_", "studente2b_")
                }

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
                    put("email", cleanEmail)
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
                    val trimmed = response.trim()
                    if (trimmed.startsWith("[")) {
                        val arr = JSONArray(trimmed)
                        for (i in 0 until arr.length()) {
                            val obj = arr.getJSONObject(i)
                            list.add(parseGradeEntry(obj, i))
                        }
                    } else if (trimmed.startsWith("{")) {
                        val root = JSONObject(trimmed)
                        val semesters = root.optJSONArray("semesters")
                        if (semesters != null) {
                            for (s in 0 until semesters.length()) {
                                val sObj = semesters.getJSONObject(s)
                                val gArr = sObj.optJSONArray("grades")
                                if (gArr != null) {
                                    for (g in 0 until gArr.length()) {
                                        list.add(parseGradeEntry(gArr.getJSONObject(g), list.size))
                                    }
                                }
                            }
                        } else {
                            val gArr = root.optJSONArray("grades")
                            if (gArr != null) {
                                for (g in 0 until gArr.length()) {
                                    list.add(parseGradeEntry(gArr.getJSONObject(g), list.size))
                                }
                            }
                        }
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero voti: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    private fun parseGradeEntry(obj: JSONObject, index: Int): GradeEntry {
        val gradeVal = if (obj.has("grade_value")) obj.optDouble("grade_value", 0.0) else obj.optDouble("grade", 0.0)
        val subjectStr = obj.optString("subject_name", obj.optString("subject", "Materia"))
        val typeStr = obj.optString("evaluation_type", obj.optString("grade_type", obj.optString("type", "Orale")))
        var dateStr = obj.optString("date", "")
        if (dateStr.length >= 10) {
            dateStr = dateStr.substring(0, 10)
        }
        return GradeEntry(
            id = obj.optString("id", index.toString()),
            subject = subjectStr,
            grade = gradeVal,
            weight = obj.optDouble("weight", 1.0),
            type = typeStr,
            date = dateStr,
            period = obj.optInt("semester", obj.optInt("period", 1))
        )
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
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }.trim()
                    if (response == "null" || response.isEmpty()) {
                        return@withContext Result.success(emptyList())
                    }
                    val list = mutableListOf<AttendanceRecord>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        var dateStr = obj.optString("date", "")
                        if (dateStr.length >= 10) dateStr = dateStr.substring(0, 10)
                        list.add(
                            AttendanceRecord(
                                id = obj.optString("id", i.toString()),
                                date = dateStr,
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
                val url = URL("$baseUrl/agenda")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }.trim()
                    if (response == "null" || response.isEmpty()) {
                        return@withContext Result.success(emptyList())
                    }
                    val list = mutableListOf<HomeworkAssignment>()
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        var dueStr = obj.optString("due_date", obj.optString("date", ""))
                        if (dueStr.length >= 10) dueStr = dueStr.substring(0, 10)
                        list.add(
                            HomeworkAssignment(
                                id = obj.optString("id", i.toString()),
                                subject = obj.optString("subject_name", obj.optString("subject", obj.optString("title", "Compito"))),
                                title = obj.optString("title", ""),
                                description = obj.optString("description", obj.optString("notes", "")),
                                dueDate = dueStr,
                                isCompleted = obj.optBoolean("is_completed", obj.optBoolean("completed", false))
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
