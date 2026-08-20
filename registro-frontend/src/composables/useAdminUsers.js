import { useAdminStore } from '../stores/admin';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useAdminUsers() {
    const store = useAdminStore();
    const $q = useQuasar();

    const onRequest = (props) => store.fetchAdmins(props.pagination);

    const openCreate = () => { /* Logic */ };

    const confirmDelete = (id) => {
        const t = i18n?.global?.t;
        $q.dialog({
            title: t ? t('composables.adminUsers.confirmTitle') : 'Conferma eliminazione',
            message: t ? t('composables.adminUsers.confirmMsg') : 'Sei sicuro di voler eliminare questo amministratore?',
            cancel: true
        }).onOk(() => {
            store.deleteAdmin(id);
        });
    };

    const resetPassword = (_id) => {
        // Logic
    };

    return { openCreate, confirmDelete, resetPassword, onRequest };
}
