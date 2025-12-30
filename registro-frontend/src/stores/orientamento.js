import { defineStore } from 'pinia';

export const useOrientamentoStore = defineStore('orientamento', {
    state: () => ({
        events: [],
        loading: false
    }),

    actions: {
        async fetchEvents() {
            this.loading = true;
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                this.events = [
                    { id: 1, title: 'University Open Day', date: '2025-02-15', location: 'Politecnico', hours: 4 },
                    { id: 2, title: 'Career Fair', date: '2025-03-01', location: 'Exhibition Center', hours: 3 }
                ];
            } finally {
                this.loading = false;
            }
        }
    }
});
