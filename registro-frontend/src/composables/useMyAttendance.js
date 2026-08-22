import { computed } from 'vue';
import { useAttendanceStore } from 'src/stores/attendance';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useMyAttendance() {
    const store = useAttendanceStore();
    const $q = useQuasar();

    const stats = computed(() => {
        const total = store.records.length;
        if (total === 0) return { present: 0, absent: 0, late: 0, percentage: 100 };

        const present = store.records.filter(r => {
            const s = (r.status || '').toLowerCase();
            return s === 'present' || s === 'late';
        }).length;
        const absent = store.records.filter(r => (r.status || '').toLowerCase() === 'absent').length;
        const late = store.records.filter(r => (r.status || '').toLowerCase() === 'late').length;

        return {
            present,
            absent,
            late,
            percentage: Math.round((present / total) * 100)
        };
    });

    const requestJustification = async (date, reason) => {
        const t = i18n?.global?.t;
        try {
            await store.requestJustification(date, reason);
            $q.notify({ type: 'positive', message: t ? t('composables.attendance.requestSent') : 'Richiesta di giustifica inviata' });
            return true;
        } catch (e) {
            $q.notify({ type: 'negative', message: t ? t('composables.attendance.requestError') : 'Impossibile inviare la richiesta di giustifica' });
            return false;
        }
    };

    return {
        stats,
        requestJustification,
        records: computed(() => store.records),
        loading: computed(() => store.loading)
    };
}
