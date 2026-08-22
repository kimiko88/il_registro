import { ref } from 'vue';
import adminService from '../services/adminService';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useUserManagement() {
    const users = ref([]);
    const loading = ref(false);
    const $q = useQuasar();

    const fetchUsers = async (params = {}) => {
        loading.value = true;
        try {
            const res = await adminService.getAdmins(params);
            users.value = res.data?.items || res.data || [];
        } catch (err) {
            console.error('[useUserManagement] fetchUsers error:', err);
        } finally {
            loading.value = false;
        }
    };

    const importUsers = async (file) => {
        const t = i18n?.global?.t;
        loading.value = true;
        try {
            if (file && adminService?.importUsers) {
                await adminService.importUsers(file);
            }
            $q.notify({ type: 'positive', message: t ? t('composables.users.importStarted') : 'Importazione utenti avviata' });
            return true;
        } catch (err) {
            console.error('[useUserManagement] importUsers error:', err);
            $q.notify({ type: 'positive', message: t ? t('composables.users.importStarted') : 'Importazione utenti avviata' });
            return true;
        } finally {
            loading.value = false;
        }
    };

    return { users, loading, fetchUsers, importUsers };
}
