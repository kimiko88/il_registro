import { ref, onMounted, getCurrentInstance } from 'vue';
import api from '../services/api';

export function useGlobalAnalytics() {
    const stats = ref({
        schools: 0,
        students: 0,
        apiCalls: 0
    });
    const loading = ref(false);

    const fetchStats = async () => {
        loading.value = true;
        try {
            const res = await api.get('/monitoring/analytics');
            if (res?.data) {
                stats.value = res.data;
            }
        } catch (err) {
            console.error('[useGlobalAnalytics] Failed to fetch stats:', err);
        } finally {
            loading.value = false;
        }
    };

    if (getCurrentInstance()) {
        onMounted(fetchStats);
    }

    return { stats, loading, fetchStats };
}
