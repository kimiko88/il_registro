import { computed } from 'vue';
import { useChildrenStore } from 'src/stores/children';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useChildrenManagement() {
    const store = useChildrenStore();
    const $q = useQuasar();

    const addChild = async (_childData) => {
        const t = i18n?.global?.t;
        // Mock Add
        $q.notify({ type: 'positive', message: t ? t('composables.children.addSuccess') : 'Figlio aggiunto con successo' });
    };

    const removeChild = async (_childId) => {
        const t = i18n?.global?.t;
        $q.dialog({
            title: t ? t('composables.children.removeConfirmTitle') : 'Conferma rimozione',
            message: t ? t('composables.children.removeConfirmMsg') : 'Sei sicuro di voler rimuovere questo profilo studente?',
            cancel: true,
            persistent: true
        }).onOk(() => {
            $q.notify({ type: 'positive', message: t ? t('composables.children.removeSuccess') : 'Profilo studente rimosso' });
        });
    };

    return {
        children: computed(() => store.children),
        addChild,
        removeChild
    };
}
