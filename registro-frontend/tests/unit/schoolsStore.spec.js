import { setActivePinia, createPinia } from 'pinia';
import { useSchoolStore } from 'src/stores/schools'; // Fixed import name
import { describe, it, expect, vi, beforeEach } from 'vitest';
import schoolService from 'src/services/schoolService';

vi.mock('src/services/schoolService');

describe('Schools Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('fetching schools updates state correctly', async () => {
        const store = useSchoolStore();
        const mockData = {
            items: [{ id: 1, name: 'Liceo' }],
            total: 1
        };
        schoolService.getSchools.mockResolvedValue({ data: mockData });

        await store.fetchSchools();

        expect(store.schools).toHaveLength(1);
        expect(store.schools[0].name).toBe('Liceo');
        expect(store.loading).toBe(false);
    });

    it('handles fetch error', async () => {
        const store = useSchoolStore();
        schoolService.getSchools.mockRejectedValue(new Error('Network Error'));

        await store.fetchSchools();

        expect(store.error).toBe('Network Error');
        expect(store.loading).toBe(false);
    });
});
