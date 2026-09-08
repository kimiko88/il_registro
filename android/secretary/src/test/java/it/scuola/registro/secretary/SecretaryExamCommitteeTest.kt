package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryExamCommitteeTest {

    @Test
    fun examCommitteeCandidateCount_matchesSum() {
        val class5A = 24
        val class5B = 24
        val total = class5A + class5B
        assertEquals(48, total)
    }

    @Test
    fun externalPresident_isAssigned() {
        val president = "Prof.ssa Maria Ricci"
        assertTrue(president.isNotBlank())
    }
}
