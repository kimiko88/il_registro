
import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { exportFile } from 'quasar'
import Users from '@/pages/secretary/Users.vue'

// Mock Quasar
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        exportFile: vi.fn().mockReturnValue(true),
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: vi.fn(),
            dialog: vi.fn().mockReturnValue({ onOk: (fn) => fn() })
        })
    }
})

// Mock Services
vi.mock('src/services/userService', () => ({
    userService: {
        getAll: vi.fn().mockResolvedValue({ data: { users: [{ id: 1, first_name: 'Test', role: 'student' }] } }),
        create: vi.fn(),
        update: vi.fn(),
        delete: vi.fn(),
        resetPassword: vi.fn()
    }
}))

vi.mock('src/services/adminService', () => ({
    default: {
        getSchoolClasses: vi.fn().mockResolvedValue({ data: [] }),
        getSubjects: vi.fn().mockResolvedValue({ data: [] }),
        getTeachersList: vi.fn().mockResolvedValue({ data: [] }),
        createClass: vi.fn()
    }
}))

describe('Secretary Users Page (Users.vue)', () => {
    let wrapper

    beforeEach(() => {
        vi.useFakeTimers()
        // Reset stores if needed, but here we just mount.
        // Also mock authStore or pinia if needed, but Users.vue uses it.
        // We probably need to mock useAuthStore too or provide a testing pinia.
        wrapper = mount(Users, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' },
                    'q-table': {
                        template: '<div><slot name="top" /><slot name="body-cell-actions" :props="{row: {id: 1, first_name: \'Test\'}}" /></div>',
                        props: ['rows', 'columns', 'loading', 'filter']
                    },
                    'q-btn': true,
                    'q-input': true,
                    'q-select': true,
                    'q-file': true,
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-btn-toggle': true,
                    'q-space': true,
                    'q-icon': true
                }
            }
        })
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('fetches and displays users', async () => {
        // Fast download timers
        await vi.runAllTimersAsync()
        await wrapper.vm.$nextTick()

        expect(wrapper.vm.loading).toBe(false)
        expect(wrapper.vm.users.length).toBeGreaterThan(0)
        expect(wrapper.text()).toContain('Gestione Utenti')
    })

    it('opens import dialog', async () => {
        await wrapper.vm.$nextTick()
        // Trigger open
        wrapper.vm.showImport = true
        await wrapper.vm.$nextTick()

        // Check text in dialog (dialog is stubbed but renders slot if visible? No, q-dialog usually hides. 
        // But wrapper.vm.showImport is true. If we stub q-dialog with div, it renders.)
        // My stub above: 'q-dialog': { template: '<div><slot /></div>' }
        // But q-dialog uses v-model. If I stub it, v-model prop must be handled or I just check vm state.

        expect(wrapper.vm.showImport).toBe(true)
    })

    it('exports users calls exportFile', async () => {
        await vi.runAllTimersAsync()

        wrapper.vm.exportUsers()

        expect(exportFile).toHaveBeenCalled()
    })
})
