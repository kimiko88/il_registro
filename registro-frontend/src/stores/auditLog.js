import { defineStore } from 'pinia';
import api from 'src/services/api';

export const useAuditLogStore = defineStore('auditLog', {
    state: () => ({
        logs: [],
        total: 0,
        page: 1,
        limit: 20,
        totalPages: 1,
        loading: false,
        error: null
    }),

    actions: {
        async fetchLogs(filters = {}, page = 1, limit = 20) {
            this.loading = true;
            this.error = null;
            try {
                const params = new URLSearchParams({
                    page: page.toString(),
                    limit: limit.toString()
                });
                if (filters.actor_id) params.append('actor_id', filters.actor_id);
                if (filters.action && filters.action !== 'all') params.append('action', filters.action);
                if (filters.entity_type) params.append('entity_type', filters.entity_type);
                if (filters.from) params.append('from', filters.from);
                if (filters.to) params.append('to', filters.to);

                const response = await api.get(`/audit-log?${params.toString()}`);
                const data = response.data || {};
                this.logs = data.data || [];
                this.total = data.total || 0;
                this.page = data.page || page;
                this.limit = data.limit || limit;
                this.totalPages = data.total_pages || 1;
                return data;
            } catch (err) {
                this.error = err.response?.data?.error || 'Errore caricamento audit log';
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async exportLogs(filters = {}) {
            const params = new URLSearchParams();
            if (filters.actor_id) params.append('actor_id', filters.actor_id);
            if (filters.action && filters.action !== 'all') params.append('action', filters.action);
            if (filters.entity_type) params.append('entity_type', filters.entity_type);
            if (filters.from) params.append('from', filters.from);
            if (filters.to) params.append('to', filters.to);

            const response = await api.get(`/audit-log/export?${params.toString()}`, {
                responseType: 'blob'
            });
            const blob = new Blob([response.data], { type: 'text/csv' });
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', `audit_logs_${new Date().toISOString().slice(0, 10)}.csv`);
            document.body.appendChild(link);
            link.click();
            link.remove();
            window.URL.revokeObjectURL(url);
        }
    }
});
