package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentGeneralMeetingsTest {

    @Test
    fun queuePosition_isPositiveInteger() {
        val position = 3
        assertTrue(position > 0)
    }

    @Test
    fun estimatedWaitTime_isCalculated() {
        val parentsAhead = 2
        val averageMinutesPerParent = 4
        val estimatedWait = parentsAhead * averageMinutesPerParent
        assertEquals(8, estimatedWait)
    }
}
