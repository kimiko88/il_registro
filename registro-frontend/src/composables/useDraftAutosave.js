import { watch, onBeforeUnmount, getCurrentInstance, ref } from 'vue'

const DRAFT_PREFIX = 'il_registro_draft_'

export function useDraftAutosave(storageKey, dataRef, debounceMs = 1500) {
    const key = DRAFT_PREFIX + storageKey
    const savedAt = ref(null)
    let debounceTimer = null

    // Check if a draft exists on init
    const hasDraft = ref(false)
    try {
        const raw = localStorage.getItem(key)
        hasDraft.value = !!raw
    } catch {
        hasDraft.value = false
    }

    // Restore draft data into the dataRef
    function restoreDraft() {
        try {
            const raw = localStorage.getItem(key)
            if (!raw) return false
            const saved = JSON.parse(raw)
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
                // Deep-clone to avoid storing reactive proxies
                try {
                    const serializable = JSON.parse(JSON.stringify(newVal))
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

    return { hasDraft, restoreDraft, clearDraft, savedAt, stop }
}
