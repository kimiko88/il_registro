import { computed } from 'vue';
import { useTeacherStore } from 'src/stores/teacher';
import { useClassesStore } from 'src/stores/classes';

export function useCoordination() {
    const teacherStore = useTeacherStore();
    const classesStore = useClassesStore();

    const isCoordinator = computed(() => teacherStore.isCoordinator);

    const coordinatedClasses = computed(() =>
        classesStore.classes.filter(c => !!c.coordinator_id)
    );

    const getProblemStudents = (_classId) => {
        // Mock logic: Find students with low grade average or high absences
        // In real app, this would process data from grades/attendance stores or API
        return [
            { id: 's2', name: 'Mario Rossi', issue: 'Low Attendance (65%)' }
        ];
    };

    return {
        isCoordinator,
        coordinatedClasses,
        getProblemStudents
    };
}
