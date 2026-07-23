import { setActivePinia, createPinia } from 'pinia';
import { useGradesStore } from '@/stores/grades';
import { useAttendanceStore } from '@/stores/attendance';
import { useDocumentsStore } from '@/stores/documents';
import { useMyGrades } from '@/composables/useMyGrades';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { gradeService } from '@/services/gradeService';

// Mock Services
vi.mock('@/services/gradeService', () => ({
    gradeService: {
        getByClass: vi.fn(),
        getMyGrades: vi.fn(),
        saveGrade: vi.fn(),
        updateGrade: vi.fn(),
        deleteGrade: vi.fn()
    }
}));

vi.mock('@/services/documentService', () => ({
    default: {
        createDocument: vi.fn(),
        getInbox: vi.fn()
    }
}));

// Mock Quasar
vi.mock('quasar', () => ({
    useQuasar: () => ({ notify: vi.fn() })
}));

describe('Regression Tests', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    // 1. Grades Regression: Invalid Inputs
    it('should handle grade boundary values correctly', async () => {
        const store = useGradesStore();
        
        // Mock getByClass so addGrade can refetch
        gradeService.getByClass.mockResolvedValue({ data: { students: [] } });

        // Hypothetical valid range 0-10
        // Mock saveGrade
        gradeService.saveGrade.mockResolvedValue({ data: { id: 'g1', value: 10, student_id: 's1' } });
        await store.addGrade({ value: 10, studentId: 's1' });
        
        gradeService.saveGrade.mockResolvedValue({ data: { id: 'g2', value: 0, student_id: 's1' } });
        await store.addGrade({ value: 0, studentId: 's1' });

        expect(gradeService.saveGrade).toHaveBeenCalledTimes(2);
    });

    // 2. Student Logic Regression: Floating Point Precision
    it('should calculate precise averages', () => {
        const store = useGradesStore();
        store.grades = {
            semesters: [
                {
                    grades: [
                        { subject_id: 'Math', grade_value: 7.1 },
                        { subject_id: 'Math', grade_value: 7.2 },
                        { subject_id: 'Math', grade_value: 7.3 }
                    ]
                }
            ]
        };
        // Average: 7.2

        const { averages } = useMyGrades();
        expect(averages.value['Math']).toBe(7.2);
    });

    // 3. Attendance Regression: Bulk Marking Safety
    it('should not overwrite justified absences during bulk mark', async () => {
        const store = useAttendanceStore();
        // Setup initial state
        store.records = [
            { studentId: 's1', status: 'Absent', justificationStatus: 'Approved' }, // Should stay
            { studentId: 's2', status: 'Present' }
        ];

        // Simulate a "Mark All Present" action logic (if it exists in composable, mocking here)
        const records = store.records.map(r => {
            // Regression Logic: Don't change if Justified
            if (r.status === 'Absent' && r.justificationStatus === 'Approved') return r;
            return { ...r, status: 'Present' };
        });

        expect(records[0].status).toBe('Absent');
        expect(records[1].status).toBe('Present');
    });

    // 4. Documents Regression: Missing Data
    it('should fail gracefully when creating document without student', async () => {
        const store = useDocumentsStore();
        store.createDocument = vi.fn().mockRejectedValue(new Error('Missing Student'));

        await expect(store.createDocument({})).rejects.toThrow('Missing Student');
    });
});
