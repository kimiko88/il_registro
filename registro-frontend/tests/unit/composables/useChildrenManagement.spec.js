import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useChildrenStore } from '@/stores/children'
import { useChildrenManagement } from '@/composables/useChildrenManagement'

const { mockNotify, mockDialog } = vi.hoisted(() => ({
    mockNotify: vi.fn(),
    mockDialog: vi.fn()
}))

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: mockNotify,
            dialog: mockDialog
        })
    }
})

describe('useChildrenManagement Composable', () => {
    let pinia

    beforeEach(() => {
        pinia = createPinia()
        setActivePinia(pinia)
        vi.clearAllMocks()
    })

    it('exposes children computed property bound to childrenStore', () => {
        const store = useChildrenStore()
        store.children = [
            { id: 'c-1', firstName: 'Matteo', lastName: 'Bianchi' },
            { id: 'c-2', firstName: 'Giulia', lastName: 'Bianchi' }
        ]

        const { children, addChild, removeChild } = useChildrenManagement()

        expect(children.value).toHaveLength(2)
        expect(children.value[0].firstName).toBe('Matteo')
        expect(typeof addChild).toBe('function')
        expect(typeof removeChild).toBe('function')
    })

    it('addChild calls notification on success', async () => {
        const { addChild } = useChildrenManagement()
        await addChild({ firstName: 'Nuovo', lastName: 'Studente' })

        expect(mockNotify).toHaveBeenCalledWith(
            expect.objectContaining({
                type: 'positive',
                message: expect.any(String)
            })
        )
    })

    it('removeChild opens confirmation dialog and notifies on confirm', async () => {
        let onOkCallback
        mockDialog.mockReturnValue({
            onOk: vi.fn().mockImplementation((cb) => {
                onOkCallback = cb
                return { onCancel: vi.fn(), onDismiss: vi.fn() }
            })
        })

        const { removeChild } = useChildrenManagement()
        await removeChild('c-1')

        expect(mockDialog).toHaveBeenCalledWith(
            expect.objectContaining({
                cancel: true,
                persistent: true
            })
        )

        // Trigger onOk callback
        expect(onOkCallback).toBeDefined()
        onOkCallback()

        expect(mockNotify).toHaveBeenCalledWith(
            expect.objectContaining({
                type: 'positive'
            })
        )
    })
})

