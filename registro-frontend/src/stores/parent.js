import { defineStore } from 'pinia';
import { api } from 'boot/axios';

export const useParentStore = defineStore('parent', {
    state: () => ({
        children: [], // List of { id, name, class, school }
        selectedChild: null,
        loading: false,
        error: null
    }),

    getters: {
        currentChildId: (state) => state.selectedChild?.id
    },

    actions: {
        async fetchChildren() {
            this.loading = true;
            try {
                const response = await api.get('/users/me/children');
                // Map backend response if needed, for instance formatting name
                this.children = response.data.map(c => ({
                    id: c.id, // Keeping Profile ID or User ID? Backend sends both. Frontend usually needs Profile ID for queries?
                    // The Backend sends: ID (Student ID), UserID, FirstName, LastName, Class, SchoolName.
                    // Let's ensure we use the student PROFILE ID for grades queries if grades service expects it. 
                    // Grades service usually takes Student ID.
                    name: `${c.first_name} ${c.last_name}`,
                    class: c.class,
                    school: c.school_name,
                    userId: c.user_id // Keep ref
                }));

                if (!this.selectedChild && this.children.length > 0) {
                    this.selectedChild = this.children[0];
                }
            } catch (err) {
                console.error(err);
                this.error = 'Failed to fetch children';
            } finally {
                this.loading = false;
            }
        },

        selectChild(child) {
            this.selectedChild = child;
        }
    }
});
