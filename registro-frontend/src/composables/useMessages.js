import { ref } from 'vue';
import { useCommunicationsStore } from 'src/stores/communications';

export function useMessages() {
    const store = useCommunicationsStore();
    const sending = ref(false);

    const sendMessage = async (threadIdOrPayload, text) => {
        let payload;
        if (typeof threadIdOrPayload === 'object' && threadIdOrPayload !== null) {
            payload = threadIdOrPayload;
        } else {
            if (!text || !text.trim()) return false;
            payload = {
                thread_id: threadIdOrPayload,
                body: text
            };
        }
        sending.value = true;
        try {
            await store.sendMessage(payload);
            return true;
        } finally {
            sending.value = false;
        }
    };

    return { sendMessage, sending };
}
