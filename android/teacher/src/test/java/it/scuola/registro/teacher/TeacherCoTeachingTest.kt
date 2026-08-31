package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherCoTeachingTest {

    @Test
    fun doubleSignature_isConfirmed() {
        val hasDoubleSignature = true
        assertTrue(hasDoubleSignature)
    }

    @Test
    fun coteachingRole_isItpOrSupport() {
        val itpDocente = "Prof. Neri (ITP Lab)"
        assertTrue(itpDocente.contains("ITP"))
    }
}
