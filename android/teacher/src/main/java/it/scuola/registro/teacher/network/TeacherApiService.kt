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
    suspend fun signLesson(token: String, classId: String, topic: String): Result<ClassSession>
    suspend fun submitRollCall(token: String, classId: String, records: List<StudentRollCall>): Result<Boolean>
    suspend fun submitGrade(token: String, proposal: GradeProposal): Result<Boolean>
    suspend fun saveDeferredScrutiny(token: String, resolution: DeferredScrutinyResolution): Result<Boolean>
}

class HttpTeacherApiService(
    private val baseUrl: String = "https://api.scuola.registro.it/api/v1"
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

    override suspend fun getClasses(token: String): Result<List<ClassSession>> =
        withContext(Dispatchers.IO) {
            try {
                val url = URL("$baseUrl/teacher/classes")
                val conn = (url.openConnection() as HttpURLConnection).apply {
                    requestMethod = "GET"
                    setRequestProperty("Authorization", "Bearer $token")
                    connectTimeout = 10000
                    readTimeout = 10000
                }
                if (conn.responseCode in 200..299) {
                    val response = BufferedReader(InputStreamReader(conn.inputStream)).use { it.readText() }
                    val arr = JSONArray(response)
                    val list = mutableListOf<ClassSession>()
                    for (i in 0 until arr.length()) {
                        val obj = arr.getJSONObject(i)
                        list.add(
                            ClassSession(
                                classId = obj.optString("id", "$i"),
                                className = obj.optString("name", obj.optString("class_name", "Classe")),
                                subject = obj.optString("subject_name", obj.optString("subject", "Materia")),
                                hourSlot = "Orario lezioni",
                                isSigned = false,
                                lessonTopic = ""
                            )
                        )
                    }
                    Result.success(list)
                } else {
                    Result.failure(Exception("Errore caricamento classi: HTTP ${conn.responseCode}"))
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
                    put("final_outcome", resolution.finalOutcome)
                    put("deliberation_notes", resolution.deliberationNotes)
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
