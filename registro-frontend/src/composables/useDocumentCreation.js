import { ref } from 'vue';
import { useDocumentsStore } from 'src/stores/documents';
import { useQuasar } from 'quasar';

export function useDocumentCreation() {
    const documentsStore = useDocumentsStore();
    const $q = useQuasar();

    const isSaving = ref(false);

    const createDraft = async (templateId, studentId, content) => {
        isSaving.value = true;
        try {
            await documentsStore.createDocument({
                templateId,
                studentId,
                content,
                title: 'New Document' // Logic to generate title
            });
            $q.notify({ type: 'positive', message: 'Draft created' });
            return true;
        } catch (err) {
            $q.notify({ type: 'negative', message: 'Failed to create draft' });
            return false;
        } finally {
            isSaving.value = false;
        }
    };

    return {
        createDraft,
        isSaving
    };
}
