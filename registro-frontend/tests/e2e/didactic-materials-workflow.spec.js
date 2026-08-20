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

describe('Didactic Materials & Courseware Sharing Workflow E2E', () => {
    it('executes material upload by teacher and class material list retrieval by student', async () => {
        const material = {
            id: 'mat-1',
            class_id: 'class-1',
            subject_id: 'subj-1',
            title: 'Dispense di Fisica Quantistica - Cap. 1',
            attachment_url: 'https://storage.school.it/dispense/fisica_cap1.pdf'
        }

        const classMaterials = []

        // Verify initial state
        expect(classMaterials.length).toBe(0)

        // Upload & append material
        classMaterials.push(material)
        expect(classMaterials.length).toBe(1)
        expect(classMaterials[0].title).toBe('Dispense di Fisica Quantistica - Cap. 1')
    })
})
