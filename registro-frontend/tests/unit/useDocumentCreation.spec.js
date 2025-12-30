import { setActivePinia, createPinia } from 'pinia';
import { useDocumentCreation } from 'src/composables/useDocumentCreation';
import { useDocumentsStore } from 'src/stores/documents';
import { describe, it, expect, vi, beforeEach } from 'vitest';

// Mock Quasar
const mockNotify = vi.fn();
vi.mock('quasar', () => ({
    useQuasar: () => ({
        notify: mockNotify
    })
}));

describe('useDocumentCreation', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        mockNotify.mockClear();
    });

    it('creates draft document', async () => {
        const { createDraft } = useDocumentCreation();
        const store = useDocumentsStore();

        // Mock store action
        store.createDocument = vi.fn();
        store.inbox = []; // Ensure inbox exists

        await createDraft('t1', 's1', 'content');

        expect(store.createDocument).toHaveBeenCalled();
        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }));
    });
});
