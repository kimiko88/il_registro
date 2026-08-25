import api from './api';

export const monitoringService = {
    getSystemHealth() {
        return api.get('/monitoring/health');
    },
    getMetrics(range) {
        return api.get('/monitoring/metrics', { params: { range } });
    },
    getActiveUsers() {
        return api.get('/monitoring/active-users');
    },
    getErrorLogs(params) {
        return api.get('/monitoring/errors', { params });
    },
    getAnalytics() {
        return api.get('/monitoring/analytics');
    }
};

export default monitoringService;

