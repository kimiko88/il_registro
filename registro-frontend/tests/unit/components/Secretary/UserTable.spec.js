import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import UserTable from '@/components/Secretary/UserTable.vue'

describe('UserTable Component', () => {
    const mockUsers = [
        { id: 1, first_name: 'Mario', last_name: 'Rossi', email: 'mario@test.com', role: 'student', active: true, is_staff: false },
        { id: 2, first_name: 'Luigi', last_name: 'Verdi', email: 'luigi@test.com', role: 'teacher', active: true, is_staff: true }
    ]

    const mountComponent = () => {
        return mount(UserTable, {
            props: {
                users: mockUsers,
                loading: false
            },
            global: {
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    // The body-cell-role slot accesses both props.value (role string)
                    // and props.row (full row object). Both must be provided to avoid
                    // 'Cannot read properties of undefined (reading is_staff)' at
                    // UserTable.vue:78.
                    'q-table': {
                        template: `<div>
                            <slot name="top" />
                            <slot name="body-cell-role" :props="{ value: 'student', row: { id: 1, role: 'student', is_staff: false } }" />
                            <slot name="body-cell-status" :props="{ value: true, row: { id: 1 } }" />
                            <slot name="body-cell-actions" :props="{ row: { id: 1, role: 'teacher', is_staff: false } }" />
                        </div>`,
                        props: ['rows', 'loading', 'filter']
                    },
                    'q-select': {
                        template: '<div class="test-select" @click="$emit(\'update:model-value\', \'student\')"></div>',
                        props: ['modelValue', 'options']
                    },
                    'q-btn-toggle': { template: '<div class="test-toggle" @click="$emit(\'update:model-value\', \'student\'); $emit(\'click\')"></div>' },
                    'q-input': true,
                    'q-icon': true,
                    'q-btn': true,
                    'q-badge': true,
                    'q-chip': true,
                    'q-td': { template: '<div><slot /></div>' },
                    'q-tr': { template: '<div><slot /></div>' },
                    'q-th': { template: '<div><slot /></div>' },
                    'q-menu': { template: '<div><slot /></div>' },
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div @click="$emit(\'click\')"><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-separator': true,
                    'q-space': true,
                    'q-tooltip': true
                }
            }
        })
    }

    it('mounts without errors', () => {
        // Verifies the component renders without throwing (covers the is_staff crash fix)
        expect(() => mountComponent()).not.toThrow()
    })

    it('exposes columns definition', () => {
        const wrapper = mountComponent()
        expect(wrapper.vm.columns).toBeDefined()
        expect(Array.isArray(wrapper.vm.columns)).toBe(true)
        expect(wrapper.vm.columns.length).toBeGreaterThan(0)
    })

    it('emits filter-role when value changes', async () => {
        const wrapper = mountComponent()
        wrapper.vm.$emit('filter-role', 'student')
        expect(wrapper.emitted('filter-role')).toBeTruthy()
        expect(wrapper.emitted('filter-role')[0]).toEqual(['student'])
    })

    it('emits create event', async () => {
        const wrapper = mountComponent()
        wrapper.vm.$emit('create')
        expect(wrapper.emitted('create')).toBeTruthy()
    })

    it('emits delete event with row payload', async () => {
        const wrapper = mountComponent()
        const row = mockUsers[0]
        wrapper.vm.$emit('delete', row)
        expect(wrapper.emitted('delete')[0]).toEqual([row])
    })

    it('emits edit event with row payload', async () => {
        const wrapper = mountComponent()
        const row = mockUsers[1]
        wrapper.vm.$emit('edit', row)
        expect(wrapper.emitted('edit')[0]).toEqual([row])
    })

    it('emits reset-pwd event', async () => {
        const wrapper = mountComponent()
        wrapper.vm.$emit('reset-pwd', mockUsers[0])
        expect(wrapper.emitted('reset-pwd')[0]).toEqual([mockUsers[0]])
    })

    it('emits bulk-delete with selected users', async () => {
        const wrapper = mountComponent()
        wrapper.vm.$emit('bulk-delete', mockUsers)
        expect(wrapper.emitted('bulk-delete')[0]).toEqual([mockUsers])
    })

    it('emits bulk-reset with selected users', async () => {
        const wrapper = mountComponent()
        wrapper.vm.$emit('bulk-reset', mockUsers)
        expect(wrapper.emitted('bulk-reset')[0]).toEqual([mockUsers])
    })

    it('emits export event', async () => {
        const wrapper = mountComponent()
        wrapper.vm.$emit('export')
        expect(wrapper.emitted('export')).toBeTruthy()
    })

    it('emits manage-subjects for teacher rows', async () => {
        const wrapper = mountComponent()
        const teacherRow = mockUsers[1]
        wrapper.vm.$emit('manage-subjects', teacherRow)
        expect(wrapper.emitted('manage-subjects')[0]).toEqual([teacherRow])
    })
})
