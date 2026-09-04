package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test
import java.io.File
import javax.xml.parsers.DocumentBuilderFactory

class LocalizationIntegrityTest {

    @Test
    fun stringResources_existForAllTargetLocales() {
        val resDir = File("src/main/res")
        val expectedLocales = listOf(
            "values",
            "values-it",
            "values-en",
            "values-de",
            "values-fr",
            "values-es",
            "values-ar",
            "values-ro",
            "values-ru",
            "values-sq",
            "values-uk",
            "values-zh"
        )

        for (locale in expectedLocales) {
            val dir = File(resDir, locale)
            val stringsFile = File(dir, "strings.xml")
            if (resDir.exists()) {
                assertTrue("strings.xml missing in $locale", stringsFile.exists())
            }
        }
    }

    @Test
    fun stringResources_haveKeyParityAcrossLocales() {
        val resDir = File("src/main/res")
        if (!resDir.exists()) return

        val baseFile = File(resDir, "values/strings.xml")
        assertTrue("Base strings.xml must exist", baseFile.exists())
        val baseKeys = parseStringKeys(baseFile)
        assertTrue("Base strings must not be empty", baseKeys.isNotEmpty())

        val locales = listOf(
            "values-it", "values-en", "values-de", "values-fr", "values-es",
            "values-ar", "values-ro", "values-ru", "values-sq", "values-uk", "values-zh"
        )

        for (locale in locales) {
            val file = File(File(resDir, locale), "strings.xml")
            assertTrue("strings.xml missing in $locale", file.exists())
            val localeKeys = parseStringKeys(file)
            val missing = baseKeys - localeKeys
            assertTrue("Missing keys in $locale: $missing", missing.isEmpty())
        }
    }

    private fun parseStringKeys(file: File): Set<String> {
        val factory = DocumentBuilderFactory.newInstance()
        val builder = factory.newDocumentBuilder()
        val doc = builder.parse(file)
        val stringNodes = doc.getElementsByTagName("string")
        val keys = mutableSetOf<String>()
        for (i in 0 until stringNodes.length) {
            val node = stringNodes.item(i)
            val name = node.attributes?.getNamedItem("name")?.nodeValue
            if (name != null) {
                keys.add(name)
            }
        }
        return keys
    }
}
