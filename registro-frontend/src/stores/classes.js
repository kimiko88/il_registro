import { defineStore } from 'pinia';
import { api } from 'boot/axios';

export const useClassesStore = defineStore('classes', {
    state: () => ({
        classes: [],
        loading: false,
        error: null
    }),

    actions: {
        async fetchClasses() {
            this.loading = true;
            try {
                const response = await api.get('/classes');
                this.classes = response.data;
            } catch (err) {
                this.error = 'Failed to fetch classes';
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
            } catch (err) {
                this.error = 'Failed to delete class';
                throw err;
            }
        }
    }
});
