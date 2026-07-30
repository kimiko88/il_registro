import { setActivePinia, createPinia } from 'pinia';
import { useCoordination } from 'src/composables/useCoordination';
import { useTeacherStore } from 'src/stores/teacher';
import { useClassesStore } from 'src/stores/classes';
import { useAuthStore } from 'src/stores/auth';
import { describe, it, expect, beforeEach } from 'vitest';

describe('useCoordination', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('identifies coordinator status', () => {
        const teacherStore = useTeacherStore();
        teacherStore.profile = { is_coordinator: true };

        const { isCoordinator } = useCoordination();
        expect(isCoordinator.value).toBe(true);
    });

    it('filters coordinated classes', () => {
        const authStore = useAuthStore();
        authStore.user = { id: 'teacher123' };

        const classesStore = useClassesStore();
        classesStore.classes = [
            { id: '1A', coordinator_id: null },
            { id: '2B', coordinator_id: 'teacher123' }
        ];

        const { coordinatedClasses } = useCoordination();
        expect(coordinatedClasses.value).toHaveLength(1);
        expect(coordinatedClasses.value[0].id).toBe('2B');
    });
});
