package it.scuola.registro.student

import it.scuola.registro.student.security.MdmConfigManager
import org.junit.Assert.*
import org.junit.Test

class MdmConfigManagerTest {

    @Test
    fun mdmManager_handlesNullContextGracefully() {
        val manager = MdmConfigManager(null)
        val bundle = manager.applyMdmPolicies()
        assertNotNull(bundle)
        assertTrue(bundle.isEmpty)
    }
}
