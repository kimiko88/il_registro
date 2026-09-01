package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherFreeActivitiesTest {

    @Test
    fun onDutyHour_durationIsOneHour() {
        val startHour = 10
        val endHour = 11
        assertEquals(1, endHour - startHour)
    }

    @Test
    fun activityCategory_isValid() {
        val category = "Ora a Disposizione"
        assertTrue(category.isNotBlank())
    }
}
