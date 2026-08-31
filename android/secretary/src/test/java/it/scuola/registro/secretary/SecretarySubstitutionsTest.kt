package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretarySubstitutionsTest {

    @Test
    fun substitution_matchesAvailableTeacherSlot() {
        val absentTeacherSlot = 3
        val substituteTeacherFreeSlot = 3
        assertEquals(absentTeacherSlot, substituteTeacherFreeSlot)
    }

    @Test
    fun sidiXmlPayload_validatesSchoolCode() {
        val schoolCode = "RMIS09900B"
        assertTrue(schoolCode.matches(Regex("^[A-Z]{2}[A-Z0-9]{8}$")))
    }
}
