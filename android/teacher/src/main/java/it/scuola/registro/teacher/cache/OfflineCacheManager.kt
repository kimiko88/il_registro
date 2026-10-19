package it.scuola.registro.teacher.cache

import it.scuola.registro.teacher.data.ClassSession
import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall

class OfflineCacheManager {
    private val cachedClasses = mutableListOf<ClassSession>()
    private val cachedRollCalls = mutableListOf<StudentRollCall>()
    private val cachedGrades = mutableListOf<GradeProposal>()
    private val cachedDeferredScrutiny = mutableListOf<DeferredScrutinyResolution>()
    private var lastSyncTimestamp: Long = 0

    fun saveClasses(classes: List<ClassSession>) {
        cachedClasses.clear()
        cachedClasses.addAll(classes)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedClasses(): List<ClassSession> = cachedClasses.toList()

    fun saveRollCalls(rollCalls: List<StudentRollCall>) {
        cachedRollCalls.clear()
        cachedRollCalls.addAll(rollCalls)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedRollCalls(): List<StudentRollCall> = cachedRollCalls.toList()

    fun saveGrades(grades: List<GradeProposal>) {
        cachedGrades.clear()
        cachedGrades.addAll(grades)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedGrades(): List<GradeProposal> = cachedGrades.toList()

    fun saveDeferredScrutiny(resolutions: List<DeferredScrutinyResolution>) {
        cachedDeferredScrutiny.clear()
        cachedDeferredScrutiny.addAll(resolutions)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedDeferredScrutiny(): List<DeferredScrutinyResolution> = cachedDeferredScrutiny.toList()

    fun isCacheValid(maxAgeMs: Long = 3600000): Boolean {
        return (System.currentTimeMillis() - lastSyncTimestamp) < maxAgeMs && cachedClasses.isNotEmpty()
    }

    fun clearCache() {
        cachedClasses.clear()
        cachedRollCalls.clear()
        cachedGrades.clear()
        cachedDeferredScrutiny.clear()
        lastSyncTimestamp = 0
    }
}
