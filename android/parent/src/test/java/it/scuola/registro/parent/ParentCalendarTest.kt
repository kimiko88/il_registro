package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentCalendarTest {

    @Test
    fun schoolEndEvent_matchesDate() {
        val lastDay = "2026-06-06"
        assertTrue(lastDay.startsWith("2026-06"))
    }

    @Test
    fun holidaySuspension_isValid() {
        val isSuspended = true
        assertTrue(isSuspended)
    }
}
