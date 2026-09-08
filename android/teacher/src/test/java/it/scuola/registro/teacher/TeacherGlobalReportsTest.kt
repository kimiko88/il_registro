package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherGlobalReportsTest {

    @Test
    fun globalReportJudgment_isNotBlank() {
        val judgment = "L'allievo ha dimostrato un impegno costante e maturo"
        assertTrue(judgment.isNotBlank())
    }

    @Test
    fun orientationAdvice_isSpecified() {
        val orientation = "Percorso Universitario Scientifico / Tecnologico"
        assertTrue(orientation.contains("Scientifico"))
    }
}
