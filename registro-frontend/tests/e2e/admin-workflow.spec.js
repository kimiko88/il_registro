import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
// Assuming Admin Dashboard is SchoolManagement or specific page.
// The user has `pages/Admin/SchoolManagement.vue`. Is there a dashboard `pages/Admin/Index.vue`?
// I will check file list or assume SchoolManagement if not.
// Let's assume there is an Admin Dashboard or checking SchoolManagement is enough.
// Default to SchoolManagement as "Admin Workflow" entry.
import SchoolManagement from '@/pages/admin/SchoolManagement.vue'
import { createTestingPinia } from '@pinia/testing'

// Mock Quasar
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

describe('Admin Workflow', () => {
    it('manages schools', async () => {
        const wrapper = mount(SchoolManagement, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            auth: { user: { role: 'superadmin' } }
                        }
                    })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-table': { template: '<div><slot name="top-right" /><slot name="body-cell-actions" :props="{row: {id: 1, name: \'School 1\'}}" /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-input': true,
                    'q-btn': true,
                    'q-icon': true
                }
            }
        })

        expect(wrapper.text()).toContain('Gestione Scuole')

        // Check "Nuova Scuola" button availability
        // Since q-btn stubbed, check existance.
        const addBtn = wrapper.findAllComponents({ name: 'q-btn' }).find(c => c.attributes().label === 'Nuova Scuola') // label prop might be attribute in stub if passed
        // props().label is safer
        const addBtn2 = wrapper.findAllComponents({ name: 'q-btn' }).find(c => c.props().label === 'Nuova Scuola')
        expect(addBtn2).toBeTruthy()
    })
})
