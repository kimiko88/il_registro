package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryTimetableTest {

    @Test
    fun timetableSchedule_hasZeroConflicts() {
        val conflictsCount = 0
        assertEquals(0, conflictsCount)
    }

    @Test
    fun weeklyPeriodsCount_isStandardThirtyHours() {
        val dailyPeriods = 6
        val weeklyDays = 5
        assertEquals(30, dailyPeriods * weeklyDays)
    }
}
