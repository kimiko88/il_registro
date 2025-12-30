import { defineStore } from 'pinia';

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
            try {
                // Mock Data
                await new Promise(resolve => setTimeout(resolve, 500));
                // Mocking records based on students in class (usually fetched from backend)
                this.records = [
                    { studentId: 's1', name: 'Giuseppe Verdi', status: 'Present', notes: '', time: '' },
                    { studentId: 's2', name: 'Mario Rossi', status: 'Absent', notes: '', time: '' },
                    { studentId: 's3', name: 'Sofia Bianchi', status: 'Late', notes: 'Bus delay', time: '08:15' },
                ];
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },

        async submitAttendance(classId, date, records) {
            // Mock API
            await new Promise(resolve => setTimeout(resolve, 500));
            this.records = records; // In real app, re-fetch or update local state
        },

        async fetchPendingJustifications(classId) {
            // Mock justifications
            this.justifications = [
                { id: 1, studentName: 'Mario Rossi', date: '2025-01-15', reason: 'Flu', status: 'Pending' }
            ];
        },

        async approveJustification(id) {
            this.justifications = this.justifications.filter(j => j.id !== id);
        },

        // Student Actions
        async fetchMyAttendance() {
            this.loading = true;
            try {
                // Mock Data
                await new Promise(resolve => setTimeout(resolve, 500));
                // Different structure for personal attendance history
                this.records = [
                    { date: '2025-01-20', status: 'Present', notes: '', time: '' },
                    { date: '2025-01-19', status: 'Absent', notes: '', time: '' },
                    { date: '2025-01-18', status: 'Late', notes: 'Traffic', time: '08:15' },
                    { date: '2025-01-15', status: 'Present', notes: '', time: '' }
                ];
            } finally {
                this.loading = false;
            }
        },

        async requestJustification(date, reason) {
            // Mock API
            await new Promise(resolve => setTimeout(resolve, 400));
            // In a real app, this would be a separate list of "requests", but for now we just log it
            console.log('Justification requested', { date, reason });
            // Update local status mock
            const record = this.records.find(r => r.date === date);
            if (record) record.justificationStatus = 'Pending';
        }
    }
});
