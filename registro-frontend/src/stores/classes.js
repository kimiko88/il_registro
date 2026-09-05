import { defineStore } from 'pinia';
import api from '../services/api';
import { i18n } from '@/i18n';

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

export const CLASSES_CACHE_TTL = 3 * 60 * 1000; // 3 minutes

export const useClassesStore = defineStore('classes', {
    state: () => ({
        classes: [],
        allSchoolClassesAndGroups: [],
        selectedClassId: null,
        loading: false,
        error: null,
        _lastFetchClasses: 0,
        _lastFetchAssigned: 0,
        _lastFetchAll: 0,
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
        invalidateCache() {
            this._lastFetchClasses = 0;
            this._lastFetchAssigned = 0;
            this._lastFetchAll = 0;
        },

        async fetchClasses(params = {}, options = {}) {
            const hasParams = params && Object.keys(params).length > 0;
            const isFresh = !options.force && !hasParams && this.classes.length > 0 && (Date.now() - this._lastFetchClasses < CLASSES_CACHE_TTL);
            if (isFresh) {
                return this.classes;
            }

            this.loading = true;
            try {
                const response = await api.get('/classes', { params });
                const raw = response.data || [];
                this.classes = raw.map(formatClassItem);
                if (!hasParams) {
                    this._lastFetchClasses = Date.now();
                }
                return this.classes;
            } catch (err) {
                this.error = 'Failed to fetch classes';
                console.error(err);
            } finally {
                this.loading = false;
            }
        },

        async fetchAssignedClasses(options = {}) {
            const isFresh = !options.force && this.classes.length > 0 && (Date.now() - this._lastFetchAssigned < CLASSES_CACHE_TTL);
            if (isFresh) {
                return this.classes;
            }

            this.loading = true;
            try {
                const response = await api.get('/teacher/classes');
                const raw = response.data || [];
                this.classes = raw.map(formatClassItem);
                this._lastFetchAssigned = Date.now();
                return this.classes;
            } catch (err) {
                this.error = 'Failed to fetch assigned classes';
                console.error(err);
            } finally {
                this.loading = false;
            }
        },

        async fetchAllSchoolClassesAndGroups(options = {}) {
            const isFresh = !options.force && this.allSchoolClassesAndGroups.length > 0 && (Date.now() - this._lastFetchAll < CLASSES_CACHE_TTL);
            if (isFresh) {
                return this.allSchoolClassesAndGroups;
            }

            this.loading = true;
            try {
                const [cRes, gRes] = await Promise.all([
                    api.get('/classes').catch(err => {
                        console.error('Error fetching /classes in fetchAllSchoolClassesAndGroups:', err);
                        return { data: [] };
                    }),
                    api.get('/groups').catch(err => {
                        console.error('Error fetching /groups in fetchAllSchoolClassesAndGroups:', err);
                        return { data: [] };
                    })
                ]);
                const rawClasses = cRes.data?.classes || cRes.data || [];
                const rawGroups = gRes.data?.groups || gRes.data || [];

                const formattedClasses = rawClasses.map(formatClassItem);
                const t = i18n?.global?.t
                const formattedGroups = rawGroups.map(g => ({
                    id: g.id,
                    name: g.name,
                    section: g.name,
                    articolazione: t ? t('classes.linguisticGroup') : 'Gruppo Linguistico / Articolazione',
                    label: t ? t('classes.linguisticGroupLabel', { name: g.name }) : `Gruppo Linguistico: ${g.name}`,
                    isGroup: true
                }));

                const combined = [...formattedClasses, ...formattedGroups];
                this.allSchoolClassesAndGroups = combined;
                this._lastFetchAll = Date.now();
                return combined;
            } catch (err) {
                console.error('Failed to fetch all school classes and groups:', err);
                return [];
            } finally {
                this.loading = false;
            }
        },

        async createClass(classData) {
            try {
                const response = await api.post('/classes', classData);
                const item = formatClassItem(response.data);
                this.classes.push(item);
                this.invalidateCache();
                return item;
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
                    this.classes[index] = formatClassItem(response.data);
                }
                this.invalidateCache();
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
                this.invalidateCache();
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
