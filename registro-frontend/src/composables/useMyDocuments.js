import { computed } from 'vue';
import { useDocumentsStore } from 'src/stores/documents';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useMyDocuments() {
    const store = useDocumentsStore();
    const $q = useQuasar();

    const downloadFile = (file) => {
        const t = i18n?.global?.t;
        const title = file?.title || file?.name || 'documento';
        // Mock Download
        $q.notify({ type: 'positive', message: t ? t('composables.documents.downloading', { title }) : `Download di ${title} in corso...` });
    };

    return {
        documents: computed(() => store.inbox),
        loading: computed(() => store.loading),
        downloadFile
    };
}
