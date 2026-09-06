import { defineStore } from 'pinia';
import schoolService from '../services/schoolService';

export const SCHOOLS_CACHE_TTL = 5 * 60 * 1000; // 5 minutes

export const useSchoolStore = defineStore('schools', {
    state: () => ({
        schools: [],
        currentSchool: null,
        loading: false,
        error: null,
        _lastFetch: 0,
        pagination: {
            page: 1,
            rowsPerPage: 10,
            rowsNumber: 0,
        },
    }),
    actions: {
        invalidateCache() {
            this._lastFetch = 0;
        },
        async fetchSchools(params = {}, options = {}) {
            const hasCustomParams = params && Object.keys(params).length > 0;
            const isFresh = !options.force && !hasCustomParams && this.schools.length > 0 && (Date.now() - this._lastFetch < SCHOOLS_CACHE_TTL);
            if (isFresh) {
                return this.schools;
            }

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
                if (!hasCustomParams) {
                    this._lastFetch = Date.now();
                }
                return this.schools;
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },
        async createSchool(data) {
            await schoolService.createSchool(data);
            this.invalidateCache();
            await this.fetchSchools({}, { force: true });
        },
        async updateSchool(id, data) {
            await schoolService.updateSchool(id, data);
            this.invalidateCache();
            await this.fetchSchools({}, { force: true });
        },
        async deleteSchool(id) {
            await schoolService.deleteSchool(id);
            this.invalidateCache();
            await this.fetchSchools({}, { force: true });
        },
    },
});
