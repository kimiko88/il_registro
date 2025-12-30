import { computed } from 'vue';
import { useDocumentsStore } from 'src/stores/documents';
import { useQuasar } from 'quasar';

export function useMyDocuments() {
    const store = useDocumentsStore();
    const $q = useQuasar();

    const downloadFile = (file) => {
        // Mock Download
        $q.notify({ type: 'positive', message: `Downloading ${file.title}...` });
    };

    return {
        documents: computed(() => store.inbox),
        loading: computed(() => store.loading),
        downloadFile
    };
}
