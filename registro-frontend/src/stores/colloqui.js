import { defineStore } from 'pinia';

export const useColloquiStore = defineStore('colloqui', {
    state: () => ({
        slots: [],
        bookings: [],
        loading: false,
        error: null
    }),

    actions: {
        async fetchSlots(_rangeStart, _rangeEnd) {
            this.loading = true;
            this.error = null;
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
            } catch (err) {
                this.error = err.message || 'Failed to fetch slots';
                console.error('Error fetching colloquio slots:', err);
            } finally {
                this.loading = false;
            }
        },

        async createSlots(slotsData) {
            this.loading = true;
            this.error = null;
            try {
                await new Promise(resolve => setTimeout(resolve, 500));
                this.slots.push(...slotsData.map(s => ({ ...s, id: Math.random(), booked: false })));
            } catch (err) {
                this.error = err.message || 'Failed to create slots';
                console.error('Error creating slots:', err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async deleteSlot(id) {
            this.loading = true;
            this.error = null;
            try {
                await new Promise(resolve => setTimeout(resolve, 300));
                this.slots = this.slots.filter(s => s.id !== id);
            } catch (err) {
                this.error = err.message || 'Failed to delete slot';
                console.error('Error deleting slot:', err);
                throw err;
            } finally {
                this.loading = false;
            }
        }
    }
});
