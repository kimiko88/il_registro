package it.scuola.registro.teacher.viewmodel

import it.scuola.registro.teacher.data.ClassSession
import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall

class TeacherViewModel {
    var activeSession: ClassSession? = null
        private set

    var rollCallList = mutableListOf<StudentRollCall>()
        private set

    var gradesList = mutableListOf<GradeProposal>()
        private set

    var deferredResolutions = mutableListOf<DeferredScrutinyResolution>()
        private set

    fun loadSessionData() {
        activeSession = ClassSession(
            classId = "3A",
            className = "Classe 3A",
            subject = "Matematica",
            hourSlot = "1a e 2a ora (08:00 - 10:00)",
            isSigned = false,
            lessonTopic = ""
        )

        rollCallList = mutableListOf(
            StudentRollCall("s1", "Banchi Andrea", "presente"),
            StudentRollCall("s2", "Bianchi Elena", "presente"),
            StudentRollCall("s3", "Ferrari Matteo", "assente"),
            StudentRollCall("s4", "Rossi Sofia", "ritardo")
        )
    }

    fun signLessonHour(topic: String): Boolean {
        if (topic.isBlank()) return false
        activeSession?.let {
            it.isSigned = true
            it.lessonTopic = topic
            return true
        }
        return false
    }

    fun updateStudentStatus(studentId: String, status: String): Boolean {
        val student = rollCallList.find { it.studentId == studentId }
        if (student != null) {
            student.status = status
            return true
        }
        return false
    }

    fun addGrade(proposal: GradeProposal): Boolean {
        if (proposal.grade < 1.0 || proposal.grade > 10.0) return false
        gradesList.add(proposal)
        return true
    }

    fun deliberateDeferredScrutiny(resolution: DeferredScrutinyResolution): Boolean {
        if (resolution.recoveryGrade < 1.0 || resolution.recoveryGrade > 10.0) return false
        deferredResolutions.add(resolution)
        return true
    }
}
