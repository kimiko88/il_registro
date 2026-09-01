package it.scuola.registro.student.cache

import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.data.HomeworkAssignment

class OfflineCacheManager {
    private val cachedGrades = mutableListOf<GradeEntry>()
    private val cachedHomework = mutableListOf<HomeworkAssignment>()
    private var lastSyncTimestamp: Long = 0

    fun saveGrades(grades: List<GradeEntry>) {
        cachedGrades.clear()
        cachedGrades.addAll(grades)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedGrades(): List<GradeEntry> = cachedGrades.toList()

    fun saveHomework(homework: List<HomeworkAssignment>) {
        cachedHomework.clear()
        cachedHomework.addAll(homework)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedHomework(): List<HomeworkAssignment> = cachedHomework.toList()

    fun isCacheValid(maxAgeMs: Long = 3600000): Boolean {
        return (System.currentTimeMillis() - lastSyncTimestamp) < maxAgeMs && cachedGrades.isNotEmpty()
    }

    fun clearCache() {
        cachedGrades.clear()
        cachedHomework.clear()
        lastSyncTimestamp = 0
    }
}
