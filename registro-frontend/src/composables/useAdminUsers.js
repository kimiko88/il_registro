import { useAdminStore } from '../stores/admin';
import { useQuasar } from 'quasar';

export function useAdminUsers() {
    const store = useAdminStore();
    const $q = useQuasar();

    const onRequest = (props) => store.fetchAdmins(props.pagination);

    const openCreate = () => { /* Logic */ };

    const confirmDelete = (id) => {
        $q.dialog({ title: 'Confirm', message: 'Delete admin?', cancel: true }).onOk(() => {
            store.deleteAdmin(id);
        });
    };

    const resetPassword = (id) => {
        // Logic
    };

    return { openCreate, confirmDelete, resetPassword, onRequest };
}
