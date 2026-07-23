import { setActivePinia, createPinia } from 'pinia';
import { useAttendanceMarking } from 'src/composables/useAttendanceMarking';
import { useAttendanceStore } from 'src/stores/attendance';
import { describe, it, expect, vi, beforeEach } from 'vitest';

// Mock Quasar
const mockNotify = vi.fn();
vi.mock('quasar', () => ({
    useQuasar: () => ({
        notify: mockNotify
    })
}));

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn(() => Promise.resolve({ data: [] })),
        post: vi.fn(() => Promise.resolve({ data: {} })),
        patch: vi.fn(() => Promise.resolve({ data: {} }))
    }
}))

vi.mock('src/services/attendanceService', () => ({
    attendanceService: {
        getByClass: vi.fn(() => Promise.resolve({ data: { records: [] } })),
        getMyAttendance: vi.fn(() => Promise.resolve({ data: [] }))
    }
}))

describe('useAttendanceMarking', () => {
    let store;

    beforeEach(() => {
        setActivePinia(createPinia());
        store = useAttendanceStore();
        store.records = [
            { status: 'Absent', notes: 'Sick', time: '' },
            { status: 'Late', notes: '', time: '09:00' }
        ];
        mockNotify.mockClear();
    });

    it('marks all students correctly', () => {
        const { markAll } = useAttendanceMarking();
        markAll('Present');

        expect(store.records[0].status).toBe('Present');
        expect(store.records[0].notes).toBe(''); // Should clear details
        expect(store.records[1].status).toBe('Present');
    });

    it('saves attendance', async () => {
        const { saveAttendance } = useAttendanceMarking();
        await saveAttendance('1A', '2025-01-01');
        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }));
    });
});
