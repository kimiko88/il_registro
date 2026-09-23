import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import TimetableConstraints from '@/pages/secretary/TimetableConstraints.vue'
import timetableGenService from '@/services/timetableGenService'
import api from '@/services/api'

vi.mock('@/services/timetableGenService', () => ({
  default: {
    getRoomRequirements: vi.fn(),
    saveRoomRequirement: vi.fn(),
    deleteRoomRequirement: vi.fn(),
    getConstraints: vi.fn(),
    saveConstraint: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      dark: { isActive: false },
      notify: vi.fn()
    })
  }
})

describe('Secretary TimetableConstraints.vue — Timetable Constraints & Room Requirements', () => {
  let wrapper

  const mockReqs = [
    {
      id: 'req-1',
      subject_id: 'sub-info',
      subject_name: 'Informatica',
      required_room_type: 'lab_informatica',
      is_mandatory: true
    }
  ]

  const mockSubjects = [
    { id: 'sub-info', name: 'Informatica' },
    { id: 'sub-chim', name: 'Chimica' }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    timetableGenService.getRoomRequirements.mockResolvedValue({ data: mockReqs })
    timetableGenService.getConstraints.mockResolvedValue({ data: [] })
    api.get.mockImplementation((url) => {
      if (url.includes('/subjects')) {
        return Promise.resolve({ data: mockSubjects })
      }
      return Promise.resolve({ data: [] })
    })

    wrapper = mount(TimetableConstraints, {
      global: {
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-chip': { template: '<span><slot /></span>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-select': { template: '<div class="q-select"><slot /></div>' },
          'q-input': { template: '<div class="q-input"><slot /></div>' },
          'q-toggle': { template: '<div class="q-toggle"><slot /></div>' },
          'q-table': { template: '<div class="q-table"><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-spinner-dots': true
        }
      }
    })
  })

  it('renders header and loads initial room requirements', async () => {
    await flushPromises()
    expect(wrapper.text()).toContain('Vincoli e Requisiti Orario Scolastico')
    expect(wrapper.text()).toContain('Assegnazione Aule Speciali / Laboratori per Materia')
    expect(timetableGenService.getRoomRequirements).toHaveBeenCalled()
    expect(wrapper.vm.roomReqs.length).toBe(1)
  })

  it('saves new room requirement calling timetableGenService.saveRoomRequirement', async () => {
    await flushPromises()
    timetableGenService.saveRoomRequirement.mockResolvedValue({ data: { id: 'req-2' } })

    wrapper.vm.openRoomReqModal()
    wrapper.vm.reqForm.subject_id = 'sub-chim'
    wrapper.vm.reqForm.required_room_type = 'lab_chimica'

    await wrapper.vm.saveReq()
    await flushPromises()

    expect(timetableGenService.saveRoomRequirement).toHaveBeenCalledWith(
      expect.objectContaining({
        subject_id: 'sub-chim',
        required_room_type: 'lab_chimica'
      })
    )
  })
})
