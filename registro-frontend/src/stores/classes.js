import { defineStore } from 'pinia';
import api from '../services/api';

function formatClassItem(c) {
    if (!c) return c;
    let nameText = c.name || `Classe ${c.id}`;
    if (c.section && !nameText.endsWith(c.section)) {
        nameText += c.section;
    }
    let parts = [nameText];
    if (c.articolazione) {
        parts.push(c.articolazione);
    }
    let fullLabel = parts.join(' - ');
    if (c.academic_year) {
        fullLabel += ` (${c.academic_year})`;
    }
    const studentCount = c.students ?? c.students_count ?? 0;
    return {
        ...c,
        students: studentCount,
        students_count: studentCount,
        label: fullLabel,
        displayName: fullLabel
    };
}

export const useClassesStore = defineStore('classes', {
    state: () => ({
        classes: [],
        selectedClassId: null,
        loading: false,
        error: null
    }),

    getters: {
        classOptions: (state) => {
            return state.classes.map(c => ({
                ...c,
                label: c.label || c.displayName || c.name,
                value: c.id
            }));
        },
        formatClassLabel: () => (c) => {
            if (!c) return '';
            let label = c.name || '';
            if (c.section && !label.endsWith(c.section)) {
                label += c.section;
            }
            if (c.articolazione) {
                label += ` - ${c.articolazione}`;
            }
            return label;
        }
    },

    actions: {
        async fetchClasses(params = {}) {
            this.loading = true;
            try {
                const response = await api.get('/classes', { params });
                const raw = response.data || [];
                this.classes = raw.map(formatClassItem);
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
                const response = await api.get('/teacher/classes');
                const raw = response.data || [];
                this.classes = raw.map(formatClassItem);
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
