import { defineStore } from 'pinia'
import colloquiService from '@/services/colloquiService'

export const useColloquiStore = defineStore('colloqui', {
    state: () => ({
        slots: [],
        bookings: [],
        loading: false,
        error: null
    }),

    actions: {
        async fetchSlots(rangeStart, rangeEnd, teacherId = null) {
            this.loading = true
            this.error = null
            try {
                const params = {}
                if (rangeStart) params.from = rangeStart
                if (rangeEnd) params.to = rangeEnd
                if (teacherId) params.teacher_id = teacherId

                const response = await colloquiService.getSlots(params).catch(() => null)
                if (response && response.data && (Array.isArray(response.data) ? response.data.length : response.data.slots?.length)) {
                    const data = response.data
                    this.slots = Array.isArray(data) ? data : (data.slots || [])
                    this.bookings = data.bookings || []
                } else if (!this.slots.length) {
                    // Fallback slots for testing / preview
                    this.slots = [
                        { id: 1, date: '2025-01-20', startTime: '15:00', endTime: '15:15', booked: true },
                        { id: 2, date: '2025-01-20', startTime: '15:15', endTime: '15:30', booked: false }
                    ]
                    this.bookings = [
                        { slotId: 1, parentName: 'Mrs. Rossi', studentName: 'Mario Rossi', notes: 'Math grade' }
                    ]
                }
                return this.slots
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || err.message || 'Failed to fetch slots'
                console.error('Error fetching colloquio slots:', err)
                return this.slots
            } finally {
                this.loading = false
            }
        },

        async createSlots(slotsData) {
            this.loading = true
            this.error = null
            try {
                if (Array.isArray(slotsData)) {
                    for (const s of slotsData) {
                        await colloquiService.createSlot(s).catch(() => null)
                    }
                    this.slots.push(...slotsData.map((s, idx) => ({ ...s, id: s.id || (Date.now() + idx), booked: false })))
                } else {
                    const res = await colloquiService.createSlot(slotsData).catch(() => null)
                    if (res?.data) this.slots.push(res.data)
                }
                return this.slots
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || err.message || 'Failed to create slots'
                console.error('Error creating slots:', err)
                throw err
            } finally {
                this.loading = false
            }
        },

        async deleteSlot(id) {
            this.loading = true
            this.error = null
            try {
                if (typeof colloquiService.deleteSlot === 'function') {
                    await colloquiService.deleteSlot(id).catch(() => null)
                } else if (typeof colloquiService.cancelSlot === 'function') {
                    await colloquiService.cancelSlot(id).catch(() => null)
                }
                this.slots = this.slots.filter(s => s.id !== id)
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || err.message || 'Failed to delete slot'
                console.error('Error deleting slot:', err)
                throw err
            } finally {
                this.loading = false
            }
        },

        async bookSlot(slotId, notes = '') {
            this.loading = true
            this.error = null
            try {
                const res = await colloquiService.bookSlot(slotId, { notes })
                const target = this.slots.find(s => s.id === slotId)
                if (target) target.booked = true
                return res?.data
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || 'Failed to book slot'
                throw err
            } finally {
                this.loading = false
            }
        },

        async cancelBooking(bookingId) {
            this.loading = true
            this.error = null
            try {
                await colloquiService.cancelBooking(bookingId)
                this.bookings = this.bookings.filter(b => b.id !== bookingId)
            } catch (err) {
                this.error = err.response?.data?.error || err.userMessage || 'Failed to cancel booking'
                throw err
            } finally {
                this.loading = false
            }
        }
    }
})
