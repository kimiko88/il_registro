import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import DocumentList from '@/components/Student/DocumentList.vue'

describe('Student/DocumentList.vue', () => {
    const documents = [
        { id: 1, title: 'Math HW', date: '2023-01-01', type: 'Homework' },
        { id: 2, title: 'History Report', date: '2023-02-01', type: 'Report' }
    ]

    it('renders list of documents', () => {
        const wrapper = mount(DocumentList, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-list': { template: '<div class="list"><slot /></div>' },
                    'q-item': { template: '<div class="item" @click="$emit(\'click\')"><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<span><slot /></span>' },
                    'q-icon': true,
                    'q-btn': true
                }
            },
            props: { documents }
        })

        expect(wrapper.findAll('.item')).toHaveLength(2)
        expect(wrapper.text()).toContain('Math HW')
        expect(wrapper.text()).toContain('History Report')
    })

    it('emits download event on item click', async () => {
        const wrapper = mount(DocumentList, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div class="item" @click="$emit(\'click\')"><slot /></div>' },
                    'q-item-section': true, 'q-item-label': true, 'q-icon': true, 'q-btn': true
                }
            },
            props: { documents }
        })

        await wrapper.findAll('.item')[0].trigger('click')
        expect(wrapper.emitted('download')).toBeTruthy()
        expect(wrapper.emitted('download')[0][0]).toEqual(documents[0])
    })
})
