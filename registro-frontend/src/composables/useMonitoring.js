import { onMounted, onUnmounted, getCurrentInstance } from 'vue';
import { useMonitoringStore } from '../stores/monitoring';

export function useMonitoring(autoStart = true, pollIntervalMs = 30000) {
    const store = useMonitoringStore();
    let interval = null;

    const refreshAll = () => {
        store.fetchHealth();
        store.fetchMetrics();
        store.fetchActiveUsers();
    };

    const startPolling = () => {
        stopPolling();
        refreshAll();
        interval = setInterval(() => {
            store.fetchHealth();
            store.fetchActiveUsers();
        }, pollIntervalMs);
    };

    const stopPolling = () => {
        if (interval) {
            clearInterval(interval);
            interval = null;
        }
    };

    if (getCurrentInstance()) {
        onMounted(() => {
            if (autoStart) {
                startPolling();
            }
        });

        onUnmounted(() => {
            stopPolling();
        });
    }

    return { store, refreshAll, startPolling, stopPolling };
}
