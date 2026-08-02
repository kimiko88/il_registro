import { defineStore } from 'pinia';
import api from '../services/api';
import { gradeService } from '../services/gradeService';

export const useGradesStore = defineStore('grades', {
    state: () => ({
        grades: null,      // ClassGradesResponse { students: [...] } or MyGradesResponse { semesters: [...] }
        loading: false,
        error: null,
        subjects: [],      // [{ id, name, teacher_id, ... }]
        // Track last fetch context to enable refetch after mutations
        _lastClassId: null,
        _lastSubjectId: null,
        _requestId: 0,
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
                    if (typeof g.grade_value === 'number' && g.grade_value > 0) {
                        sum += g.grade_value;
                        count++;
                    }
                });
            });
            if (count === 0) return 0;
            return Math.round((sum / count) * 10) / 10;
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
            const currentReqId = ++this._requestId;
            try {
                const response = await gradeService.getByClass(classId, subjectId);
                // Ignore stale response if a newer request was dispatched
                if (currentReqId === this._requestId) {
                    this.grades = response.data || null;
                }
            } catch (err) {
                if (currentReqId === this._requestId) {
                    this.error = {
                        message: err.response?.data?.error || err.message,
                        status: err.response?.status || 500
                    };
                    console.error("Error fetching grades:", err);
                }
            } finally {
                if (currentReqId === this._requestId) {
                    this.loading = false;
                }
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
                this.error = {
                    message: err.response?.data?.error || err.message,
                    status: err.response?.status || 500
                };
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
            this.loading = true;
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
            } finally {
                this.loading = false;
            }
        },

        async deleteGrade(id) {
            this.loading = true;
            try {
                await gradeService.deleteGrade(id);
                // Refetch to keep the full class view consistent
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId);
                }
            } catch (err) {
                console.error("Error deleting grade:", err);
                throw err;
            } finally {
                this.loading = false;
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
                const response = await gradeService.getSemesterReport(semester);
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
            this.loading = true;
            let url = null;
            let link = null;
            try {
                const response = await gradeService.downloadReportCardPDF(semester);
                const blob = new Blob([response.data], { type: 'application/pdf' });
                url = window.URL.createObjectURL(blob);
                link = document.createElement('a');
                link.href = url;
                link.setAttribute('download', `pagella_q${semester}.pdf`);
                document.body.appendChild(link);
                link.click();
            } catch (err) {
                console.error("Error downloading report card PDF:", err);
                throw err;
            } finally {
                if (link && link.parentNode) {
                    link.remove();
                }
                if (url) {
                    window.URL.revokeObjectURL(url);
                }
                this.loading = false;
            }
        }
    }
});
