import { setActivePinia, createPinia } from 'pinia';
import { useCoordination } from 'src/composables/useCoordination';
import { useTeacherStore } from 'src/stores/teacher';
import { useClassesStore } from 'src/stores/classes';
import { describe, it, expect, beforeEach } from 'vitest';

describe('useCoordination', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('identifies coordinator status', () => {
        const teacherStore = useTeacherStore();
        teacherStore.profile = { isCoordinator: true };

        const { isCoordinator } = useCoordination();
        expect(isCoordinator.value).toBe(true);
    });

    it('filters coordinated classes', () => {
        const classesStore = useClassesStore();
        classesStore.classes = [
            { id: '1A', coordinator: false },
            { id: '2B', coordinator: true }
        ];

        const { coordinatedClasses } = useCoordination();
        expect(coordinatedClasses.value).toHaveLength(1);
        expect(coordinatedClasses.value[0].id).toBe('2B');
    });
});
