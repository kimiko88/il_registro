import { defineStore } from 'pinia';
import adminService from '../services/adminService';

export const useAdminStore = defineStore('admin', {
    state: () => ({
        admins: [],
        loading: false,
        error: null,
        pagination: {
            page: 1,
            rowsPerPage: 10,
            rowsNumber: 0,
        }
    }),
    actions: {
        async fetchAdmins(params = {}) {
            this.loading = true;
            this.error = null;
            try {
                const response = await adminService.getAdmins({
                    page: params.page || this.pagination.page,
                    limit: params.rowsPerPage || this.pagination.rowsPerPage,
                    ...params
                });
                this.admins = response.data.items;
                this.pagination.rowsNumber = response.data.total;
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },
        // CRUD actions delegated to service but handled here for UI state
    }
});
