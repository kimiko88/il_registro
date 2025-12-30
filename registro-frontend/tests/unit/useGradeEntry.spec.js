import { setActivePinia, createPinia } from 'pinia';
import { useGradeEntry } from 'src/composables/useGradeEntry';
import { describe, it, expect, vi, beforeEach } from 'vitest';

// Mock Quasar
const mockNotify = vi.fn();
vi.mock('quasar', () => ({
    useQuasar: () => ({
        notify: mockNotify
    })
}));

describe('useGradeEntry', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        mockNotify.mockClear();
    });

    it('fails validation on empty fields', async () => {
        const { submitGrade } = useGradeEntry();
        const result = await submitGrade({}); // Empty
        expect(result).toBe(false);
        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'warning' }));
    });

    it('submits successfully with valid data', async () => {
        const { submitGrade } = useGradeEntry();

        const validData = {
            studentId: 's1',
            value: 8,
            type: 'Written',
            date: '2025-01-01'
        };

        const result = await submitGrade(validData);
        expect(result).toBe(true);
        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }));
    });
});
