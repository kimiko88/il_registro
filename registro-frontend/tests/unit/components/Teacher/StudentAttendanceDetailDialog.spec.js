import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import StudentAttendanceDetailDialog from '@/components/Teacher/StudentAttendanceDetailDialog.vue'
import api from '@/services/api'
import attendanceService from '@/services/attendanceService'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

vi.mock('@/services/attendanceService', () => ({
  default: {
    getStudentSummary: vi.fn()
  }
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key, params) => {
      if (params?.hour) return `Ora ${params.hour}: ${params.status || ''}`
      return key
    }
  })
}))

describe('StudentAttendanceDetailDialog.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders correctly when modelValue is true and fetches student info and summary', async () => {
    const student = { id: 's-123', first_name: 'Mario', last_name: 'Rossi' }
    api.get.mockResolvedValueOnce({
      data: {
        id: 's-123',
        first_name: 'Mario',
        last_name: 'Rossi',
        fiscal_code: 'RSSMRA80A01H501U',
        class_name: '3A',
        email: 'mario.rossi@scuola.it',
        phone_number: '+39 333 1234567',
        date_of_birth: '2008-05-14'
      }
    })
    attendanceService.getStudentSummary.mockResolvedValueOnce({
      total_absences: 4,
      total_lates: 2,
      total_early_exits: 1,
      justified_count: 5,
      absence_rate: 12.5,
      risk_level: 'medium'
    })

    const wrapper = mount(StudentAttendanceDetailDialog, {
      props: {
        modelValue: true,
        student,
        allTodayAttendance: [
          { student_id: 's-123', hour: 1, status: 'Present' },
          { student_id: 's-123', hour: 2, status: 'Late', entry_time: '09:15' }
        ]
      },
      global: {
        stubs: {
          'q-dialog': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-bar': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-spinner': true,
          'q-icon': true,
          'q-space': true,
          'q-btn': true,
          'q-list': { template: '<div><slot /></div>' },
          'q-item': { template: '<div><slot /></div>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-linear-progress': true,
          'q-chip': { template: '<div><slot /></div>' },
          'q-badge': { template: '<div><slot /></div>' },
          'q-tooltip': true
        }
      }
    })

    await vi.waitFor(() => {
      expect(api.get).toHaveBeenCalledWith('/users/s-123')
      expect(attendanceService.getStudentSummary).toHaveBeenCalledWith('s-123')
    })
    await flushPromises()

    expect(wrapper.text()).toContain('Mario')
    expect(wrapper.text()).toContain('Rossi')
    expect(wrapper.text()).toContain('RSSMRA80A01H501U')
    expect(wrapper.text()).toContain('14/05/2008')
    expect(wrapper.text()).toContain('12.5%')
  })

  it('does not fetch when modelValue is false', () => {
    mount(StudentAttendanceDetailDialog, {
      props: {
        modelValue: false,
        student: { id: 's-123', first_name: 'Mario', last_name: 'Rossi' }
      },
      global: {
        stubs: {
          'q-dialog': true
        }
      }
    })

    expect(api.get).not.toHaveBeenCalled()
    expect(attendanceService.getStudentSummary).not.toHaveBeenCalled()
  })

  it('emits update:modelValue when dialog is closed', async () => {
    const wrapper = mount(StudentAttendanceDetailDialog, {
      props: {
        modelValue: true,
        student: { id: 's-123' }
      },
      global: {
        stubs: {
          'q-dialog': {
            template: '<div @close="$emit(\'update:modelValue\', false)"><slot /></div>'
          },
          'q-card': true,
          'q-bar': true,
          'q-card-section': true
        }
      }
    })

    wrapper.vm.$emit('update:modelValue', false)
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
  })
})
