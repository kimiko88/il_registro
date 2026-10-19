package it.scuola.registro.teacher.network

import it.scuola.registro.teacher.data.ClassSession
import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONArray
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

interface TeacherApiService {
    suspend fun login(email: String, password: String): Result<String>
    suspend fun getClasses(token: String): Result<List<ClassSession>>
    suspend fun getClassStudents(token: String, classId: String): Result<List<StudentRollCall>>
    suspend fun signLesson(token: String, classId: String, topic: String): Result<ClassSession>
    suspend fun submitRollCall(token: String, classId: String, records: List<StudentRollCall>): Result<Boolean>
    suspend fun submitGrade(token: String, proposal: GradeProposal): Result<Boolean>
    suspend fun saveDeferredScrutiny(token: String, resolution: DeferredScrutinyResolution): Result<Boolean>
}

class HttpTeacherApiService(
    private val baseUrl: String = "https://registro-backend-fdu2.onrender.com/api/v1"
) : TeacherApiService {

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
                    if (role != "teacher" && role != "coordinator") {
                        Result.failure(Exception("Accesso non consentito: questo account non appartiene a un docente."))
                    } else {
                        val token = json.optString("access_token", json.optString("token"))
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

    override suspend fun getClasses(token: String): Result<List<ClassSession>> =
        withContext(Dispatchers.IO) {
            try {
                var url = URL("$baseUrl/teacher/classes")
                var conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    connectTimeout = 10000
                    readTimeout = 10000
                }
                var response = if (conn.responseCode in 200..299) {
                    BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }.trim()
                } else {
                    ""
                }
                if (response.isEmpty() || response == "[]" || response == "null") {
                    url = URL("$baseUrl/classes")
                    conn = (url.openConnection() as HttpURLConnection).apply {
                        requestMethod = "GET"
                        setRequestProperty("Authorization", "Bearer $token")
                        connectTimeout = 10000
                        readTimeout = 10000
                    }
                    if (conn.responseCode in 200..299) {
                        response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }.trim()
                    }
                }

                val list = mutableListOf<ClassSession>()
                if (response.startsWith("[")) {
                    val arr = JSONArray(response)
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            ClassSession(
                                classId = obj.optString("id", "$i"),
                                className = obj.optString("name", obj.optString("class_name", "Classe")),
                                subject = obj.optString("subject_name", obj.optString("subject", "Matematica")),
                                hourSlot = "Orario lezioni",
                                isSigned = false,
                                lessonTopic = ""
                            )
                        )
                    }
                }
                Result.success(list)
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun getClassStudents(token: String, classId: String): Result<List<StudentRollCall>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/users?role=student&class_id=$classId")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Accept", "application/json")
                    connectTimeout = 10000
                    readTimeout = 10000
                }
                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }.trim()
                    val list = mutableListOf<StudentRollCall>()
                    val arr = if (response.startsWith("{")) {
                        val root = JSONObject(response)
                        root.optJSONArray("users") ?: JSONArray()
                    } else if (response.startsWith("[")) {
                        JSONArray(response)
                    } else {
                        JSONArray()
                    }

                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        val fName = obj.optString("first_name", "")
                        val lName = obj.optString("last_name", "")
                        val fullName = "$fName $lName".trim().ifEmpty { obj.optString("email", "Studente") }
                        list.add(
                            StudentRollCall(
                                studentId = obj.optString("id", "$i"),
                                fullName = fullName,
                                status = "presente"
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore recupero studenti classe: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun signLesson(token: String, classId: String, topic: String): Result<ClassSession> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/lessons")
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
                    put("class_id", classId)
                    put("topic", topic)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val obj = JSONObject(response)
                    Result.success(
                        ClassSession(
                            classId = obj.optString("class_id", classId),
                            className = obj.optString("class_name", "Classe"),
                            subject = obj.optString("subject_name", obj.optString("subject", "Materia")),
                            hourSlot = obj.optString("hour_slot", "1a ora"),
                            isSigned = true,
                            lessonTopic = topic
                        )
                    )
                } else {
                    Result.failure(Exception("Errore firma lezione: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun submitRollCall(token: String, classId: String, records: List<StudentRollCall>): Result<Boolean> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/attendance")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "POST"
                    setRequestProperty("Authorization", "Bearer $token")
                    setRequestProperty("Content-Type", "application/json; charset=UTF-8")
                    doOutput = true
                    doInput = true
                    connectTimeout = 10000
                    readTimeout = 10000
                }

                val arr = JSONArray()
                records.forEach { r ->
                    arr.put(
                        JSONObject().apply {
                            put("student_id", r.studentId)
                            put("status", r.status)
                            put("reason", r.note)
                        }
                    )
                }

                val payload = JSONObject().apply {
                    put("class_id", classId)
                    put("records", arr)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    Result.success(true)
                } else {
                    Result.failure(Exception("Errore salvataggio presenze: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun submitGrade(token: String, proposal: GradeProposal): Result<Boolean> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/grades")
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
                    put("student_id", proposal.studentId)
                    put("subject_id", proposal.subjectId)
                    put("grade", proposal.grade)
                    put("weight", proposal.weight)
                    put("type", proposal.type)
                    put("notes", proposal.comment)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    Result.success(true)
                } else {
                    Result.failure(Exception("Errore inserimento voto: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }

    override suspend fun saveDeferredScrutiny(token: String, resolution: DeferredScrutinyResolution): Result<Boolean> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/scrutiny/deferred-resolution")
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
                    put("student_id", resolution.studentId)
                    put("subject_id", resolution.subjectId)
                    put("recovery_grade", resolution.recoveryGrade)
                    put("final_outcome", resolution.status)
                    put("deliberation_notes", resolution.finalDecision)
                }

                OutputStreamWriter(conn.outputStream).use { it.write(payload.toString()) }

                if (conn.responseCode in 200..299) {
                    Result.success(true)
                } else {
                    Result.failure(Exception("Errore scrutinio differito: HTTP ${conn.responseCode}"))
                }
            } catch (e: Exception) {
                Result.failure(e)
            }
        }
}
