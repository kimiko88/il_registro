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
        const fromStore = classesStore.classes.filter(c => c.coordinator_id === uid);
        const assignedClassIds = (authStore.user?.assignments || [])
            .filter(a => (a.assignment_type === 'coordinatore_classe' || a.assignment_type === 'coordinator') && a.is_active !== false && a.scope_id)
            .map(a => a.scope_id);
        
        if (assignedClassIds.length === 0) return fromStore;

        const combined = [...fromStore];
        for (const cid of assignedClassIds) {
            if (!combined.some(c => c.id === cid)) {
                const found = classesStore.classes.find(c => c.id === cid);
                if (found) combined.push(found);
            }
        }
        return combined;
    });

    const isCoordinator = computed(() => {
        if (teacherStore.isCoordinator) return true;
        if (coordinatedClasses.value.length > 0) return true;
        return (authStore.user?.assignments || []).some(
            a => (a.assignment_type === 'coordinatore_classe' || a.assignment_type === 'coordinator') && a.is_active !== false
        );
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
