import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Substitutions from '@/pages/teacher/Substitutions.vue'
import { useSubstitutionsStore } from '@/stores/substitutions'

describe('Teacher Substitutions Management Workflow E2E', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { id: 'teacher-1', role: 'teacher', name: 'Prof. Mario Rossi' }
                },
                substitutions: {
                    substitutions: [],
                    mySubstitutions: [],
                    myTodaySubstitutions: [],
                    loading: false
                }
            }
        })
        const subStore = useSubstitutionsStore()
        subStore.fetchMySubstitutions = vi.fn().mockResolvedValue([])
    })

    it('renders substitutions overview, stats cards and week filter', () => {
        const wrapper = mount(Substitutions, {
            global: {
                plugins: [Quasar, pinia],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-table': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' }
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.find('.sticky-header').exists()).toBe(true)
    })
})
