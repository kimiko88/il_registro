import { defineStore } from 'pinia';
import api from '../services/api';
import { gradeService } from '../services/gradeService';

export const useGradesStore = defineStore('grades', {
    state: () => ({
        grades: null,      // ClassGradesResponse { students: [...] } or MyGradesResponse { semesters: [...] }
        loading: false,
        error: null,
        subjects: [],      // [{ id, name, teacher_id, ... }]
        _cacheMap: {},     // In-memory cache by `${classId}:${subjectId}`
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
        classAverage: (state, getters) => (semester = 0) => {
            return getters.classAverageForSemester(semester);
        },
        classAverageForSemester: (state) => (semester = 0) => {
            if (!state.grades || !state.grades.students) return 0;
            let sum = 0;
            let count = 0;
            state.grades.students.forEach(s => {
                s.grades.forEach(g => {
                    if (semester > 0 && g.semester && Number(g.semester) !== Number(semester)) {
                        return;
                    }
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
        async fetchGrades(classId, subjectId, force = false, isBackgroundRefresh = false) {
            const cacheKey = `${classId}:${subjectId || 'all'}`;
            const cached = this._cacheMap[cacheKey];
            const CACHE_TTL = 120000; // 2 minutes TTL

            if (!force && !isBackgroundRefresh && cached && (Date.now() - cached.timestamp < CACHE_TTL)) {
                this.grades = cached.data;
                this._lastClassId = classId;
                this._lastSubjectId = subjectId;
                return;
            }

            if (!isBackgroundRefresh) {
                this.loading = true;
            }
            this.error = null;
            this._lastClassId = classId;
            this._lastSubjectId = subjectId;
            const currentReqId = ++this._requestId;
            try {
                const response = await gradeService.getByClass(classId, subjectId);
                // Ignore stale response if a newer request was dispatched
                if (currentReqId === this._requestId) {
                    const data = response.data || null;
                    this.grades = data;
                    if (data) {
                        this._cacheMap[cacheKey] = { data, timestamp: Date.now() };
                    }
                }
            } catch (err) {
                if (currentReqId === this._requestId) {
                    this.error = err.response?.data?.error || err.message || 'Errore durante il recupero dei voti';
                    console.error("Error fetching grades:", err);
                }
            } finally {
                if (currentReqId === this._requestId && !isBackgroundRefresh) {
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
                this.grades = response.data;
            } catch (err) {
                this.error = err.response?.data?.error || err.message || 'Errore durante il recupero dei voti';
            } finally {
                this.loading = false;
            }
        },

        async addGrade(gradeData) {
            this.loading = true;
            this.error = null;
            try {
                const response = await gradeService.saveGrade(gradeData);
                const targetKey = `${this._lastClassId}:${this._lastSubjectId || 'all'}`;
                delete this._cacheMap[targetKey];
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId, true, true);
                }
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || err.message || 'Errore durante la registrazione del voto';
                console.error("Error adding grade:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async updateGrade(id, updates) {
            this.loading = true;
            this.error = null;
            try {
                const response = await gradeService.updateGrade(id, updates);
                const targetKey = `${this._lastClassId}:${this._lastSubjectId || 'all'}`;
                delete this._cacheMap[targetKey];
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId, true, true);
                }
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || err.message || 'Errore durante la modifica del voto';
                console.error("Error updating grade:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async deleteGrade(id) {
            this.loading = true;
            this.error = null;
            try {
                const response = await gradeService.deleteGrade(id);
                const targetKey = `${this._lastClassId}:${this._lastSubjectId || 'all'}`;
                delete this._cacheMap[targetKey];
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId, true, true);
                }
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || err.message || 'Errore durante l\'eliminazione del voto';
                console.error("Error deleting grade:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async createClassTest(testData) {
            this.loading = true;
            this.error = null;
            try {
                const response = await gradeService.createClassTest(testData);
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId, true, true);
                }
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || err.message || 'Errore durante la creazione della verifica';
                console.error("Error creating class test:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async updateClassTest(id, testData) {
            this.loading = true;
            this.error = null;
            try {
                const response = await gradeService.updateClassTest(id, testData);
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId, true, true);
                }
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || err.message || 'Errore durante la modifica della verifica';
                console.error("Error updating class test:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async deleteClassTest(id) {
            this.loading = true;
            this.error = null;
            try {
                const response = await gradeService.deleteClassTest(id);
                if (this._lastClassId) {
                    await this.fetchGrades(this._lastClassId, this._lastSubjectId, true, true);
                }
                return response.data;
            } catch (err) {
                this.error = err.response?.data?.error || err.message || 'Errore durante l\'eliminazione della verifica';
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
            this.error = null;
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
                this.error = err.response?.data?.error || err.message || 'Errore durante il download del PDF della pagella';
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
        },

        clearCache() {
            this._cacheMap = {};
            this.grades = null;
            this._lastClassId = null;
            this._lastSubjectId = null;
        }
    }
});
