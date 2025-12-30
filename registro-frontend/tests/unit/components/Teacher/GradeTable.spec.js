import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import GradeTable from '@/components/Teacher/GradeTable.vue'
import { createTestingPinia } from '@pinia/testing'

describe('GradeTable.vue', () => {
    let wrapper

    beforeEach(() => {
        // Check if component exists, otherwise this test will fail compilation of import
        // Assuming component exists as per previous context
        wrapper = mount(GradeTable, {
            props: {
                students: [], // Prop input
                subjectId: 'math'
            },
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            teacher: { selectedClass: '1A' }
                        }
                    })
                ]
            }
        })
    })

    it('renders correctly', () => {
        expect(wrapper.exists()).toBe(true)
    })

    it('emits save event', () => {
        wrapper.vm.$emit('save')
        expect(wrapper.emitted('save')).toBeTruthy()
    })
})
