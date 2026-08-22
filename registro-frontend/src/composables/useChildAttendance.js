import { computed, watch } from 'vue';
import { useChildrenStore } from 'src/stores/children';
import { useAttendanceStore } from 'src/stores/attendance';
import { useMyAttendance } from 'src/composables/useMyAttendance';

export function useChildAttendance() {
    const childrenStore = useChildrenStore();
    const attendanceStore = useAttendanceStore();

    const { stats, records, loading } = useMyAttendance();

    const fetchAttendanceForChild = async (_studentId) => {
        if (typeof attendanceStore.fetchMyAttendance === 'function') {
            await attendanceStore.fetchMyAttendance();
        }
    };

    const requestChildJustification = async (date, reason) => {
        const targetId = childrenStore.selectedChildId;
        return attendanceStore.requestJustification(date, reason, targetId);
    };

    watch(() => childrenStore.selectedChildId, (newId) => {
        if (newId) fetchAttendanceForChild(newId);
    }, { immediate: true });

    return {
        selectedChild: computed(() => childrenStore.selectedChild),
        stats,
        records,
        loading,
        fetchAttendanceForChild,
        requestChildJustification
    };
}
