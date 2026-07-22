import { computed, watch } from 'vue';
import { useChildrenStore } from 'src/stores/children';
import { useAttendanceStore } from 'src/stores/attendance';
import { useMyAttendance } from 'src/composables/useMyAttendance'; // Reuse!

export function useChildAttendance() {
    const childrenStore = useChildrenStore();
    const attendanceStore = useAttendanceStore();

    const { stats, records, loading } = useMyAttendance();

    const fetchAttendanceForChild = async (_studentId) => {
        // In real app, pass studentId to fetchMyAttendance or a specific fetchChildAttendance
        await attendanceStore.fetchMyAttendance();
    };

    watch(() => childrenStore.selectedChildId, (newId) => {
        if (newId) fetchAttendanceForChild(newId);
    }, { immediate: true });

    return {
        selectedChild: computed(() => childrenStore.selectedChild),
        stats,
        records,
        loading
    };
}
