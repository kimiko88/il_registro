import { watch, onBeforeUnmount, getCurrentInstance, ref, computed, toRaw } from 'vue'

const DRAFT_PREFIX = 'il_registro_draft_'
/** Bozze più vecchie di 24 ore vengono scartate automaticamente al restore. */
const DRAFT_TTL_MS = 24 * 60 * 60 * 1000

export function useDraftAutosave(storageKey, dataRef, debounceMs = 1500) {
    const key = DRAFT_PREFIX + storageKey
    const savedAt = ref(null)
    let debounceTimer = null

    // Check if a non-expired draft exists on init
    const hasDraft = ref(false)
    try {
        const raw = localStorage.getItem(key)
        if (raw) {
            const parsed = JSON.parse(raw)
            let isExpired = false
            if (parsed?.savedAt) {
                const savedTime = new Date(parsed.savedAt).getTime()
                if (!isNaN(savedTime) && (Date.now() - savedTime > DRAFT_TTL_MS)) {
                    isExpired = true
                }
            }
            hasDraft.value = !isExpired
            if (isExpired) {
                localStorage.removeItem(key)
            }
        }
    } catch {
        hasDraft.value = false
    }

    /** Human-readable age of the saved draft, e.g. "5 minuti fa" */
    const draftAge = computed(() => {
        if (!savedAt.value) return null
        const diffMs = Date.now() - savedAt.value.getTime()
        const diffMins = Math.floor(diffMs / 60_000)
        if (diffMins < 1) return 'pochi secondi fa'
        if (diffMins < 60) return `${diffMins} minut${diffMins === 1 ? 'o' : 'i'} fa`
        const diffHours = Math.floor(diffMins / 60)
        return `${diffHours} or${diffHours === 1 ? 'a' : 'e'} fa`
    })

    // Restore draft data into the dataRef
    function restoreDraft() {
        try {
            const raw = localStorage.getItem(key)
            if (!raw) return false
            const saved = JSON.parse(raw)
            // Discard drafts older than DRAFT_TTL_MS to avoid surfacing stale data.
            if (saved?.savedAt) {
                const age = Date.now() - new Date(saved.savedAt).getTime()
                if (age > DRAFT_TTL_MS) {
                    localStorage.removeItem(key)
                    hasDraft.value = false
                    return false
                }
                savedAt.value = new Date(saved.savedAt)
            }
            if (saved && saved.data !== undefined) {
                if (typeof dataRef.value === 'object' && dataRef.value !== null) {
                    Object.assign(dataRef.value, saved.data)
                } else {
                    dataRef.value = saved.data
                }
                return true
            }
        } catch (err) {
            console.warn('[useDraftAutosave] restore failed:', err)
        }
        return false
    }

    // Clear the saved draft (call after successful API save)
    function clearDraft() {
        try {
            localStorage.removeItem(key)
            hasDraft.value = false
            savedAt.value = null
        } catch { /* noop */ }
    }

    // Internal save function
    function saveDraft(value) {
        try {
            localStorage.setItem(key, JSON.stringify({
                data: value,
                savedAt: new Date().toISOString()
            }))
            savedAt.value = new Date()
            hasDraft.value = true
        } catch (err) {
            // LocalStorage can be full or unavailable
            console.warn('[useDraftAutosave] save failed:', err)
        }
    }

    // Watch for changes with debounce
    const stopWatcher = watch(
        dataRef,
        (newVal) => {
            clearTimeout(debounceTimer)
            debounceTimer = setTimeout(() => {
                // Efficient clone using structuredClone / toRaw to avoid heavy JSON stringify
                try {
                    const rawVal = toRaw(newVal)
                    const serializable = typeof structuredClone === 'function' ? structuredClone(rawVal) : JSON.parse(JSON.stringify(rawVal))
                    saveDraft(serializable)
                } catch {
                    saveDraft(newVal)
                }
            }, debounceMs)
        },
        { deep: true }
    )

    function stop() {
        if (stopWatcher) stopWatcher()
        if (debounceTimer) {
            clearTimeout(debounceTimer)
            debounceTimer = null
        }
    }

    // Clean up on unmount if within component setup
    if (getCurrentInstance()) {
        onBeforeUnmount(() => {
            stop()
        })
    }

    return { hasDraft, draftAge, restoreDraft, clearDraft, savedAt, stop }
}
