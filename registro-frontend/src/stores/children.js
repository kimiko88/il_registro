import { defineStore } from 'pinia';
import api from '../services/api';

export const useChildrenStore = defineStore('children', {
    state: () => ({
        children: [],
        selectedChildId: null,
        loading: false
    }),

    getters: {
        selectedChild: (state) => state.children.find(c => c.id === state.selectedChildId),
        hasMultipleChildren: (state) => state.children.length > 1
    },

    actions: {
        async fetchChildren() {
            this.loading = true;
            try {
                const response = await api.get('/users/me/children');
                this.children = (response.data || []).map(c => ({
                    id: c.id,
                    userId: c.user_id,
                    firstName: c.first_name,
                    lastName: c.last_name,
                    className: c.class,
                    school: c.school_name,
                    avatar: 'https://cdn.quasar.dev/img/boy-avatar.png'
                }));
                // Auto-select first if none selected
                if (!this.selectedChildId && this.children.length > 0) {
                    this.selectedChildId = this.children[0].id;
                }
            } catch (err) {
                console.error("Error fetching children:", err);
            } finally {
                this.loading = false;
            }
        },

        selectChild(id) {
            this.selectedChildId = id;
        }
    }
});
