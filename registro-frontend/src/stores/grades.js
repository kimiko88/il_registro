import { defineStore } from 'pinia';
import api from '../services/api';
import { gradeService } from 'src/services/gradeService';

export const useGradesStore = defineStore('grades', {
    state: () => ({
        grades: null,      // ClassGradesResponse { students: [...] } or MyGradesResponse { semesters: [...] }
        loading: false,
        error: null,
        subjects: [],      // [{ id, name, teacher_id, ... }]
        // Track last fetch context to enable refetch after mutations
        _lastClassId: null,
        _lastSubjectId: null,
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
                    if (typeof g.grade_value === 'number' && g.grade_value >= 0) {
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
            this._lastClassId = classId;
            this._lastSubjectId = subjectId;
            try {
                const response = await gradeService.getByClass(classId, subjectId);
                this.grades = response.data || null;
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
                // Refetch to keep the full class view consistent
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId);
                }
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
                // Refetch to keep the full class view consistent
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId);
                }
                return response.data;
            } catch (err) {
                console.error("Error updating grade:", err);
                throw err;
            }
        },

        async deleteGrade(id) {
            try {
                await gradeService.deleteGrade(id);
                // Refetch to keep the full class view consistent
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId);
                }
            } catch (err) {
                console.error("Error deleting grade:", err);
                throw err;
            }
        },

        async createClassTest(testData) {
            this.loading = true;
            try {
                const response = await gradeService.createClassTest(testData);
                // Refetch grades to show the new grades in the register
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId);
                }
                return response.data;
            } catch (err) {
                console.error("Error creating class test:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async updateClassTest(id, testData) {
            this.loading = true;
            try {
                const response = await gradeService.updateClassTest(id, testData);
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId);
                }
                return response.data;
            } catch (err) {
                console.error("Error updating class test:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async deleteClassTest(id) {
            this.loading = true;
            try {
                const response = await gradeService.deleteClassTest(id);
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId);
                }
                return response.data;
            } catch (err) {
                console.error("Error deleting class test:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async fetchSemesterReport(semester = 1) {
            this.loading = true;
            this.error = null;
            try {
                const response = await api.get(`/grades/my-grades/semester/${semester}`);
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || 'Errore durante il recupero della pagella';
                console.error(err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async downloadReportCardPDF(semester = 1) {
            try {
                const response = await api.get(`/grades/my-grades/semester/${semester}/pdf`, {
                    responseType: 'blob'
                });
                const blob = new Blob([response.data], { type: 'application/pdf' });
                const url = window.URL.createObjectURL(blob);
                const link = document.createElement('a');
                link.href = url;
                link.setAttribute('download', `pagella_q${semester}.pdf`);
                document.body.appendChild(link);
                link.click();
                link.remove();
                window.URL.revokeObjectURL(url);
            } catch (err) {
                console.error("Error downloading report card PDF:", err);
                throw err;
            }
        }
    }
});
