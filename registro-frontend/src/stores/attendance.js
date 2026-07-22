import { defineStore } from 'pinia';
import api from '../services/api';
import { attendanceService } from '../services/attendanceService';
import { useAuthStore } from './auth';
import { useChildrenStore } from './children';

export const useAttendanceStore = defineStore('attendance', {
    state: () => ({
        records: [],
        loading: false,
        error: null,
        justifications: [], // Pending justification requests
    }),

    getters: {
        presentCount: (state) => state.records.filter(r => r.status === 'Present').length,
        absentCount: (state) => state.records.filter(r => r.status === 'Absent').length,
        lateCount: (state) => state.records.filter(r => r.status === 'Late').length,
    },

    actions: {
        async fetchDailyAttendance(classId, date) {
            this.loading = true;
            this.error = null;
            try {
                const response = await attendanceService.getByClass(classId, date);
                const records = response.data?.records || [];
                this.records = records.map(r => ({
                    studentId: r.student_id,
                    name: r.student_name || `Student (${r.student_id})`,
                    status: r.status,
                    notes: r.notes || '',
                    time: r.entry_time || ''
                }));
            } catch (err) {
                this.error = err.message;
                console.error("Error fetching daily attendance:", err);
            } finally {
                this.loading = false;
            }
        },

        async submitAttendance(classId, date, records, hour = 1, subjectId = '00000000-0000-0000-0000-000000000000') {
            this.loading = true;
            try {
                const payload = {
                    class_id: classId,
                    date: date,
                    hour: hour,
                    subject_id: subjectId,
                    statuses: records.map(r => ({
                        student_id: r.studentId,
                        status: r.status,
                        entry_time: r.status === 'Late' ? r.time : null,
                        exit_time: r.status === 'LeftEarly' ? r.time : null,
                        hour: hour,
                        subject_id: subjectId
                    }))
                };
                const response = await api.post('/attendance/mark-bulk', payload);
                this.records = records;
                return response.data;
            } catch (err) {
                console.error("Error submitting attendance:", err);
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async fetchPendingJustifications(classId) {
            try {
                const response = await api.get('/attendance/pending-justifications', {
                    params: { class_id: classId }
                });
                this.justifications = response.data || [];
            } catch (err) {
                console.error("Error fetching justifications:", err);
            }
        },

        async approveJustification(id) {
            try {
                await api.post(`/attendance/justification/${id}/process`, { approve: true });
                this.justifications = this.justifications.filter(j => j.id !== id);
            } catch (err) {
                console.error("Error approving justification:", err);
                throw err;
            }
        },

        // Student Actions
        async fetchMyAttendance() {
            this.loading = true;
            try {
                const response = await attendanceService.getMyAttendance();
                const data = response.data || [];
                this.records = data.map(r => ({
                    date: r.date,
                    status: r.status,
                    notes: r.notes || '',
                    time: r.entry_time || '',
                    justificationStatus: r.is_justified ? 'Approved' : 'PendingJustification'
                }));
            } catch (err) {
                console.error("Error fetching my attendance:", err);
            } finally {
                this.loading = false;
            }
        },

        async requestJustification(date, reason, studentId) {
            try {
                const authStore = useAuthStore();
                const childrenStore = useChildrenStore();
                
                let targetStudentId = studentId;
                if (!targetStudentId && authStore.user) {
                    targetStudentId = authStore.user.student_id;
                }
                if (!targetStudentId) {
                    targetStudentId = childrenStore.selectedChildId;
                }

                const payload = {
                    student_id: targetStudentId,
                    start_date: date,
                    end_date: date,
                    reason: reason
                };
                await api.post('/attendance/justify', payload);
                const record = this.records.find(r => r.date === date);
                if (record) record.justificationStatus = 'Pending';
            } catch (err) {
                console.error("Error requesting justification:", err);
                throw err;
            }
        },

        /**
         * Fetches the monthly attendance breakdown for a student.
         * @param {string} studentID - The student UUID
         * @param {string} [schoolYear] - Format "2024-2025"; defaults to current school year
         * @param {boolean} [isParent] - If true, uses the parent child-attendance endpoint
         * @returns {Promise<Object>} MonthlyBreakdownResponse with 'months' array
         */
        async fetchMonthlyBreakdown(studentID, schoolYear = '', isParent = false) {
            try {
                const params = schoolYear ? { school_year: schoolYear } : {};
                const url = isParent
                    ? `/attendance/child-attendance/${studentID}/monthly-breakdown`
                    : `/attendance/students/${studentID}/monthly-breakdown`;
                const response = await api.get(url, { params });
                return response.data;
            } catch (err) {
                console.error('Error fetching monthly breakdown:', err);
                throw err;
            }
        }
    }
});

