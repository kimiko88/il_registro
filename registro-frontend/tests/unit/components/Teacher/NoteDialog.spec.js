import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar, Notify } from 'quasar'
import NoteDialog from '@/components/Teacher/NoteDialog.vue'

// Mock services using vi.hoisted to avoid reference errors
const { mockCreateNote } = vi.hoisted(() => ({
    mockCreateNote: vi.fn()
}))

vi.mock('@/services/notesService', () => ({
    default: {
        createNote: mockCreateNote
    }
}))

// Mock Quasar
const mockNotify = vi.hoisted(() => vi.fn())
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: mockNotify
        }),
        Notify: {
            create: mockNotify
        }
    }
})

describe('NoteDialog.vue', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        wrapper = mount(NoteDialog, {
            global: {
                plugins: [
                    [Quasar, {}]
                ],
                stubs: {
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-select': {
                        template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><option value="generic">Generic</option></select>',
                        props: ['modelValue']
                    },
                    'q-input': {
                        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
                        props: ['modelValue']
                    },
                    'q-btn': { template: '<button :type="type" @click="$emit(\'click\')"></button>', props: ['type'] }
                }
            },
            props: {
                modelValue: true,
                student: { id: 'S1', name: 'Mario' },
                classId: 'C1'
            }
        })
    })

    it('renders correctly when visible', () => {
        expect(wrapper.text()).toContain('Nuova Nota')
        expect(wrapper.text()).toContain('Mario')
    })

    it('resets form when opened', async () => {
        // Simulate closing and reopening
        await wrapper.setProps({ modelValue: false })
        wrapper.vm.noteData.note = 'Old Note'

        await wrapper.setProps({ modelValue: true })
        expect(wrapper.vm.noteData.note).toBe('')
    })

    it('submits form and calls service', async () => {
        wrapper.vm.noteData.note = 'Behavior check'
        wrapper.vm.noteData.type = 'disciplinary'

        await wrapper.find('form').trigger('submit')

        expect(mockCreateNote).toHaveBeenCalledWith({
            student_id: 'S1',
            class_id: 'C1',
            type: 'disciplinary',
            note: 'Behavior check',
            date: expect.any(String)
        })

        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({
            type: 'positive'
        }))

        expect(wrapper.emitted('saved')).toBeTruthy()
        expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
    })

    it('handles submission error', async () => {
        mockCreateNote.mockRejectedValue(new Error('API Error'))

        await wrapper.find('form').trigger('submit')

        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({
            type: 'negative'
        }))
    })
})
