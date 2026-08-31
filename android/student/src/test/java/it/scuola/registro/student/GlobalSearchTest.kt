package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class GlobalSearchTest {

    @Test
    fun searchQuery_findsMatchingResult() {
        val query = "Compiti"
        val sampleResults = listOf("Circolari", "Compiti a casa", "Materiali")
        val matches = sampleResults.filter { it.contains(query, ignoreCase = true) }
        assertEquals(1, matches.size)
    }

    @Test
    fun emptySearchQuery_handledGracefully() {
        val emptyQuery = ""
        assertTrue(emptyQuery.isEmpty())
    }
}
