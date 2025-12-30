import { defineStore } from 'pinia';

export const useGradesStore = defineStore('grades', {
    state: () => ({
        grades: [],
        loading: false,
        error: null,
        subjects: ['Mathematics', 'Physics', 'History'], // Mock subjects
    }),

    getters: {
        getGradesByStudent: (state) => (studentId) => {
            return state.grades.filter(g => g.studentId === studentId);
        },
        classAverage: (state) => {
            if (state.grades.length === 0) return 0;
            const sum = state.grades.reduce((acc, curr) => acc + curr.value, 0);
            return (sum / state.grades.length).toFixed(1);
        }
    },

    actions: {
        // Teacher Actions
        async fetchGrades(classId, subject) {
            this.loading = true;
            try {
                // Mock Data for Teacher
                await new Promise(resolve => setTimeout(resolve, 500));
                this.grades = [
                    { id: 1, studentId: 's1', value: 8, type: 'Written', date: '2025-01-10', description: 'Algebra Test' },
                    { id: 2, studentId: 's1', value: 7.5, type: 'Oral', date: '2025-01-15', description: 'Interrogation' },
                    { id: 3, studentId: 's2', value: 6, type: 'Written', date: '2025-01-10', description: 'Algebra Test' },
                ];
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },

        // Student Actions
        async fetchMyGrades(studentId) {
            this.loading = true;
            try {
                // Mock Data for Student
                await new Promise(resolve => setTimeout(resolve, 600));
                this.grades = [ // Reusing 'grades' state is fine if we consider the store context (Student view vs Teacher view)
                    { id: 1, subject: 'Mathematics', value: 8, type: 'Written', date: '2025-01-10', description: 'Algebra Test' },
                    { id: 2, subject: 'Mathematics', value: 7.5, type: 'Oral', date: '2025-01-15', description: 'Polynomials' },
                    { id: 4, subject: 'Physics', value: 9, type: 'Written', date: '2025-01-12', description: 'Thermodynamics' },
                    { id: 5, subject: 'History', value: 7, type: 'Oral', date: '2025-01-18', description: 'WWII' },
                    { id: 6, subject: 'Physics', value: 8.5, type: 'Lab', date: '2025-01-20', description: 'Heat Capacity' }
                ];
            } finally {
                this.loading = false;
            }
        },

        async addGrade(gradeData) {
            // Mock API call
            await new Promise(resolve => setTimeout(resolve, 300));
            const newGrade = {
                ...gradeData,
                id: Math.random().toString(36).substr(2, 9)
            };
            this.grades.push(newGrade);
            return newGrade;
        },

        async updateGrade(id, updates) {
            const index = this.grades.findIndex(g => g.id === id);
            if (index !== -1) {
                this.grades[index] = { ...this.grades[index], ...updates };
            }
        },

        async deleteGrade(id) {
            this.grades = this.grades.filter(g => g.id !== id);
        }
    }
});
