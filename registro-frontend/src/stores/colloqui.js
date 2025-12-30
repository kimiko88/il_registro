import { defineStore } from 'pinia';

export const useColloquiStore = defineStore('colloqui', {
    state: () => ({
        slots: [],
        bookings: [],
        loading: false,
        error: null
    }),

    actions: {
        async fetchSlots(rangeStart, rangeEnd) {
            this.loading = true;
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                // Mock slots
                this.slots = [
                    { id: 1, date: '2025-01-20', startTime: '15:00', endTime: '15:15', booked: true },
                    { id: 2, date: '2025-01-20', startTime: '15:15', endTime: '15:30', booked: false }
                ];
                this.bookings = [
                    { slotId: 1, parentName: 'Mrs. Rossi', studentName: 'Mario Rossi', notes: 'Math grade' }
                ];
            } finally {
                this.loading = false;
            }
        },

        async createSlots(slotsData) {
            // Mock API
            await new Promise(resolve => setTimeout(resolve, 500));
            this.slots.push(...slotsData.map(s => ({ ...s, id: Math.random(), booked: false })));
        },

        async deleteSlot(id) {
            this.slots = this.slots.filter(s => s.id !== id);
        }
    }
});
