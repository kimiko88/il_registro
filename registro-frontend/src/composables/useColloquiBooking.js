import { ref } from 'vue';
import { useQuasar } from 'quasar';
import { useColloquiStore } from 'src/stores/colloqui';

export function useColloquiBooking() {
    const store = useColloquiStore();
    const $q = useQuasar();
    const loading = ref(false);

    const fetchAvailableSlots = async (teacherId) => {
        // Mock fetch
    };

    const bookSlot = async (slotId) => {
        loading.value = true;
        try {
            await new Promise(resolve => setTimeout(resolve, 500));
            // Mock booking
            store.bookings.push({
                id: Math.random(),
                slotId,
                teacherName: 'Prof. Bianchi',
                subject: 'Mathematics',
                date: '2025-01-25 15:00'
            });
            $q.notify({ type: 'positive', message: 'Booking Confirmed' });
            return true;
        } finally {
            loading.value = false;
        }
    };

    const cancelBooking = async (bookingId) => {
        store.bookings = store.bookings.filter(b => b.id !== bookingId);
        $q.notify({ type: 'positive', message: 'Booking Cancelled' });
    };

    return {
        bookings: store.bookings, // This should be reactive
        bookSlot,
        cancelBooking,
        loading
    };
}
