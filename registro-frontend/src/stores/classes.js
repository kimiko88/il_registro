import { defineStore } from 'pinia';
import api from '../services/api';

export const useClassesStore = defineStore('classes', {
    state: () => ({
        classes: [],
        selectedClassId: null,
        loading: false,
        error: null
    }),

    actions: {
        async fetchClasses(params = {}) {
            this.loading = true;
            try {
                const response = await api.get('/classes', { params });
                this.classes = response.data;
            } catch (err) {
                this.error = 'Failed to fetch classes';
                console.error(err);
            } finally {
                this.loading = false;
            }
        },

        async fetchAssignedClasses() {
            this.loading = true;
            try {
                // Endpoint for classes assigned to the current teacher
                const response = await api.get('/teacher/classes');
                this.classes = response.data;
            } catch (err) {
                this.error = 'Failed to fetch assigned classes';
                console.error(err);
            } finally {
                this.loading = false;
            }
        },

        async createClass(classData) {
            try {
                const response = await api.post('/classes', classData);
                this.classes.push(response.data);
                return response.data;
            } catch (err) {
                this.error = 'Failed to create class';
                throw err;
            }
        },

        async updateClass(id, classData) {
            try {
                const response = await api.put(`/classes/${id}`, classData);
                const index = this.classes.findIndex(c => c.id === id);
                if (index !== -1) {
                    this.classes[index] = response.data;
                }
            } catch (err) {
                this.error = 'Failed to update class';
                throw err;
            }
        },

        async deleteClass(id) {
            try {
                await api.delete(`/classes/${id}`);
                this.classes = this.classes.filter(c => c.id !== id);
                if (this.selectedClassId === id) {
                    this.selectedClassId = null;
                }
            } catch (err) {
                this.error = 'Failed to delete class';
                throw err;
            }
        },

        selectClass(id) {
            this.selectedClassId = id;
        }
    }
});
