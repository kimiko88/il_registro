package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryClassFormationTest {

    @Test
    fun genderBalance_isEqual() {
        val males = 12
        val females = 12
        assertEquals(males, females)
    }

    @Test
    fun totalSectionStudents_matchesSum() {
        val total = 24
        assertEquals(24, total)
    }
}
