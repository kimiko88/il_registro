import { setActivePinia, createPinia } from 'pinia';
import { useStudentStore } from 'src/stores/student';
import { useGradesStore } from 'src/stores/grades';
import { useMyGrades } from 'src/composables/useMyGrades';
import { useAttendanceStore } from 'src/stores/attendance';
import { useMyAttendance } from 'src/composables/useMyAttendance';
import { describe, it, expect, beforeEach, vi } from 'vitest';

// Mock Quasar
vi.mock('quasar', () => ({
    useQuasar: () => ({ notify: vi.fn() })
}));

describe('Student Logic', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('calculates averages per subject in useMyGrades', () => {
        const gradesStore = useGradesStore();
        gradesStore.grades = [
            { subject: 'Math', value: 8 },
            { subject: 'Math', value: 6 },
            { subject: 'History', value: 9 }
        ];

        const { averages } = useMyGrades();
        expect(averages.value['Math']).toBe('7.0');
        expect(averages.value['History']).toBe('9.0');
    });

    it('calculates attendance stats in useMyAttendance', () => {
        const attStore = useAttendanceStore();
        attStore.records = [
            { status: 'Present' },
            { status: 'Absent' },
            { status: 'Late' }, // Late counts as present in calculation often, or separate. My logic: Present + Late
            { status: 'Present' }
        ];

        const { stats } = useMyAttendance();
        // Present(2) + Late(1) = 3 / 4 = 75%
        expect(stats.value.percentage).toBe(75);
        expect(stats.value.absent).toBe(1);
    });
});
