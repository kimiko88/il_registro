import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import RoomBooking from '@/pages/teacher/RoomBooking.vue'
import roomsService from '@/services/roomsService'
import api from '@/services/api'

vi.mock('@/services/roomsService', () => ({
  default: {
    getBuildings: vi.fn(),
    getRooms: vi.fn(),
    getBookings: vi.fn(),
    createBooking: vi.fn(),
    cancelBooking: vi.fn()
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
      notify: vi.fn(),
      dialog: vi.fn().mockReturnValue({ onOk: (fn) => fn() })
    })
  }
})

describe('Teacher RoomBooking.vue — Room Booking Workflow', () => {
  let wrapper

  const mockBuildings = [
    { id: 'bld-1', name: 'Sede Centrale' },
    { id: 'bld-2', name: 'Succursale' }
  ]

  const mockRooms = [
    { id: 'r-1', name: 'Lab Info 1', building_id: 'bld-1', capacity: 25, equipment: ['PC'] },
    { id: 'r-2', name: 'Palestra', building_id: 'bld-2', capacity: 40, equipment: [] }
  ]

  const mockBookings = [
    {
      id: 'bk-1',
      room_name: 'Lab Info 1',
      building_name: 'Sede Centrale',
      class_name: '3A',
      booking_date: '2026-10-15',
      hour_index: 2,
      status: 'confirmed',
      is_recurring: false
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    roomsService.getBuildings.mockResolvedValue({ data: mockBuildings })
    roomsService.getRooms.mockResolvedValue({ data: mockRooms })
    roomsService.getBookings.mockResolvedValue({ data: mockBookings })
    api.get.mockImplementation((url) => {
      if (url.includes('/classes')) {
        return Promise.resolve({ data: [{ id: 'c-1', name: '3A' }] })
      }
      if (url.includes('/subjects')) {
        return Promise.resolve({ data: [{ id: 's-1', name: 'Informatica' }] })
      }
      return Promise.resolve({ data: [] })
    })

    wrapper = mount(RoomBooking, {
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
          'q-radio': { template: '<div class="q-radio"><slot /></div>' },
          'q-toggle': { template: '<div class="q-toggle"><slot /></div>' },
          'q-table': { template: '<div class="q-table"><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-spinner-dots': true
        }
      }
    })
  })

  it('renders header and loads initial user bookings', async () => {
    await flushPromises()
    expect(wrapper.text()).toContain('Prenotazione Aule e Laboratori')
    expect(wrapper.text()).toContain('Le Tue Prenotazioni')
    expect(roomsService.getBookings).toHaveBeenCalled()
  })

  it('filters available rooms based on selected building in form', async () => {
    await flushPromises()
    expect(wrapper.vm.filteredRooms.length).toBe(2)

    wrapper.vm.filterBuilding = 'bld-1'
    await flushPromises()
    expect(wrapper.vm.filteredRooms.length).toBe(1)
    expect(wrapper.vm.filteredRooms[0].id).toBe('r-1')
  })

  it('submits a booking request calling roomsService.createBooking', async () => {
    await flushPromises()
    roomsService.createBooking.mockResolvedValue({
      data: { id: 'new-bk', room_id: 'r-1', hour_index: 2 }
    })

    wrapper.vm.bookingForm.room_id = 'r-1'
    wrapper.vm.bookingForm.class_id = 'c-1'
    wrapper.vm.bookingForm.hour_index = 2
    wrapper.vm.bookingForm.booking_date = '2026-10-20'

    await wrapper.vm.submitBooking()
    await flushPromises()

    expect(roomsService.createBooking).toHaveBeenCalledWith(
      expect.objectContaining({
        room_id: 'r-1',
        class_id: 'c-1',
        hour_index: 2,
        booking_date: '2026-10-20'
      })
    )
  })
})
