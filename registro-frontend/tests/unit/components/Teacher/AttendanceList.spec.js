import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import AttendanceList from '@/components/Teacher/AttendanceList.vue'
import { useAttendanceStore } from '@/stores/attendance'

describe('AttendanceList.vue', () => {
    let wrapper
    let store

    beforeEach(() => {
        wrapper = mount(AttendanceList, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            attendance: {
                                records: [
                                    { studentId: '1', name: 'Mario Rossi', status: 'Present', time: '', notes: '' },
                                    { studentId: '2', name: 'Luca Bianchi', status: 'Absent', time: '', notes: 'Malattia' }
                                ]
                            }
                        },
                        stubActions: false
                    })
                ],
                stubs: {
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-btn-group': { template: '<div><slot /></div>' },
                    'q-btn': { template: '<button @click="$emit(\'click\')">{{ label }}</button>', props: ['label'] },
                    'q-dialog': { template: '<div><slot /></div>' }, // Render dialog content inline for testing
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-input': {
                        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
                        props: ['modelValue']
                    }
                }
            }
        })
        store = useAttendanceStore()
    })

    it('renders list of students', () => {
        expect(wrapper.text()).toContain('Mario Rossi')
        expect(wrapper.text()).toContain('Luca Bianchi')
    })

    it('updates status to Present', async () => {
        const buttons = wrapper.findAll('button')
        const presentBtn = buttons.find(b => b.text() === 'P')

        await presentBtn.trigger('click')

        expect(store.records[0].status).toBe('Present')
    })

    it('updates status to Absent', async () => {
        const buttons = wrapper.findAll('button')
        const absentBtn = buttons.find(b => b.text() === 'A')

        await absentBtn.trigger('click')

        expect(store.records[0].status).toBe('Absent')
    })

    it('opens note dialog and saves note', async () => {
        // Edit note button has icon but no label in main list, usually.
        // In AttendanceList.vue: <q-btn flat round icon="edit_note" ... />
        // Our stub renders {{ label }}. If label is undefined, text is empty.
        // So look for empty text button? Or stub icon?
        // Let's find button that is NOT P, A, or L.

        const buttons = wrapper.findAll('button')
        const editBtn = buttons.find(b => !['P', 'A', 'L'].includes(b.text()))

        await editBtn.trigger('click')
        expect(wrapper.vm.showNoteDialog).toBe(true)

        // Set note value
        wrapper.vm.currentNote = 'Test Note'
        await wrapper.vm.saveNote()

        expect(store.records[0].notes).toBe('Test Note')
        expect(wrapper.vm.showNoteDialog).toBe(false)
    })
})
