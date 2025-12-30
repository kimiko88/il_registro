import { onMounted, onUnmounted } from 'vue';
import { useMonitoringStore } from '../stores/monitoring';

export function useMonitoring() {
    const store = useMonitoringStore();
    let interval;

    onMounted(() => {
        store.fetchHealth();
        store.fetchMetrics();
        store.fetchActiveUsers();

        // Poll every 30 seconds
        interval = setInterval(() => {
            store.fetchHealth();
            store.fetchActiveUsers();
        }, 30000);
    });

    onUnmounted(() => {
        clearInterval(interval);
    });

    return { store };
}
