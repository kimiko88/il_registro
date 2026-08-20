import { ref, computed } from 'vue';
import { useQuasar } from 'quasar';
import { useColloquiStore } from 'src/stores/colloqui';
import { i18n } from '@/i18n';

export function useColloquiBooking() {
    const store = useColloquiStore();
    const $q = useQuasar();
    const loading = ref(false);

    const fetchAvailableSlots = async (_teacherId) => {
        // Mock fetch
    };

    const bookSlot = async (slotId) => {
        loading.value = true;
        const t = i18n?.global?.t;
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
            $q.notify({ type: 'positive', message: t ? t('composables.colloqui.bookingConfirmed') : 'Prenotazione colloquio confermata' });
            return true;
        } finally {
            loading.value = false;
        }
    };

    const cancelBooking = async (bookingId) => {
        const t = i18n?.global?.t;
        store.bookings = store.bookings.filter(b => b.id !== bookingId);
        $q.notify({ type: 'positive', message: t ? t('composables.colloqui.bookingCancelled') : 'Prenotazione colloquio annullata' });
    };

    return {
        bookings: computed(() => store.bookings),
        fetchAvailableSlots,
        bookSlot,
        cancelBooking,
        loading
    };
}
