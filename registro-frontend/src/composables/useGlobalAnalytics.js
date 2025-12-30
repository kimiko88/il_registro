import { ref, onMounted } from 'vue';
import api from '../services/api';

export function useGlobalAnalytics() {
    const stats = ref({
        schools: 0,
        students: 0,
        apiCalls: 0
    });

    const fetchStats = async () => {
        const res = await api.get('/monitoring/analytics');
        stats.value = res.data;
    };

    onMounted(fetchStats);

    return { stats };
}
