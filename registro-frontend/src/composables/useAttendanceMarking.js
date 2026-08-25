import { ref } from 'vue';
import { useAttendanceStore } from 'src/stores/attendance';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useAttendanceMarking() {
    const attendanceStore = useAttendanceStore();
    const $q = useQuasar();

    const isSaving = ref(false);

    const markAll = (status) => {
        attendanceStore.records.forEach(r => {
            r.status = status;
            if (status === 'Present') {
                r.time = '';
                r.notes = '';
            }
        });
    };

    const saveAttendance = async (classId, date) => {
        isSaving.value = true;
        const t = i18n?.global?.t;
        try {
            await attendanceStore.submitAttendance(classId, date, attendanceStore.records);
            $q.notify({
                type: 'positive',
                message: t ? t('composables.attendance.saveSuccess') : 'Presenze salvate con successo'
            });
        } catch (err) {
            $q.notify({
                type: 'negative',
                message: t ? t('composables.attendance.saveError') : 'Impossibile salvare le presenze'
            });
        } finally {
            isSaving.value = false;
        }
    };

    return {
        markAll,
        saveAttendance,
        isSaving
    };
}
