import { ref, effectScope, nextTick } from 'vue'
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { useDraftAutosave } from 'src/composables/useDraftAutosave'

describe('useDraftAutosave composable', () => {
    let scope

    beforeEach(() => {
        localStorage.clear()
        vi.useFakeTimers()
        scope = effectScope()
    })

    afterEach(() => {
        scope.stop()
        localStorage.clear()
        vi.useRealTimers()
    })

    it('should detect when no draft exists initially', () => {
        scope.run(() => {
            const formRef = ref({ title: 'Nota Lezione', content: '' })
            const { hasDraft } = useDraftAutosave('test-key', formRef)

            expect(hasDraft.value).toBe(false)
        })
    })

    it('should auto-save draft to localStorage after debounce timer', async () => {
        let hasDraftRef, savedAtRef, formRef

        scope.run(() => {
            formRef = ref({ title: 'Nota Lezione', content: '' })
            const { hasDraft, savedAt } = useDraftAutosave('test-key', formRef, 500)
            hasDraftRef = hasDraft
            savedAtRef = savedAt
        })

        // Modify reactive value
        formRef.value.content = 'Spiegazione derivate'
        await nextTick()

        // Fast-forward timer by 500ms
        vi.advanceTimersByTime(500)

        expect(hasDraftRef.value).toBe(true)
        expect(savedAtRef.value).not.toBeNull()

        const stored = JSON.parse(localStorage.getItem('il_registro_draft_test-key'))
        expect(stored.data.content).toBe('Spiegazione derivate')
    })

    it('should restore draft into ref data', () => {
        const initialDraft = {
            data: { title: 'Compito Recuperato', content: 'Pagina 45 n. 12' },
            savedAt: new Date().toISOString()
        }
        localStorage.setItem('il_registro_draft_test-key', JSON.stringify(initialDraft))

        scope.run(() => {
            const formRef = ref({ title: '', content: '' })
            const { hasDraft, restoreDraft } = useDraftAutosave('test-key', formRef)

            expect(hasDraft.value).toBe(true)

            const restored = restoreDraft()
            expect(restored).toBe(true)
            expect(formRef.value.title).toBe('Compito Recuperato')
            expect(formRef.value.content).toBe('Pagina 45 n. 12')
        })
    })

    it('should clear draft on demand', async () => {
        let hasDraftRef, clearDraftFn, formRef

        scope.run(() => {
            formRef = ref({ title: 'Bozza da cancellare', content: 'Testo' })
            const { hasDraft, clearDraft } = useDraftAutosave('test-key', formRef, 100)
            hasDraftRef = hasDraft
            clearDraftFn = clearDraft
        })

        formRef.value.content = 'Aggiornamento'
        await nextTick()
        vi.advanceTimersByTime(100)
        expect(hasDraftRef.value).toBe(true)

        clearDraftFn()

        expect(hasDraftRef.value).toBe(false)
        expect(localStorage.getItem('il_registro_draft_test-key')).toBeNull()
    })
})
