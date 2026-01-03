import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import MessageThread from '@/components/Teacher/MessageThread.vue'

const mockSendMessage = vi.fn()
const mockSending = ref(false)

import { ref } from 'vue'

vi.mock('src/composables/useMessages', () => ({
    useMessages: () => ({
        sendMessage: mockSendMessage,
        sending: mockSending
    })
}))

describe('Teacher/MessageThread.vue', () => {
    it('renders messages and sends new one', async () => {
        const messages = [
            { id: 1, from: 'Me', text: 'Hi', date: 'now' },
            { id: 2, from: 'Parent', text: 'Hello', date: 'later' }
        ]

        mockSendMessage.mockResolvedValue(true)

        const wrapper = mount(MessageThread, {
            props: { threadId: 100, messages },
            global: {
                stubs: {
                    'q-scroll-area': { template: '<div><slot /></div>' },
                    'q-chat-message': true,
                    'q-input': {
                        template: '<input @keyup.enter="$emit(\'keyup.enter\')" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
                        props: ['modelValue']
                    },
                    'q-btn': true
                }
            }
        })

        // Check messages rendering
        const chats = wrapper.findAllComponents({ name: 'q-chat-message' })
        expect(chats).toHaveLength(2)

        // Type and send
        const input = wrapper.find('input')
        await input.setValue('Reply text')
        await input.trigger('keyup.enter')

        // Wait for async send and reactivity
        await wrapper.vm.$nextTick()
        await new Promise(resolve => setTimeout(resolve, 0))

        expect(mockSendMessage).toHaveBeenCalledWith(100, 'Reply text')
        // Input should be cleared if success
        expect(input.element.value).toBe('')
    })
})
