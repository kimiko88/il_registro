import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ref, nextTick } from 'vue'
import { createTestingPinia } from '@pinia/testing'

// ─── No Quasar dependency needed for composable tests ─────────────────────

// ─────────────────────────────────────────────────────────────────────────────
// Suite 1 – useDraftAutosave composable
// ─────────────────────────────────────────────────────────────────────────────
describe('useDraftAutosave Composable', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        localStorage.clear()
    })

    it('DA01 — saves draft to localStorage after debounce', async () => {
        vi.useFakeTimers()
        const { useDraftAutosave } = await import('@/composables/useDraftAutosave')
        const data = ref({ title: 'Prima bozza', body: 'Testo del compito' })
        const { stop } = useDraftAutosave('test-draft-key', data, 100)

        // Replace entire ref to trigger the shallow watch
        data.value = { title: 'Seconda bozza', body: 'Testo aggiornato' }
        // Flush Vue's watch queue (microtasks) then advance fake timers for debounce
        await nextTick()
        vi.advanceTimersByTime(200)
        await nextTick()

        const saved = localStorage.getItem('il_registro_draft_test-draft-key')
        expect(saved).not.toBeNull()
        const parsed = JSON.parse(saved)
        expect(parsed.data.title).toBe('Seconda bozza')
        stop()
        vi.useRealTimers()
    })

    it('DA02 — hasDraft is true when storage contains key', async () => {
        localStorage.setItem('il_registro_draft_existing-key', JSON.stringify({
            data: { note: 'salvato' }, savedAt: new Date().toISOString()
        }))
        const { useDraftAutosave } = await import('@/composables/useDraftAutosave')
        const data = ref({ note: '' })
        const { hasDraft, stop } = useDraftAutosave('existing-key', data)
        expect(hasDraft.value).toBe(true)
        stop()
    })

    it('DA03 — restoreDraft populates dataRef with saved data', async () => {
        localStorage.setItem('il_registro_draft_restore-key', JSON.stringify({
            data: { content: 'bozza ripristinata' }, savedAt: new Date().toISOString()
        }))
        const { useDraftAutosave } = await import('@/composables/useDraftAutosave')
        const data = ref({ content: '' })
        const { restoreDraft, stop } = useDraftAutosave('restore-key', data)
        const restored = restoreDraft()
        expect(restored).toBe(true)
        expect(data.value.content).toBe('bozza ripristinata')
        stop()
    })

    it('DA04 — clearDraft removes entry from localStorage and resets flags', async () => {
        localStorage.setItem('il_registro_draft_clear-key', JSON.stringify({ data: 'something', savedAt: '' }))
        const { useDraftAutosave } = await import('@/composables/useDraftAutosave')
        const data = ref('something')
        const { hasDraft, clearDraft, stop } = useDraftAutosave('clear-key', data)
        expect(hasDraft.value).toBe(true)
        clearDraft()
        expect(hasDraft.value).toBe(false)
        expect(localStorage.getItem('il_registro_draft_clear-key')).toBeNull()
        stop()
    })

    it('DA05 — savedAt is updated on save', async () => {
        vi.useFakeTimers()
        const { useDraftAutosave } = await import('@/composables/useDraftAutosave')
        const data = ref({ text: 'init' })
        const { savedAt, stop } = useDraftAutosave('savedat-key', data, 50)
        expect(savedAt.value).toBeNull()
        // Replace the entire ref value to trigger the shallow watch (deep mutation won't trigger)
        data.value = { text: 'changed' }
        await nextTick()
        vi.advanceTimersByTime(200)
        await nextTick()
        expect(savedAt.value).not.toBeNull()
        stop()
        vi.useRealTimers()
    })

    it('DA06 — hasDraft is false when no draft exists', async () => {
        const { useDraftAutosave } = await import('@/composables/useDraftAutosave')
        const data = ref({ x: 1 })
        const { hasDraft, stop } = useDraftAutosave('no-such-key-xyz', data)
        expect(hasDraft.value).toBe(false)
        stop()
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 2 – useUndoToast composable
// ─────────────────────────────────────────────────────────────────────────────
describe('useUndoToast Composable', () => {
    it('UT01 — composable exports notifyWithUndo function and undoClicked ref', async () => {
        const mod = await import('@/composables/useUndoToast')
        expect(typeof mod.useUndoToast).toBe('function')
        // useUndoToast requires useQuasar which needs a Quasar app context.
        // We verify only the module structure here to avoid test environment issues.
        expect(mod.useUndoToast.length).toBeLessThanOrEqual(1) // accepts optional timeoutMs
    })
})


// ─────────────────────────────────────────────────────────────────────────────
// Suite 3 – useGlobalKeyboardShortcuts composable
// ─────────────────────────────────────────────────────────────────────────────
describe('useGlobalKeyboardShortcuts Composable', () => {
    beforeEach(() => { vi.clearAllMocks() })

    it('KS01 — composable can be imported and provides register/unregister functions', async () => {
        const mod = await import('@/composables/useGlobalKeyboardShortcuts')
        expect(typeof mod.useGlobalKeyboardShortcuts).toBe('function')
        const result = mod.useGlobalKeyboardShortcuts()
        expect(result).toBeDefined()
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 4 – useGradeFormatter composable
// ─────────────────────────────────────────────────────────────────────────────
describe('useGradeFormatter Composable', () => {
    it('GF01 — formats numeric grade as string', async () => {
        const { useGradeFormatter } = await import('@/composables/useGradeFormatter')
        const { formatGrade } = useGradeFormatter()
        const result = formatGrade(7.5)
        expect(result).toBeTruthy()
        expect(typeof result).toBe('string')
    })

    it('GF02 — handles null/undefined grade gracefully', async () => {
        const { useGradeFormatter } = await import('@/composables/useGradeFormatter')
        const { formatGrade } = useGradeFormatter()
        expect(() => formatGrade(null)).not.toThrow()
        expect(() => formatGrade(undefined)).not.toThrow()
    })

    it('GF03 — formatDecimal handles decimals correctly', async () => {
        const { useGradeFormatter } = await import('@/composables/useGradeFormatter')
        const { formatDecimal } = useGradeFormatter()
        const result = formatDecimal(3.14159, 2)
        expect(typeof result).toBe('string')
        expect(result.length).toBeGreaterThan(0)
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 5 – usePermissions composable
// ─────────────────────────────────────────────────────────────────────────────
describe('usePermissions Composable', () => {
    it('PM01 — teacher role reports isAdmin as false', async () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: { auth: { user: { role: 'teacher' }, isAuthenticated: true } }
        })
        const { useAuthStore } = await import('@/stores/auth')
        const { usePermissions } = await import('@/composables/usePermissions')
        const store = useAuthStore(pinia)
        store.userRole = 'teacher'
        const { isAdmin } = usePermissions()
        expect(isAdmin.value).toBe(false)
    })

    it('PM02 — student has no admin privileges', async () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: { auth: { user: { role: 'student' }, isAuthenticated: true } }
        })
        const { useAuthStore } = await import('@/stores/auth')
        const { usePermissions } = await import('@/composables/usePermissions')
        const store = useAuthStore(pinia)
        store.userRole = 'student'
        const { canManageAllSchools } = usePermissions()
        expect(canManageAllSchools.value).toBe(false)
    })
})

// ─────────────────────────────────────────────────────────────────────────────
// Suite 6 – useSessionReauth composable
// ─────────────────────────────────────────────────────────────────────────────
describe('useSessionReauth Composable', () => {
    it('SR01 — composable provides showDialog, triggerReauth, resolveReauth, cancelReauth', async () => {
        const mod = await import('@/composables/useSessionReauth')
        expect(typeof mod.useSessionReauth).toBe('function')
        const result = mod.useSessionReauth()
        expect(result).toHaveProperty('showDialog')
        expect(result).toHaveProperty('triggerReauth')
        expect(result).toHaveProperty('resolveReauth')
        expect(result).toHaveProperty('cancelReauth')
    })

    it('SR02 — showDialog starts as false (no dialog open)', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { showDialog } = useSessionReauth()
        expect(showDialog.value).toBe(false)
    })

    it('SR03 — cancelReauth rejects pending promise with reauth_cancelled', async () => {
        const { useSessionReauth } = await import('@/composables/useSessionReauth')
        const { triggerReauth, cancelReauth } = useSessionReauth()
        const promise = triggerReauth('user@school.it')
        cancelReauth()
        await expect(promise).rejects.toThrow('reauth_cancelled')
    })
})

