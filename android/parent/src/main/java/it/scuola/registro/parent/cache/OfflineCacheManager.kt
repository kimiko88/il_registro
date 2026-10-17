package it.scuola.registro.parent.cache

import it.scuola.registro.parent.data.ColloquioBooking
import it.scuola.registro.parent.data.ParentChild
import it.scuola.registro.parent.data.PendingAbsence

class OfflineCacheManager {
    private val cachedChildren = mutableListOf<ParentChild>()
    private val cachedAbsences = mutableListOf<PendingAbsence>()
    private val cachedBookings = mutableListOf<ColloquioBooking>()
    private var lastSyncTimestamp: Long = 0

    fun saveChildren(children: List<ParentChild>) {
        cachedChildren.clear()
        cachedChildren.addAll(children)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedChildren(): List<ParentChild> = cachedChildren.toList()

    fun saveAbsences(absences: List<PendingAbsence>) {
        cachedAbsences.clear()
        cachedAbsences.addAll(absences)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedAbsences(): List<PendingAbsence> = cachedAbsences.toList()

    fun saveBookings(bookings: List<ColloquioBooking>) {
        cachedBookings.clear()
        cachedBookings.addAll(bookings)
        lastSyncTimestamp = System.currentTimeMillis()
    }

    fun getCachedBookings(): List<ColloquioBooking> = cachedBookings.toList()

    fun isCacheValid(maxAgeMs: Long = 3600000): Boolean {
        return (System.currentTimeMillis() - lastSyncTimestamp) < maxAgeMs && cachedChildren.isNotEmpty()
    }

    fun clearCache() {
        cachedChildren.clear()
        cachedAbsences.clear()
        cachedBookings.clear()
        lastSyncTimestamp = 0
    }
}
