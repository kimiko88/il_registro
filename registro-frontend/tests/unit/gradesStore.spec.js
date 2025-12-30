import { setActivePinia, createPinia } from 'pinia';
import { useGradesStore } from 'src/stores/grades';
import { describe, it, expect, beforeEach } from 'vitest';

describe('Grades Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('calculates class average correctly', () => {
        const store = useGradesStore();
        store.grades = [
            { value: 6 },
            { value: 8 },
            { value: 10 }
        ];

        expect(store.classAverage).toBe('8.0');
    });

    it('filters grades by student', () => {
        const store = useGradesStore();
        store.grades = [
            { studentId: 's1', value: 6 },
            { studentId: 's2', value: 8 }
        ];

        const s1Grades = store.getGradesByStudent('s1');
        expect(s1Grades).toHaveLength(1);
        expect(s1Grades[0].value).toBe(6);
    });

    it('adds a grade', async () => {
        const store = useGradesStore();
        await store.addGrade({ studentId: 's1', value: 9 });
        expect(store.grades).toHaveLength(1);
        expect(store.grades[0].value).toBe(9);
    });
});
