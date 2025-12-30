import { ref } from 'vue';
import { useAttendanceStore } from 'src/stores/attendance';
import { useQuasar } from 'quasar';

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
        try {
            await attendanceStore.submitAttendance(classId, date, attendanceStore.records);
            $q.notify({
                type: 'positive',
                message: 'Attendance saved successfully'
            });
        } catch (err) {
            $q.notify({
                type: 'negative',
                message: 'Failed to save attendance'
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
