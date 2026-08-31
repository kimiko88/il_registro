package it.scuola.registro.teacher.viewmodel

import it.scuola.registro.teacher.data.ClassSession
import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall
import it.scuola.registro.teacher.network.HttpTeacherApiService
import it.scuola.registro.teacher.network.TeacherApiService

class TeacherViewModel(
    private val apiService: TeacherApiService = HttpTeacherApiService()
) {
    var activeSession: ClassSession? = null
        private set

    var currentSession: ClassSession? = null
        private set

    var rollCallList = mutableListOf<StudentRollCall>()
        private set

    var gradesList = mutableListOf<GradeProposal>()
        private set

    var deferredResolutions = mutableListOf<DeferredScrutinyResolution>()
        private set

    var isLoading: Boolean = false
        private set

    var errorMessage: String? = null
        private set

    suspend fun signLessonViaApi(token: String, classId: String, topic: String): Boolean {
        isLoading = true
        errorMessage = null
        val result = apiService.signLesson(token, classId, topic)
        isLoading = false
        return if (result.isSuccess) {
            activeSession = result.getOrNull()
            currentSession = activeSession
            true
        } else {
            errorMessage = result.exceptionOrNull()?.localizedMessage
            false
        }
    }

    suspend fun submitRollCallViaApi(token: String, classId: String, records: List<StudentRollCall>): Boolean {
        isLoading = true
        errorMessage = null
        val result = apiService.submitRollCall(token, classId, records)
        isLoading = false
        return if (result.isSuccess) {
            true
        } else {
            errorMessage = result.exceptionOrNull()?.localizedMessage
            false
        }
    }

    suspend fun submitGradeViaApi(token: String, proposal: GradeProposal): Boolean {
        isLoading = true
        errorMessage = null
        val result = apiService.submitGrade(token, proposal)
        isLoading = false
        return if (result.isSuccess) {
            gradesList.add(proposal)
            true
        } else {
            errorMessage = result.exceptionOrNull()?.localizedMessage
            false
        }
    }

    fun loadSessionData() {
        activeSession = ClassSession(
            classId = "3A",
            className = "Classe 3A",
            subject = "Matematica",
            hourSlot = "1a e 2a ora (08:00 - 10:00)",
            isSigned = false,
            lessonTopic = ""
        )
        currentSession = activeSession

        rollCallList = mutableListOf(
            StudentRollCall("s1", "Banchi Andrea", "presente"),
            StudentRollCall("s2", "Bianchi Elena", "presente"),
            StudentRollCall("s3", "Ferrari Matteo", "assente"),
            StudentRollCall("s4", "Rossi Sofia", "ritardo")
        )
    }

    fun loadSampleSession() = loadSessionData()

    fun signLessonHour(topic: String): Boolean {
        if (topic.isBlank()) return false
        activeSession?.let {
            it.isSigned = true
            it.lessonTopic = topic
            currentSession = it
            return true
        }
        return false
    }

    fun signCurrentLesson(topic: String): Boolean = signLessonHour(topic)

    fun updateStudentStatus(studentId: String, status: String): Boolean {
        val student = rollCallList.find { it.studentId == studentId }
        if (student != null) {
            student.status = status
            return true
        }
        return false
    }

    fun toggleStudentAttendance(studentId: String, status: String): Boolean = updateStudentStatus(studentId, status)

    fun addGrade(proposal: GradeProposal): Boolean {
        if (proposal.grade < 1.0 || proposal.grade > 10.0) return false
        gradesList.add(proposal)
        return true
    }

    fun submitGradeForStudent(studentId: String, grade: Double, weight: Double, type: String, comment: String): Boolean {
        return addGrade(GradeProposal(studentId, "mat1", grade, weight, type, comment))
    }

    fun deliberateDeferredScrutiny(resolution: DeferredScrutinyResolution): Boolean {
        if (resolution.recoveryGrade < 1.0 || resolution.recoveryGrade > 10.0) return false
        deferredResolutions.add(resolution)
        return true
    }

    fun resolveDeferredStudent(studentId: String, recoveryGrade: Double, outcome: String, notes: String): Boolean {
        return deliberateDeferredScrutiny(DeferredScrutinyResolution(studentId, "mat1", recoveryGrade, outcome, notes))
    }
}
