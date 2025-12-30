import { defineStore } from 'pinia';

export const useCommunicationsStore = defineStore('communications', {
    state: () => ({
        threads: [],
        loading: false,
        activeThreadId: null
    }),

    getters: {
        activeThread: (state) => state.threads.find(t => t.id === state.activeThreadId),
        unreadCount: (state) => state.threads.reduce((acc, t) => acc + (t.unread ? 1 : 0), 0)
    },

    actions: {
        async fetchThreads() {
            this.loading = true;
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                // Mock Threads
                this.threads = [
                    { id: 1, subject: 'Regarding Mario', recipient: 'Mrs. Rossi', lastMessage: 'Please confirm availability', date: '2025-01-20', unread: true, messages: [] },
                    { id: 2, subject: 'Class Trip', recipient: 'Principal', lastMessage: 'Approved', date: '2025-01-18', unread: false, messages: [] }
                ];
            } finally {
                this.loading = false;
            }
        },

        async fetchMessages(threadId) {
            // Mock fetch details
            const thread = this.threads.find(t => t.id === threadId);
            if (thread && thread.messages.length === 0) {
                thread.messages = [
                    { id: 1, from: 'Me', text: 'Hello, I would like to discuss Mario\'s progress.', date: '2025-01-19 10:00' },
                    { id: 2, from: 'Mrs. Rossi', text: 'Sure, when are you available?', date: '2025-01-19 14:00' }
                ];
            }
            this.activeThreadId = threadId;
            if (thread) thread.unread = false;
        },

        async sendMessage(threadId, text) {
            // Mock send
            await new Promise(resolve => setTimeout(resolve, 300));
            const thread = this.threads.find(t => t.id === threadId);
            if (thread) {
                thread.messages.push({ id: Math.random(), from: 'Me', text, date: new Date().toISOString() });
            }
        }
    }
});
