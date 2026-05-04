import { setActivePinia, createPinia } from 'pinia';
import { useChildrenStore } from 'src/stores/children';
import { useChildGrades } from 'src/composables/useChildGrades';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { gradeService } from '@/services/gradeService';

// Mock gradeService
vi.mock('@/services/gradeService', () => ({
    gradeService: {
        getMyGrades: vi.fn().mockResolvedValue({ data: [] }),
        getByClass: vi.fn().mockResolvedValue({ data: [] })
    }
}))

describe('Parent Logic', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('manages selected child correctly in children store', async () => {
        const store = useChildrenStore();
        await store.fetchChildren();

        expect(store.children.length).toBeGreaterThan(0);
        expect(store.selectedChildId).not.toBeNull();

        const initialId = store.selectedChildId;
        const otherId = store.children.find(c => c.id !== initialId)?.id;

        if (otherId) {
            store.selectChild(otherId);
            expect(store.selectedChildId).toBe(otherId);
        }
    });

    it('useChildGrades reacts to selected child change', async () => {
        const childrenStore = useChildrenStore();
        // Mock fetch
        childrenStore.children = [{ id: '1' }, { id: '2' }];
        childrenStore.selectedChildId = '1';

        const { selectedChild } = useChildGrades();
        expect(selectedChild.value.id).toBe('1');

        childrenStore.selectChild('2');
        expect(selectedChild.value.id).toBe('2');
    });
});
