import { ref, computed } from 'vue';
import { useQuasar } from 'quasar';
import { useColloquiStore } from 'src/stores/colloqui';
import colloquiService from '@/services/colloquiService';
import { i18n } from '@/i18n';

export function useColloquiBooking() {
    const store = useColloquiStore();
    const $q = useQuasar();
    const loading = ref(false);

    const fetchAvailableSlots = async (teacherId) => {
        loading.value = true;
        try {
            const res = await colloquiService.getSlots(teacherId ? { teacher_id: teacherId } : {});
            if (res.data) {
                store.slots = Array.isArray(res.data) ? res.data : res.data.slots || [];
            }
        } catch (err) {
            console.error('[useColloquiBooking] fetchAvailableSlots error:', err);
        } finally {
            loading.value = false;
        }
    };

    const bookSlot = async (slotId, notes = '') => {
        loading.value = true;
        const t = i18n?.global?.t;
        try {
            await colloquiService.bookSlot({ slot_id: slotId, notes });
            $q.notify({ type: 'positive', message: t ? t('composables.colloqui.bookingConfirmed') : 'Prenotazione colloquio confermata' });
            return true;
        } catch (err) {
            const errorMsg = err?.response?.data?.error || err?.message || (t ? t('common.error') : 'Errore durante la prenotazione');
            $q.notify({ type: 'negative', message: errorMsg });
            return false;
        } finally {
            loading.value = false;
        }
    };

    const cancelBooking = async (bookingId) => {
        const t = i18n?.global?.t;
        try {
            await colloquiService.cancelSlot(bookingId);
            store.bookings = store.bookings.filter(b => b.id !== bookingId && b.slotId !== bookingId);
            $q.notify({ type: 'positive', message: t ? t('composables.colloqui.bookingCancelled') : 'Prenotazione colloquio annullata' });
            return true;
        } catch (err) {
            const errorMsg = err?.response?.data?.error || err?.message || (t ? t('common.error') : 'Errore durante la cancellazione');
            $q.notify({ type: 'negative', message: errorMsg });
            return false;
        }
    };

    return {
        bookings: computed(() => store.bookings),
        slots: computed(() => store.slots),
        fetchAvailableSlots,
        bookSlot,
        cancelBooking,
        loading
    };
}
