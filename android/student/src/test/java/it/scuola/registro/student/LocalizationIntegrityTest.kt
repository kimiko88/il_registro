package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test
import java.io.File

class LocalizationIntegrityTest {

    @Test
    fun stringResources_existForAllTargetLocales() {
        val resDir = File("src/main/res")
        val expectedLocales = listOf(
            "values", // it (default)
            "values-en", // en
            "values-de", // de
            "values-fr", // fr
            "values-es", // es
            "values-ar", // ar
            "values-ro", // ro
            "values-ru", // ru
            "values-sq", // sq
            "values-uk", // uk
            "values-zh"  // zh
        )

        for (locale in expectedLocales) {
            val dir = File(resDir, locale)
            val stringsFile = File(dir, "strings.xml")
            if (resDir.exists()) {
                assertTrue("strings.xml missing in $locale", stringsFile.exists())
            }
        }
    }
}
