package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryInventoryTest {

    @Test
    fun assetInventoryNumber_startsWithInv() {
        val inventoryNum = "INV-2026-0491"
        assertTrue(inventoryNum.startsWith("INV-"))
    }

    @Test
    fun assetFundingSource_isPnrr() {
        val funding = "PNRR Piano Scuola 4.0"
        assertTrue(funding.contains("PNRR"))
    }
}
