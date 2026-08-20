import { ref, onMounted } from 'vue';
import api from '../services/api';

export function useGlobalAnalytics() {
    const stats = ref({
        schools: 0,
        students: 0,
        apiCalls: 0
    });

    const fetchStats = async () => {
        try {
            const res = await api.get('/monitoring/analytics');
            if (res?.data) {
                stats.value = res.data;
            }
        } catch (err) {
            console.error('[useGlobalAnalytics] Failed to fetch stats:', err);
        }
    };

    onMounted(fetchStats);

    return { stats };
}
