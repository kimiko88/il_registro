import { setActivePinia, createPinia } from 'pinia';
import { useAttendanceStore } from 'src/stores/attendance';
import { describe, it, expect, beforeEach } from 'vitest';

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
