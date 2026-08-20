import { ref } from 'vue';
import adminService from '../services/adminService'; // Reuse or create specialized service
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useUserManagement() {
    const users = ref([]);
    const loading = ref(false);
    const $q = useQuasar();

    const fetchUsers = async () => {
        loading.value = true;
        try {
            const res = await adminService.getAdmins({}); // Placeholder, should accept role filter
            users.value = res.data?.items || res.data || [];
        } catch (err) {
            console.error('[useUserManagement] fetchUsers error:', err);
        } finally {
            loading.value = false;
        }
    };

    const importUsers = async (_file) => {
        const t = i18n?.global?.t;
        // Mock import
        $q.notify({ type: 'positive', message: t ? t('composables.users.importStarted') : 'Importazione utenti avviata' });
    };

    return { users, loading, fetchUsers, importUsers };
}
