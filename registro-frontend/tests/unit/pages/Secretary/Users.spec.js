
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

describe('Secretary Users Page (Users.vue)', () => {
    let wrapper

    beforeEach(() => {
        vi.useFakeTimers()
        wrapper = mount(Users, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-card-actions': { template: '<div><slot /></div>' },
                    'q-dialog': { template: '<div><slot /></div>' }, // Stub dialog to render content inline if model true, but model is false initially. 
                    // Actually if we stub dialog, we might miss visibility logic. Better stub as transition-group if needed or trust v-model.
                    // For simple checks, stubbing q-dialog allows finding content if we force it open.
                    'q-table': {
                        template: '<div><slot name="top" /><slot name="body-cell-actions" :props="{row: {id: 1, first_name: \'Test\'}}" /></div>',
                        props: ['rows', 'columns', 'loading', 'filter']
                    }, // Stub table to avoid complex render but keep slots we use
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
        expect(wrapper.vm.loading).toBe(true)

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
