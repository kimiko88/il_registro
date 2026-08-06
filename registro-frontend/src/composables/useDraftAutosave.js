import { watch, onBeforeUnmount } from 'vue'

/**
 * useDraftAutosave — Composable per il salvataggio automatico delle bozze in localStorage.
 *
 * Salva automaticamente un ref/reactive dopo `debounceMs` millisecondi di inattività.
 * Ripristina automaticamente la bozza salvata al montaggio del componente (se presente).
 * Cancella la bozza dopo un salvataggio riuscito (`clearDraft()`).
 *
 * @param {string}  storageKey   — Chiave localStorage univoca (es. 'agenda-event-draft-{classId}')
 * @param {Ref}     dataRef      — Ref o reactive da sorvegliare e salvare
 * @param {number}  debounceMs   — Debounce in ms prima di scrivere nel localStorage (default 1500)
 *
 * @returns {object}
 *   hasDraft    — boolean: true se esiste una bozza salvata
 *   restoreDraft — funzione: ripristina la bozza nel dataRef
 *   clearDraft  — funzione: elimina la bozza salvata (da chiamare dopo salvataggio riuscito)
 *   savedAt     — Ref<Date|null>: timestamp dell'ultimo autosave
 *
 * Usage:
 *   const { hasDraft, restoreDraft, clearDraft, savedAt } = useDraftAutosave(
 *     `lesson-draft-${classId}`,
 *     lessonForm
 *   )
 *   // On mount:
 *   if (hasDraft.value) {
 *     const confirmed = await $q.dialog({ message: 'Esiste una bozza non salvata. Ripristinarla?' })
 *     if (confirmed) restoreDraft()
 *   }
 *   // On successful save:
 *   clearDraft()
 */
import { ref } from 'vue'

const DRAFT_PREFIX = 'registrov2_draft_'

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

    // Clean up on unmount
    onBeforeUnmount(() => {
        stopWatcher()
        clearTimeout(debounceTimer)
    })

    return { hasDraft, restoreDraft, clearDraft, savedAt }
}
