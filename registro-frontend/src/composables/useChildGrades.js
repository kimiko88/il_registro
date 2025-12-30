import { computed, watch } from 'vue';
import { useChildrenStore } from 'src/stores/children';
import { useGradesStore } from 'src/stores/grades';
import { useMyGrades } from 'src/composables/useMyGrades'; // Reuse student logic!

export function useChildGrades() {
    const childrenStore = useChildrenStore();
    const gradesStore = useGradesStore();

    // Reuse the logic from useMyGrades which relies on store.grades
    // We just need to make sure store.grades is populated with the SELECTED child's grades
    const { gradesBySubject, averages, getTrend, loading } = useMyGrades();

    const fetchGradesForChild = async (studentId) => {
        await gradesStore.fetchMyGrades(studentId); // Reusing this action is fine as it fetches by ID
    };

    // Watch for selected child changes to refetch
    watch(() => childrenStore.selectedChildId, (newId) => {
        if (newId) fetchGradesForChild(newId);
    }, { immediate: true });

    return {
        selectedChild: computed(() => childrenStore.selectedChild),
        gradesBySubject,
        averages,
        getTrend,
        loading
    };
}
