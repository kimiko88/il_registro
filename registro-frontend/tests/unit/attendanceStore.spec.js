import { setActivePinia, createPinia } from 'pinia';
import { describe, it, expect, beforeEach, vi } from 'vitest';

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn(() => Promise.resolve({ data: [] })),
        post: vi.fn(() => Promise.resolve({ data: {} })),
        patch: vi.fn(() => Promise.resolve({ data: {} }))
    }
}))

vi.mock('src/services/attendanceService', () => ({
    attendanceService: {
        getByClass: vi.fn(() => Promise.resolve({
            data: {
                records: [
                    { student_id: 's1', student_name: 'Giuseppe Verdi', status: 'Present', notes: '', entry_time: '' },
                    { student_id: 's2', student_name: 'Mario Rossi', status: 'Absent', notes: '', entry_time: '' }
                ]
            }
        }))
    }
}))

import { useAttendanceStore } from 'src/stores/attendance';

describe('Attendance Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('counts present/absent correctly', () => {
        const store = useAttendanceStore();
        store.records = [
            { status: 'Present' },
            { status: 'Present' },
            { status: 'Absent' },
            { status: 'Late' }
        ];

        expect(store.presentCount).toBe(2);
        expect(store.absentCount).toBe(1);
        expect(store.lateCount).toBe(1);
    });

    it('updates records on fetch', async () => {
        const store = useAttendanceStore();
        await store.fetchDailyAttendance('1A', '2025-01-01');
        expect(store.records.length).toBeGreaterThan(0);
        expect(store.loading).toBe(false);
    });
});
