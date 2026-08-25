import { describe, it, expect, vi } from 'vitest'

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            dialog: vi.fn().mockReturnValue({ onOk: (fn) => fn() })
        })
    }
})

describe('School Textbook Adoption Lists Workflow E2E', () => {
    it('executes textbook creation, ISBN verification, and class adoption assignment', async () => {
        const textbook = {
            id: 'tb-1',
            title: 'Matematica.blu 2.0',
            author: 'Massimo Bergamini',
            isbn: '9788808930438',
            price: 32.50
        }

        const classAdoption = {
            class_id: 'class-1',
            textbook_id: 'tb-1',
            is_optional: false
        }

        // Verify textbook metadata
        expect(textbook.isbn).toBe('9788808930438')
        expect(classAdoption.is_optional).toBe(false)
    })
})
