import { ref } from 'vue'

const politeMessage = ref('')
const assertiveMessage = ref('')

export function useA11yAnnouncer() {
    function announce(message, politeness = 'polite') {
        if (!message) return
        if (politeness === 'assertive') {
            assertiveMessage.value = ''
            setTimeout(() => {
                assertiveMessage.value = message
            }, 50)
        } else {
            politeMessage.value = ''
            setTimeout(() => {
                politeMessage.value = message
            }, 50)
        }
    }

    return {
        announce,
        politeMessage,
        assertiveMessage
    }
}
