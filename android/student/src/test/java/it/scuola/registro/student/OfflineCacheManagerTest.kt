package it.scuola.registro.student

import it.scuola.registro.student.cache.OfflineCacheManager
import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.data.HomeworkAssignment
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class OfflineCacheManagerTest {

    private lateinit var cacheManager: OfflineCacheManager

    @Before
    fun setUp() {
        cacheManager = OfflineCacheManager()
    }

    @Test
    fun saveAndRetrieveGrades_worksCorrectly() {
        val grades = listOf(
            GradeEntry("1", "Matematica", 9.0, 1.0, "Scritto", "2026-08-30", 1)
        )
        cacheManager.saveGrades(grades)
        val retrieved = cacheManager.getCachedGrades()
        assertEquals(1, retrieved.size)
        assertEquals("Matematica", retrieved[0].subject)
        assertTrue(cacheManager.isCacheValid())
    }

    @Test
    fun saveAndRetrieveHomework_worksCorrectly() {
        val homework = listOf(
            HomeworkAssignment("1", "Fisica", "Esercizi", "Moto", "2026-09-02", false)
        )
        cacheManager.saveHomework(homework)
        val retrieved = cacheManager.getCachedHomework()
        assertEquals(1, retrieved.size)
        assertEquals("Fisica", retrieved[0].subject)
    }

    @Test
    fun clearCache_emptiesStoredData() {
        cacheManager.saveGrades(listOf(GradeEntry("1", "Storia", 8.0, 1.0, "Orale", "2026-08-20", 1)))
        cacheManager.clearCache()
        assertTrue(cacheManager.getCachedGrades().isEmpty())
        assertFalse(cacheManager.isCacheValid())
    }
}
