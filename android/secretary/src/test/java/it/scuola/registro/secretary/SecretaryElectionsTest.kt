package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryElectionsTest {

    @Test
    fun voterTurnout_isAboveFiftyPercent() {
        val totalVoters = 1230.0
        val castVotes = 842.0
        val turnout = (castVotes / totalVoters) * 100.0
        assertTrue(turnout > 50.0)
    }

    @Test
    fun pollingStationStatus_isOpen() {
        val status = "Seggio Aperto"
        assertEquals("Seggio Aperto", status)
    }
}
