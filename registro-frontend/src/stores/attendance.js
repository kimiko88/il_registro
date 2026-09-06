import { defineStore } from 'pinia';
import api from '../services/api';
import { attendanceService } from '../services/attendanceService';
import { useAuthStore } from './auth';
import { useChildrenStore } from './children';
import { i18n } from '@/i18n';

function saveAttendanceCache(classId, date, records) {
    try {
        if (typeof localStorage !== 'undefined' && classId && date && Array.isArray(records)) {
            localStorage.setItem(`registro_attendance_${classId}_${date}`, JSON.stringify(records));
        }
    } catch { /* storage quota */ }
}

function loadAttendanceCache(classId, date) {
    try {
        if (typeof localStorage === 'undefined' || !classId || !date) return null;
        const raw = localStorage.getItem(`registro_attendance_${classId}_${date}`);
        return raw ? JSON.parse(raw) : null;
    } catch {
        return null;
    }
}

export const useAttendanceStore = defineStore('attendance', {
    state: () => ({
        records: [],
        loading: false,
        error: null,
        justifications: [], // Pending justification requests
        _requestId: 0,
        // Tracked to allow reconnect-refresh in websocket.js
        currentClassId: null,
        currentDate: null,
    }),

    getters: {
        presentCount: (state) => state.records.filter(r => (r.status || '').toLowerCase() === 'present').length,
        absentCount: (state) => state.records.filter(r => (r.status || '').toLowerCase() === 'absent').length,
        lateCount: (state) => state.records.filter(r => (r.status || '').toLowerCase() === 'late').length,
        unjustifiedCount: (state) => state.records.filter(r => {
            const status = (r.status || '').toLowerCase();
            const jStatus = (r.justificationStatus || r.justification_status || '').toLowerCase();
            const isJustifiedFlag = r.is_justified || r.isJustified || r.justified;
            if (status !== 'absent') return false;
            if (isJustifiedFlag) return false;
            if (jStatus === 'justified' || jStatus === 'pending' || jStatus === 'pendingapproval' || jStatus === 'pending_approval') return false;
            return true;
        }).length,
    },

    actions: {
        async fetchDailyAttendance(classId, date) {
            this.loading = true;
            this.error = null;
            // Track context so reconnect-refresh (websocket.js) can re-call correctly.
            this.currentClassId = classId;
            this.currentDate = date;
            const currentReqId = ++this._requestId;
            try {
                const response = await attendanceService.getByClass(classId, date);
                if (currentReqId === this._requestId) {
                    const records = response.data?.records || [];
                    const t = i18n?.global?.t;
                    const studentLabel = t ? t('gradesPage.student') : 'Studente';
                    this.records = records.map(r => ({
                        studentId: r.student_id,
                        name: r.student_name || `${studentLabel} (${r.student_id})`,
                        status: r.status,
                        notes: r.notes || '',
                        time: r.entry_time || ''
                    }));
                    saveAttendanceCache(classId, date, this.records);
                }
            } catch (err) {
                if (currentReqId === this._requestId) {
                    const cached = loadAttendanceCache(classId, date);
                    if (cached && Array.isArray(cached) && cached.length > 0) {
                        this.records = cached;
                    } else {
                        const t = i18n?.global?.t;
                        this.error = err.response?.data?.error || err.message || (t ? t('common.error') : 'Errore durante il recupero delle presenze');
                    }
                    console.error("Error fetching daily attendance:", err);
                }
            } finally {
                if (currentReqId === this._requestId) {
                    this.loading = false;
                }
            }
        },

        // Alias used by websocket.js reconnect-refresh (keeps last-known classId + date).
        async fetchAttendance(classId) {
            if (classId && this.currentDate) {
                return this.fetchDailyAttendance(classId, this.currentDate);
            }
        },

        async submitAttendance(classId, date, records, hour = 1, subjectId = null) {
            this.loading = true;
            try {
                const payload = {
                    class_id: classId,
                    date: date,
                    hour: hour,
                    subject_id: subjectId || null,
                    statuses: records.map(r => ({
                        student_id: r.studentId,
                        status: r.status,
                        entry_time: r.status === 'Late' ? r.time : null,
                        exit_time: r.status === 'LeftEarly' ? r.time : null
                    }))
                };
                const response = await api.post('/attendance/mark-bulk', payload);
                if (response.data && Array.isArray(response.data.records)) {
                    this.records = response.data.records;
                } else {
                    this.records = records;
                }
                return response.data;
            } catch (err) {
                const t = i18n?.global?.t;
                this.error = err.response?.data?.error || err.message || (t ? t('common.error') : 'Errore nella registrazione presenze');
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

        async approveJustification(id, classId = null) {
            try {
                await api.post(`/attendance/justification/${id}/process`, { approve: true });
                this.justifications = this.justifications.filter(j => j.id !== id);
                const authStore = useAuthStore();
                if (authStore.userRole === 'student' || authStore.userRole === 'parent') {
                    await this.fetchMyAttendance();
                } else if (classId) {
                    await this.fetchPendingJustifications(classId);
                }
            } catch (err) {
                console.error("Error approving justification:", err);
                throw err;
            }
        },

        // Student Actions
        async fetchMyAttendance() {
            this.loading = true;
            this.error = null;
            try {
                const response = await attendanceService.getMyAttendance();
                const data = response.data || [];
                this.records = data.map(r => ({
                    date: r.date,
                    status: r.status,
                    notes: r.notes || '',
                    time: r.entry_time || '',
                    justificationStatus: r.is_justified ? 'Justified' : (r.parent_justified ? 'PendingApproval' : 'Unjustified')
                }));
            } catch (err) {
                const t = i18n?.global?.t;
                this.error = err.response?.data?.error || err.message || (t ? t('common.error') : 'Errore durante il recupero delle mie presenze');
                console.error("Error fetching my attendance:", err);
            } finally {
                this.loading = false;
            }
        },

        async requestJustification(date, reason, studentId) {
            const cleanReason = (reason || '').trim();
            if (!cleanReason || cleanReason.length < 3) {
                throw new Error("La motivazione della giustificazione deve contenere almeno 3 caratteri");
            }
            if (cleanReason.length > 500) {
                throw new Error("La motivazione della giustificazione non può superare 500 caratteri");
            }

            try {
                const authStore = useAuthStore();
                const childrenStore = useChildrenStore();
                
                let targetStudentId = studentId;

                if (authStore.userRole === 'parent') {
                    const validChildIds = (childrenStore.children || []).map(c => c.id || c.student_id).filter(Boolean);
                    if (targetStudentId && validChildIds.length > 0 && !validChildIds.includes(targetStudentId)) {
                        throw new Error("Impossibile richiedere giustificazione per uno studente non associato");
                    }
                    if (!targetStudentId) {
                        targetStudentId = childrenStore.selectedChildId || validChildIds[0];
                    }
                } else if (!targetStudentId) {
                    targetStudentId = authStore.user?.student_id;
                }

                if (!targetStudentId) {
                    throw new Error("Impossibile determinare lo studente per la richiesta di giustificazione");
                }

                const payload = {
                    student_id: targetStudentId,
                    start_date: date,
                    end_date: date,
                    reason: cleanReason
                };
                const record = this.records.find(r => r.date === date);
                const prevStatus = record ? record.justificationStatus : null;
                if (record) record.justificationStatus = 'Pending';
                try {
                    await api.post('/attendance/justify', payload);
                    await this.fetchMyAttendance();
                } catch (err) {
                    if (record && prevStatus) record.justificationStatus = prevStatus;
                    throw err;
                }
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
        },

        async fetchUnjustified(studentID) {
            try {
                const response = await api.get(`/attendance/child/${studentID}/unjustified`);
                return response.data || [];
            } catch (err) {
                console.error('Error fetching unjustified absences:', err);
                return [];
            }
        },

        async justifyAbsence(studentID, attendanceID, reason, notes = '') {
            try {
                const response = await api.post(`/attendance/child/${studentID}/justify/${attendanceID}`, {
                    reason,
                    notes
                });
                return response.data;
            } catch (err) {
                console.error('Error justifying absence:', err);
                throw err;
            }
        },

        async fetchChildStats(studentID) {
            try {
                const response = await api.get(`/attendance/child/${studentID}/stats`);
                return response.data;
            } catch (err) {
                console.error('Error fetching child attendance stats:', err);
                throw err;
            }
        }
    }
});

