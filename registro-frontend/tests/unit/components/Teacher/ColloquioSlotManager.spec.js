import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ColloquioSlotManager from '@/components/Teacher/ColloquioSlotManager.vue'
import { ref } from 'vue'

// Mock composable
const mockGenerateSlots = vi.fn()
const mockGenerating = ref(false)

vi.mock('src/composables/useColloquiScheduling', () => ({
    useColloquiScheduling: () => ({
        generateSlots: mockGenerateSlots,
        generating: mockGenerating
    })
}))

describe('Teacher/ColloquioSlotManager.vue', () => {
    it('submits form to generate slots with defaults', async () => {
        const wrapper = mount(ColloquioSlotManager)

        // Directly call method
        await wrapper.vm.onSubmit()

        expect(mockGenerateSlots).toHaveBeenCalledWith(expect.objectContaining({
            startTime: '15:00',
            endTime: '17:00'
        }))
    })

    it('updates data and submits correct numbers', async () => {
        const wrapper = mount(ColloquioSlotManager)

        // Modify reactive form data directly
        wrapper.vm.form.duration = 20
        wrapper.vm.form.break = 5

        await wrapper.vm.onSubmit()

        expect(mockGenerateSlots).toHaveBeenCalledWith(expect.objectContaining({
            duration: 20,
            break: 5
        }))
    })
})
