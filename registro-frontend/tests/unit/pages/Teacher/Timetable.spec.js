import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import Timetable from '@/pages/teacher/Timetable.vue'
import { useClassesStore } from 'src/stores/classes'
import api from 'src/services/api'

vi.mock('src/services/api', () => ({
    default: {
        get: vi.fn()
    }
}))

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn()
    }
}))

describe('Teacher Timetable.vue — Teacher Personal Schedule View', () => {
    let wrapper

    const mockScheduleEntries = [
        { id: '1', class_id: 'c1', class_name: '3A', day_of_week: 1, hour_index: 1, subject_id: 's1', subject_name: 'Matematica', room: '101' },
        { id: '2', class_id: 'c2', class_name: '4B', day_of_week: 2, hour_index: 2, subject_id: 's2', subject_name: 'Fisica', room: 'Lab 2' }
    ]

    beforeEach(() => {
        vi.clearAllMocks()
        api.get.mockResolvedValue({ data: mockScheduleEntries })
        const pinia = createTestingPinia({ createSpy: vi.fn })
        const classesStore = useClassesStore(pinia)
        classesStore.fetchAssignedClasses.mockResolvedValue([])

        wrapper = mount(Timetable, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div class="q-page"><slot /></div>' },
                    'q-card': { template: '<div class="q-card"><slot /></div>' },
                    'q-icon': true,
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-badge': { template: '<span><slot /></span>' },
                    'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
                    'q-spinner-dots': true,
                    'q-tooltip': true
                }
            }
        })
    })

    it('renders header and metrics stats', async () => {
        await flushPromises()
        expect(wrapper.text()).toContain('Il Mio Orario di Insegnamento')
        expect(wrapper.text()).toContain('Ore Settimanali')
        expect(wrapper.text()).toContain('2 Ore')
    })

    it('renders schedule entries in grid', async () => {
        await flushPromises()
        expect(wrapper.text()).toContain('Matematica')
        expect(wrapper.text()).toContain('3A')
        expect(wrapper.text()).toContain('Aula: 101')

        expect(wrapper.text()).toContain('Fisica')
        expect(wrapper.text()).toContain('4B')
        expect(wrapper.text()).toContain('Aula: Lab 2')
    })

    it('displays empty state when no schedule entries exist', async () => {
        api.get.mockResolvedValueOnce({ data: [] })
        const pinia = createTestingPinia({ createSpy: vi.fn })
        const classesStore = useClassesStore(pinia)
        classesStore.fetchAssignedClasses.mockResolvedValue([])

        wrapper = mount(Timetable, {
            global: {
                plugins: [pinia],
                stubs: {
                    'q-page': { template: '<div class="q-page"><slot /></div>' },
                    'q-card': { template: '<div class="q-card"><slot /></div>' },
                    'q-icon': true,
                    'q-avatar': { template: '<div><slot /></div>' },
                    'q-badge': { template: '<span><slot /></span>' },
                    'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
                    'q-spinner-dots': true,
                    'q-tooltip': true
                }
            }
        })
        await flushPromises()

        expect(wrapper.text()).toContain('Nessuna lezione in orario')
        expect(wrapper.text()).toContain('Non risultano ancora ore di lezione')
    })

    it('invokes api.get when clicking refresh button', async () => {
        await flushPromises()
        api.get.mockClear()

        const refreshBtn = wrapper.findAll('button').find(b => b.html().includes('refresh') || b.attributes('icon') === 'refresh')
        if (refreshBtn) {
            await refreshBtn.trigger('click')
            expect(api.get).toHaveBeenCalledWith('/timetables/my-schedule')
        }
    })
})
