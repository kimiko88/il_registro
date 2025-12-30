import { computed } from 'vue';
import { useAttendanceStore } from 'src/stores/attendance';
import { useQuasar } from 'quasar';

export function useMyAttendance() {
    const store = useAttendanceStore();
    const $q = useQuasar();

    const stats = computed(() => {
        const total = store.records.length;
        if (total === 0) return { present: 0, absent: 0, late: 0, percentage: 100 };

        const present = store.records.filter(r => r.status === 'Present' || r.status === 'Late').length;
        const absent = store.records.filter(r => r.status === 'Absent').length;
        const late = store.records.filter(r => r.status === 'Late').length;

        return {
            present,
            absent,
            late,
            percentage: Math.round((present / total) * 100)
        };
    });

    const requestJustification = async (date, reason) => {
        try {
            await store.requestJustification(date, reason);
            $q.notify({ type: 'positive', message: 'Request sent' });
            return true;
        } catch (e) {
            $q.notify({ type: 'negative', message: 'Failed to send request' });
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
