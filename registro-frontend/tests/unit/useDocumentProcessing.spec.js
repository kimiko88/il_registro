import { setActivePinia, createPinia } from 'pinia';
import { useDocumentProcessing } from 'src/composables/useDocumentProcessing';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useDocumentsStore } from 'src/stores/documents';

// Mock Quasar
const mockNotify = vi.fn();
vi.mock('quasar', () => ({
    useQuasar: () => ({
        notify: mockNotify
    })
}));

describe('useDocumentProcessing', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        mockNotify.mockClear();
    });

    it('opens review dialog', () => {
        const { openReview, showReviewDialog, selectedDoc } = useDocumentProcessing();
        const doc = { id: 1, title: 'Test Doc' };

        openReview(doc);

        expect(showReviewDialog.value).toBe(true);
        expect(selectedDoc.value).toEqual(doc);
    });

    it('submits review', async () => {
        const { submitReview, selectedDoc, reviewNotes } = useDocumentProcessing();
        const store = useDocumentsStore();
        store.reviewDocument = vi.fn();

        selectedDoc.value = { id: 1 };
        reviewNotes.value = 'Looks good';

        await submitReview('approved');

        expect(store.reviewDocument).toHaveBeenCalledWith(1, 'approved', 'Looks good');
        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }));
    });
});
