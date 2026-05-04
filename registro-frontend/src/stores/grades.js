import { defineStore } from 'pinia';
import api from '../services/api';
import { gradeService } from 'src/services/gradeService';

export const useGradesStore = defineStore('grades', {
    state: () => ({
        grades: [],
        loading: false,
        error: null,
        subjects: [], // Now stored as [{id, name, teacher_id}]
    }),

    getters: {
        getGradesByStudent: (state) => (studentId) => {
            if (!state.grades || !state.grades.students) return [];
            const student = state.grades.students.find(s => s.student_id === studentId);
            return student ? student.grades : [];
        },
        classAverage: (state) => {
            if (!state.grades || !state.grades.students) return 0;
            let sum = 0;
            let count = 0;
            state.grades.students.forEach(s => {
                s.grades.forEach(g => {
                    if (typeof g.grade_value === 'number') {
                        sum += g.grade_value;
                        count++;
                    }
                });
            });
            if (count === 0) return 0;
            return (sum / count).toFixed(1);
        }
    },

    actions: {
        async fetchClassSubjects(classId) {
            try {
                const response = await api.get(`/classes/${classId}/subjects`);
                this.subjects = response.data || [];
            } catch (err) {
                console.error("Error fetching class subjects:", err);
            }
        },
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
