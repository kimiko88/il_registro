import { defineStore } from 'pinia';

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
                await new Promise(resolve => setTimeout(resolve, 500));
                // Mock Data
                this.children = [
                    { id: 's1', firstName: 'Mario', lastName: 'Rossi', className: '5A', school: 'Liceo Scientifico', avatar: 'https://cdn.quasar.dev/img/boy-avatar.png' },
                    { id: 's3', firstName: 'Sofia', lastName: 'Rossi', className: '3B', school: 'Liceo Classico', avatar: 'https://cdn.quasar.dev/img/avatar6.jpg' }
                ];
                // Auto-select first if none selected
                if (!this.selectedChildId && this.children.length > 0) {
                    this.selectedChildId = this.children[0].id;
                }
            } finally {
                this.loading = false;
            }
        },

        selectChild(id) {
            this.selectedChildId = id;
        }
    }
});
