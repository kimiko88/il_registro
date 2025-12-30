import { defineStore } from 'pinia';

export const useClassesStore = defineStore('classes', {
    state: () => ({
        classes: [], // List of classes assigned to teacher
        selectedClassId: null, // Currently active class ID
        currentClassDetails: null, // Full details (students, etc.)
        loading: false,
        error: null
    }),

    getters: {
        selectedClass: (state) => state.classes.find(c => c.id === state.selectedClassId),
        hasSelectedClass: (state) => !!state.selectedClassId,
    },

    actions: {
        async fetchAssignedClasses() {
            this.loading = true;
            try {
                // Mock Data
                await new Promise(resolve => setTimeout(resolve, 600));
                this.classes = [
                    { id: '1A', name: '1A', type: 'Scientifico', studentsCount: 20, coordinator: false },
                    { id: '2B', name: '2B', type: 'Linguistico', studentsCount: 22, coordinator: true }, // Coordinator example
                    { id: '5A', name: '5A', type: 'Scientifico', studentsCount: 18, coordinator: false }
                ];

                // Auto-select first if none selected
                if (!this.selectedClassId && this.classes.length > 0) {
                    this.selectedClassId = this.classes[0].id;
                }
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },

        selectClass(classId) {
            if (classId === this.selectedClassId) return;
            this.selectedClassId = classId;
            this.fetchClassDetails(classId); // Fetch details when selected
        },

        async fetchClassDetails(classId) {
            // In real app, fetch students list here
            // For now, mock it strictly when needed by components
            console.log(`Fetching details for class ${classId}`);
        }
    }
});
