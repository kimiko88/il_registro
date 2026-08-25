import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { ref } from 'vue'
import { useDraftAutosave } from '@/composables/useDraftAutosave'
import { useGlobalAnalytics } from '@/composables/useGlobalAnalytics'
import { useMonitoring } from '@/composables/useMonitoring'
import { useMyAttendance } from '@/composables/useMyAttendance'
import { useChildAttendance } from '@/composables/useChildAttendance'
import { useColloquiBooking } from '@/composables/useColloquiBooking'
import { useUserManagement } from '@/composables/useUserManagement'
import { useAttendanceStore } from '@/stores/attendance'
import { useChildrenStore } from '@/stores/children'
import api from '@/services/api'
import colloquiService from '@/services/colloquiService'
import adminService from '@/services/adminService'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn()
  }
}))

vi.mock('@/services/colloquiService', () => ({
  default: {
    getSlots: vi.fn(),
    bookSlot: vi.fn(),
    cancelSlot: vi.fn()
  }
}))

vi.mock('@/services/adminService', () => ({
  default: {
    getAdmins: vi.fn()
  }
}))

describe('Composables Audit & Hardening Test Suite', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('useDraftAutosave lifecycle & persistence', () => {
    it('executes outside component setup without throwing or warning', async () => {
      const form = ref({ subject: 'Math', title: 'Lesson 1' })
      const { hasDraft, clearDraft, stop } = useDraftAutosave('test-key', form, 50)

      expect(hasDraft.value).toBe(false)
      form.value.title = 'Lesson 2 Updated'

      // Wait for debounce
      await new Promise(r => setTimeout(r, 80))
      expect(hasDraft.value).toBe(true)

      const restored = ref({ subject: '', title: '' })
      const { restoreDraft: restoreInto } = useDraftAutosave('test-key', restored, 50)
      expect(restoreInto()).toBe(true)
      expect(restored.value.title).toBe('Lesson 2 Updated')

      clearDraft()
      expect(hasDraft.value).toBe(false)
      stop()
    })
  })

  describe('useGlobalAnalytics and useMonitoring controls', () => {
    it('useGlobalAnalytics fetches stats on demand without active instance', async () => {
      api.get.mockResolvedValueOnce({
        data: { schools: 5, students: 120, apiCalls: 500 }
      })

      const { stats, fetchStats, loading } = useGlobalAnalytics()
      expect(loading.value).toBe(false)
      await fetchStats()
      expect(stats.value.schools).toBe(5)
      expect(stats.value.students).toBe(120)
    })

    it('useMonitoring provides manual start and stop polling controls', () => {
      const { startPolling, stopPolling } = useMonitoring(false, 10000)
      expect(startPolling).toBeDefined()
      expect(stopPolling).toBeDefined()

      // Start and immediately stop without errors
      startPolling()
      stopPolling()
    })
  })

  describe('useMyAttendance & useChildAttendance status normalization', () => {
    it('correctly aggregates stats regardless of status casing', () => {
      const store = useAttendanceStore()
      store.records = [
        { id: '1', status: 'present', date: '2026-03-01' },
        { id: '2', status: 'Present', date: '2026-03-02' },
        { id: '3', status: 'late', date: '2026-03-03' },
        { id: '4', status: 'Absent', date: '2026-03-04' }
      ]

      const { stats } = useMyAttendance()
      // Total 4: present/late = 3, absent = 1, percentage = 75%
      expect(stats.value.present).toBe(3)
      expect(stats.value.absent).toBe(1)
      expect(stats.value.late).toBe(1)
      expect(stats.value.percentage).toBe(75)
    })

    it('useChildAttendance links with childrenStore and attendanceStore', async () => {
      const childrenStore = useChildrenStore()
      const attendanceStore = useAttendanceStore()
      attendanceStore.fetchMyAttendance = vi.fn().mockResolvedValue()
      attendanceStore.requestJustification = vi.fn().mockResolvedValue()

      childrenStore.children = [{ id: 'child-99', firstName: 'Marco' }]
      childrenStore.selectedChildId = 'child-99'

      const { selectedChild, requestChildJustification } = useChildAttendance()
      expect(selectedChild.value.firstName).toBe('Marco')

      await requestChildJustification('2026-03-01', 'Visita medica')
      expect(attendanceStore.requestJustification).toHaveBeenCalledWith('2026-03-01', 'Visita medica', 'child-99')
    })
  })

  describe('useColloquiBooking & useUserManagement service integrations', () => {
    it('useColloquiBooking calls colloquiService.bookSlot and updates state', async () => {
      colloquiService.bookSlot.mockResolvedValueOnce({ data: { success: true } })

      const { bookSlot } = useColloquiBooking()
      const result = await bookSlot('slot-123', 'Parent notes')
      expect(result).toBe(true)
      expect(colloquiService.bookSlot).toHaveBeenCalledWith({ slot_id: 'slot-123', notes: 'Parent notes' })
    })

    it('useUserManagement passes query parameters to adminService', async () => {
      adminService.getAdmins.mockResolvedValueOnce({
        data: { items: [{ id: 'u1', email: 'user@school.it' }] }
      })

      const { users, fetchUsers } = useUserManagement()
      await fetchUsers({ role: 'teacher', search: 'user' })

      expect(adminService.getAdmins).toHaveBeenCalledWith({ role: 'teacher', search: 'user' })
      expect(users.value).toHaveLength(1)
    })
  })
})
