import { ref } from 'vue';
import { useDocumentsStore } from 'src/stores/documents';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useDocumentCreation() {
    const documentsStore = useDocumentsStore();
    const $q = useQuasar();

    const isSaving = ref(false);

    const createDraft = async (templateId, studentId, content) => {
        isSaving.value = true;
        const t = i18n?.global?.t;
        try {
            await documentsStore.createDocument({
                templateId,
                studentId,
                content,
                title: 'New Document' // Logic to generate title
            });
            $q.notify({ type: 'positive', message: t ? t('composables.documents.draftCreated') : 'Bozza documento creata' });
            return true;
        } catch (err) {
            $q.notify({ type: 'negative', message: t ? t('composables.documents.draftError') : 'Impossibile creare la bozza' });
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
