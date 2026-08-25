import { ref } from 'vue';
import { useDocumentsStore } from '../stores/documents';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

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
        const t = i18n?.global?.t;
        try {
            await store.reviewDocument(selectedDoc.value.id, decision, reviewNotes.value);
            $q.notify({ type: 'positive', message: t ? t('composables.documents.reviewSuccess', { decision }) : `Documento ${decision}` });
            showReviewDialog.value = false;
        } catch (e) {
            $q.notify({ type: 'negative', message: t ? t('composables.documents.reviewError') : 'Revisione fallita' });
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
