package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretarySafetyTest {

    @Test
    fun evacuationDrillTime_isUnderTarget() {
        val actualSeconds = 165 // 2 min 45 sec
        val targetSeconds = 180 // 3 min
        assertTrue(actualSeconds < targetSeconds)
    }

    @Test
    fun totalEvacuatedCount_matchesSum() {
        val students = 520
        val staff = 68
        val total = students + staff
        assertEquals(588, total)
    }
}
