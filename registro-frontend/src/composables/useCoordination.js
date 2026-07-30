import { computed } from 'vue';
import { useTeacherStore } from 'src/stores/teacher';
import { useClassesStore } from 'src/stores/classes';
import { useAuthStore } from 'src/stores/auth';

export function useCoordination() {
    const teacherStore = useTeacherStore();
    const classesStore = useClassesStore();
    const authStore = useAuthStore();

    const currentUserId = computed(() => authStore.user?.id || teacherStore.profile?.id);

    const coordinatedClasses = computed(() => {
        const uid = currentUserId.value;
        if (!uid) return [];
        return classesStore.classes.filter(c => c.coordinator_id === uid);
    });

    const isCoordinator = computed(() => {
        if (teacherStore.isCoordinator) return true;
        return coordinatedClasses.value.length > 0;
    });

    const getProblemStudents = (_classId) => {
        // Mock logic: Find students with low grade average or high absences
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
