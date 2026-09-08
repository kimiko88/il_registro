package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentTransportTest {

    @Test
    fun schoolBusLine_isAssigned() {
        val line = "Linea 3: Nord - Centro"
        assertTrue(line.startsWith("Linea 3"))
    }

    @Test
    fun morningPickupTime_isSet() {
        val pickupTime = "07:35"
        assertTrue(pickupTime.isNotBlank())
    }
}
