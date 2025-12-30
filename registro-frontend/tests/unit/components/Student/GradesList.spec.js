import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import GradesList from '@/components/Student/GradesList.vue'
import { createTestingPinia } from '@pinia/testing'

describe('GradesList.vue', () => {
    let wrapper

    beforeEach(() => {
        wrapper = mount(GradesList, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            grades: { list: [] }
                        }
                    })
                ]
            }
        })
    })

    it('renders', () => {
        expect(wrapper.exists()).toBe(true)
    })
})
