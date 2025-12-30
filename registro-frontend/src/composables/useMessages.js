import { ref } from 'vue';
import { useCommunicationsStore } from 'src/stores/communications';

export function useMessages() {
    const store = useCommunicationsStore();
    const sending = ref(false);

    const sendMessage = async (threadId, text) => {
        if (!text.trim()) return;
        sending.value = true;
        try {
            await store.sendMessage(threadId, text);
            return true;
        } finally {
            sending.value = false;
        }
    };

    return { sendMessage, sending };
}
