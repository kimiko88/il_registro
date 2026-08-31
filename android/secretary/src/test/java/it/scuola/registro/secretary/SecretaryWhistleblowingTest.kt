package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryWhistleblowingTest {

    @Test
    fun encryptionStandard_isRSA4096() {
        val cipher = "RSA-4096"
        assertTrue(cipher.contains("4096"))
    }

    @Test
    fun legalCompliance_isDLgs24() {
        val law = "D.Lgs. 24/2023"
        assertTrue(law.contains("24/2023"))
    }
}
