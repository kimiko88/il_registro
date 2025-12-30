import { defineStore } from 'pinia';
import schoolService from '../services/schoolService';

export const useSchoolStore = defineStore('schools', {
    state: () => ({
        schools: [],
        currentSchool: null,
        loading: false,
        error: null,
        pagination: {
            page: 1,
            rowsPerPage: 10,
            rowsNumber: 0,
        },
    }),
    actions: {
        async fetchSchools(params = {}) {
            this.loading = true;
            try {
                const response = await schoolService.getSchools({
                    page: params.page || this.pagination.page,
                    limit: params.rowsPerPage || this.pagination.rowsPerPage,
                    ...params
                });
                this.schools = response.data.items;
                this.pagination.rowsNumber = response.data.total;
                this.pagination.page = params.page || this.pagination.page;
                this.pagination.rowsPerPage = params.rowsPerPage || this.pagination.rowsPerPage;
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },
        async createSchool(data) {
            await schoolService.createSchool(data);
            await this.fetchSchools();
        },
        async updateSchool(id, data) {
            await schoolService.updateSchool(id, data);
            await this.fetchSchools();
        },
        async deleteSchool(id) {
            await schoolService.deleteSchool(id);
            await this.fetchSchools();
        },
    },
});
