package it.scuola.registro.secretary.cache

import it.scuola.registro.secretary.data.AuditLogEntry
import it.scuola.registro.secretary.data.CertificateRequest
import it.scuola.registro.secretary.data.DashboardStats
import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.data.ScrutinyClassStatus

class OfflineCacheManager {
    private val cachedUsers = mutableListOf<ManagedUser>()
    private val cachedCertificates = mutableListOf<CertificateRequest>()
    private val cachedScrutiny = mutableListOf<ScrutinyClassStatus>()
    private val cachedAuditLogs = mutableListOf<AuditLogEntry>()
    private var cachedStats: DashboardStats? = null
    private var lastSyncTimestamp: Long = 0

    fun saveUsers(users: List<ManagedUser>) {
        cachedUsers.clear()
        cachedUsers.addAll(users)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedUsers(): List<ManagedUser> = cachedUsers.toList()

    fun saveCertificates(certificates: List<CertificateRequest>) {
        cachedCertificates.clear()
        cachedCertificates.addAll(certificates)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedCertificates(): List<CertificateRequest> = cachedCertificates.toList()

    fun saveScrutiny(scrutinyList: List<ScrutinyClassStatus>) {
        cachedScrutiny.clear()
        cachedScrutiny.addAll(scrutinyList)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedScrutiny(): List<ScrutinyClassStatus> = cachedScrutiny.toList()

    fun saveAuditLogs(logs: List<AuditLogEntry>) {
        cachedAuditLogs.clear()
        cachedAuditLogs.addAll(logs)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedAuditLogs(): List<AuditLogEntry> = cachedAuditLogs.toList()

    fun saveDashboardStats(stats: DashboardStats) {
        cachedStats = stats
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedDashboardStats(): DashboardStats? = cachedStats

    fun isCacheValid(maxAgeMs: Long = 3600000): Boolean {
        return (System.currentTimeMillis() - lastSyncTimestamp) < maxAgeMs && (cachedUsers.isNotEmpty() || cachedCertificates.isNotEmpty())
    }

    fun clearCache() {
        cachedUsers.clear()
        cachedCertificates.clear()
        cachedScrutiny.clear()
        cachedAuditLogs.clear()
        cachedStats = null
        lastSyncTimestamp = 0
    }
}
