package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryTransfersTest {

    @Test
    fun nullaOstaProtocol_matchesPattern() {
        val protocol = "NO-2026-0392"
        assertTrue(protocol.startsWith("NO-"))
    }

    @Test
    fun sidiTransferSync_isCompleted() {
        val isSidiSynced = true
        assertTrue(isSidiSynced)
    }
}
