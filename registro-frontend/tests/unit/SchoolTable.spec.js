import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import SchoolTable from 'src/components/Admin/SchoolTable.vue';
import { describe, it, expect, vi } from 'vitest';
import { Quasar } from 'quasar';

describe('SchoolTable.vue', () => {
    it('emits create event when new school button is clicked', async () => {
        const wrapper = mount(SchoolTable, {
            global: {
                plugins: [createTestingPinia({ createSpy: vi.fn })],
                stubs: {
                    'q-table': {
                        template: '<div><slot name="top-right"></slot><slot name="body-cell-actions" :row="{id: 1, name: \'Test\'}"></slot></div>',
                        props: ['rows', 'columns', 'loading', 'pagination']
                    },
                    'q-btn': {
                        template: '<button :class="icon === \'add\' ? \'btn-add\' : \'\'" @click="$emit(\'click\')"></button>',
                        props: ['icon']
                    },
                    'q-input': { template: '<input class="input-filter" @input="$emit(\'update:modelValue\', $event.target.value)" />', props: ['modelValue'] },
                    'q-td': { template: '<div><slot /></div>' },
                    'q-space': true,
                    'q-icon': true
                }
            }
        })

        await wrapper.find('.btn-add').trigger('click')
        expect(wrapper.emitted('create')).toBeTruthy()
    })

    it('emits edit and delete events with row data', async () => {
        const wrapper = mount(SchoolTable, {
            global: {
                plugins: [createTestingPinia({ createSpy: vi.fn })],
                stubs: {
                    'q-table': {
                        template: '<div><slot name="body-cell-actions" :row="{id: 1, name: \'Test\'}"></slot></div>'
                    },
                    'q-btn': {
                        template: '<button :class="\'btn-\' + icon" @click="$emit(\'click\')"></button>',
                        props: ['icon']
                    },
                    'q-td': { template: '<div><slot /></div>' }
                }
            }
        })

        await wrapper.find('.btn-edit').trigger('click')
        expect(wrapper.emitted('edit')).toBeTruthy()
        expect(wrapper.emitted('edit')[0][0]).toEqual({ id: 1, name: 'Test' })

        await wrapper.find('.btn-delete').trigger('click')
        expect(wrapper.emitted('delete')).toBeTruthy()
        expect(wrapper.emitted('delete')[0][0]).toBe(1)
    })

    it('handles filter updates', async () => {
        const wrapper = mount(SchoolTable, {
            global: {
                plugins: [createTestingPinia({ createSpy: vi.fn })],
                stubs: {
                    'q-table': { template: '<div><slot name="top-right"></slot></div>' },
                    'q-input': { template: '<input class="input-filter" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />', props: ['modelValue'] },
                    'q-btn': true, 'q-space': true, 'q-icon': true
                }
            }
        })

        const input = wrapper.find('.input-filter')
        await input.setValue('High School')

        expect(wrapper.vm.filter).toBe('High School')
    })

    it('emits request event for pagination', async () => {
        const wrapper = mount(SchoolTable, {
            global: {
                plugins: [createTestingPinia({ createSpy: vi.fn })],
                stubs: {
                    'q-table': {
                        template: '<div class="q-table-stub" @click="$emit(\'request\', { page: 2 })"></div>',
                        emits: ['request']
                    }
                }
            }
        })

        await wrapper.find('.q-table-stub').trigger('click')
        expect(wrapper.emitted('request')).toBeTruthy()
        expect(wrapper.emitted('request')[0][0]).toEqual({ page: 2 })
    })
})
