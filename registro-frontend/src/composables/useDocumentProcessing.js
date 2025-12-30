import { ref } from 'vue';
import { useDocumentsStore } from '../stores/documents';
import { useQuasar } from 'quasar';

export function useDocumentProcessing() {
    const store = useDocumentsStore();
    const $q = useQuasar();

    const showReviewDialog = ref(false);
    const selectedDoc = ref(null);
    const reviewNotes = ref('');

    const openReview = (doc) => {
        selectedDoc.value = doc;
        reviewNotes.value = '';
        showReviewDialog.value = true;
    };

    const submitReview = async (decision) => {
        try {
            await store.reviewDocument(selectedDoc.value.id, decision, reviewNotes.value);
            $q.notify({ type: 'positive', message: `Document ${decision}` });
            showReviewDialog.value = false;
        } catch (e) {
            $q.notify({ type: 'negative', message: 'Review failed' });
        }
    };

    return {
        store,
        showReviewDialog,
        selectedDoc,
        reviewNotes,
        openReview,
        submitReview
    };
}
