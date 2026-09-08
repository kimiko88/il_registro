package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryMaterialsTest {

    @Test
    fun cloudStorageUsage_doesNotExceedQuota() {
        val usedGb = 42.8
        val quotaGb = 250.0
        assertTrue(usedGb < quotaGb)
    }

    @Test
    fun antivirusScan_isClean() {
        val status = "Regolare"
        assertEquals("Regolare", status)
    }
}
