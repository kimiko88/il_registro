import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import ColloquioBookingList from '@/components/Teacher/ColloquioBookingList.vue'
import { useColloquiStore } from '@/stores/colloqui'

describe('Teacher/ColloquioBookingList.vue', () => {
    it('renders slots correctly', () => {
        const slots = [
            { id: 1, startTime: '15:00', endTime: '15:15', booked: false },
            { id: 2, startTime: '15:15', endTime: '15:30', booked: true }
        ]

        const wrapper = mount(ColloquioBookingList, {
            props: { slots },
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        colloqui: {
                            bookings: [
                                { slotId: 2, studentName: 'Mario', parentName: 'Papà' }
                            ]
                        }
                    }
                })],
                stubs: {
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-btn': true
                }
            }
        })

        // Check text content
        expect(wrapper.text()).toContain('15:00 - 15:15')
        // Check booking details for slot 2
        // Since we stubbed q-item-label simply yielding slot, verify computed text
        expect(wrapper.text()).toContain('Mario')
    })

    it('emits delete event for available slots', async () => {
        const slots = [
            { id: 1, startTime: '15:00', endTime: '15:15', booked: false }
        ]
        const wrapper = mount(ColloquioBookingList, {
            props: { slots },
            global: {
                plugins: [createTestingPinia()],
                stubs: {
                    'q-list': { template: '<div><slot /></div>' },
                    'q-item': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<div><slot /></div>' },
                    'q-btn': { template: '<button @click="$emit(\'click\')"></button>', props: ['disable'] }
                }
            }
        })

        const btn = wrapper.find('button')
        await btn.trigger('click')
        expect(wrapper.emitted('delete')).toBeTruthy()
        expect(wrapper.emitted('delete')[0]).toEqual([1])
    })
})
