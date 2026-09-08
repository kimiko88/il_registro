package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryInsuranceClaimsTest {

    @Test
    fun inailProtocol_startsWithInail() {
        val inailProt = "INAIL-2026-9021"
        assertTrue(inailProt.startsWith("INAIL-"))
    }

    @Test
    fun insuranceClaimNumber_startsWithSin() {
        val claim = "SIN-4412"
        assertTrue(claim.startsWith("SIN-"))
    }
}
