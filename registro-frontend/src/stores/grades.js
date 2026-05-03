import { defineStore } from 'pinia';
import { gradeService } from 'src/services/gradeService';

export const useGradesStore = defineStore('grades', {
    state: () => ({
        grades: [],
        loading: false,
        error: null,
        subjects: ['Mathematics', 'Physics', 'History'], // Fallback subjects if not fetched
    }),

    getters: {
        getGradesByStudent: (state) => (studentId) => {
            return state.grades.filter(g => g.studentId === studentId);
        },
        classAverage: (state) => {
            const validGrades = state.grades.filter(g => typeof g.grade_value === 'number');
            if (validGrades.length === 0) return 0;
            const sum = validGrades.reduce((acc, curr) => acc + curr.grade_value, 0);
            return (sum / validGrades.length).toFixed(1);
        }
    },

    actions: {
        // Teacher Actions
        async fetchGrades(classId, subjectId) {
            this.loading = true;
            this.error = null;
            try {
                const response = await gradeService.getByClass(classId, subjectId);
                this.grades = response.data || [];
            } catch (err) {
                this.error = err.message;
                console.error("Error fetching grades:", err);
            } finally {
                this.loading = false;
            }
        },

        // Student Actions
        async fetchMyGrades() {
            this.loading = true;
            this.error = null;
            try {
                const response = await gradeService.getMyGrades();
                // Map the complex semesters response if needed, or store as is
                this.grades = response.data;
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },

        async addGrade(gradeData) {
            this.loading = true;
            try {
                const response = await gradeService.saveGrade(gradeData);
                this.grades.push(response.data);
                return response.data;
            } catch (err) {
                console.error("Error adding grade:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async updateGrade(id, updates) {
            try {
                const response = await gradeService.updateGrade(id, updates);
                const index = this.grades.findIndex(g => g.id === id);
                if (index !== -1) {
                    this.grades[index] = response.data;
                }
            } catch (err) {
                console.error("Error updating grade:", err);
                throw err;
            }
        },

        async deleteGrade(id) {
            try {
                await gradeService.deleteGrade(id);
                this.grades = this.grades.filter(g => g.id !== id);
            } catch (err) {
                console.error("Error deleting grade:", err);
                throw err;
            }
        }
    }
});
