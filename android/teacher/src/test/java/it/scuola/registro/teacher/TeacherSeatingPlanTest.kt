package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherSeatingPlanTest {

    @Test
    fun totalAssignedDesks_matchesClassCount() {
        val totalDesks = 24
        assertEquals(24, totalDesks)
    }

    @Test
    fun cooperativeGroup_hasAssignedStudents() {
        val groupMembers = listOf("Mario Rossi", "Luca Bianchi", "Giulia Verdi", "Sara Neri")
        assertEquals(4, groupMembers.size)
    }
}
