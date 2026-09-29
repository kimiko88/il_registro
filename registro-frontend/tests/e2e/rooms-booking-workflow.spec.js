import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import Rooms from '@/pages/secretary/Rooms.vue'
import RoomBooking from '@/pages/teacher/RoomBooking.vue'
import roomsService from '@/services/roomsService'
import api from '@/services/api'

vi.mock('@/services/roomsService', () => ({
  default: {
    getBuildings: vi.fn(),
    getRooms: vi.fn(),
    getBookings: vi.fn(),
    createBuilding: vi.fn(),
    createRoom: vi.fn(),
    createBooking: vi.fn(),
    cancelBooking: vi.fn(),
    getRoomAvailability: vi.fn()
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
      loading: { show: vi.fn(), hide: vi.fn() },
      notify: vi.fn(),
      dialog: vi.fn().mockReturnValue({
        onOk: vi.fn(cb => { if (cb) cb('series'); return { onCancel: vi.fn(), onDismiss: vi.fn() } })
      })
    })
  }
})

describe('E2E Workflow: Multi-Building Bookable Rooms & Teacher Reservations', () => {
  let piniaSecretary
  let piniaTeacher

  const createdBuilding = { id: 'bld-fermi', name: 'Plesso Fermi', address: 'Via Roma 10' }
  const createdLabRoom = {
    id: 'room-lab-1',
    name: 'Laboratorio Multimediale',
    room_type: 'lab_informatica',
    building_id: 'bld-fermi',
    building_name: 'Plesso Fermi',
    capacity: 28,
    equipment: ['LIM', 'PC'],
    is_active: true
  }

  beforeEach(() => {
    vi.clearAllMocks()

    piniaSecretary = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { user: { role: 'secretary', id: 'sec-1', school_id: 'school-1' }, isAuthenticated: true }
      }
    })

    piniaTeacher = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { user: { role: 'teacher', id: 'teacher-1', school_id: 'school-1' }, isAuthenticated: true }
      }
    })

    roomsService.getBuildings.mockResolvedValue({ data: [createdBuilding] })
    roomsService.getRooms.mockResolvedValue({ data: [createdLabRoom] })
    roomsService.getBookings.mockResolvedValue({ data: [] })
    api.get.mockImplementation((url) => {
      if (url.includes('/classes')) return Promise.resolve({ data: [{ id: 'c-3A', name: '3A' }] })
      if (url.includes('/subjects')) return Promise.resolve({ data: [{ id: 's-inf', name: 'Informatica' }] })
      return Promise.resolve({ data: [] })
    })
  })

  it('Step 1: Secretary creates school building (plesso) and registers computer lab', async () => {
    const wrapper = mount(Rooms, {
      global: {
        plugins: [piniaSecretary],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-tabs': { template: '<div><slot /></div>' },
          'q-tab': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-chip': { template: '<span><slot /></span>' },
          'q-select': { template: '<div><slot /></div>' },
          'q-input': { template: '<div><slot /></div>' },
          'q-toggle': { template: '<div><slot /></div>' },
          'q-table': { template: '<div><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' }
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Aule Prenotabili & Plessi')

    // Secretary triggers room creation
    roomsService.createRoom.mockResolvedValue({ data: createdLabRoom })
    wrapper.vm.roomForm = {
      name: 'Laboratorio Multimediale',
      building_id: 'bld-fermi',
      room_type: 'lab_informatica',
      capacity: 28,
      equipment: ['LIM', 'PC'],
      requires_booking: true,
      is_active: true
    }

    await wrapper.vm.saveRoom()
    await flushPromises()

    expect(roomsService.createRoom).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'Laboratorio Multimediale',
        room_type: 'lab_informatica',
        capacity: 28
      })
    )
  })

  it('Step 2: Teacher navigates to booking page, selects building and books recurring weekly slots', async () => {
    const wrapper = mount(RoomBooking, {
      global: {
        plugins: [piniaTeacher],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-chip': { template: '<span><slot /></span>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-select': { template: '<div><slot /></div>' },
          'q-input': { template: '<div><slot /></div>' },
          'q-radio': { template: '<div><slot /></div>' },
          'q-toggle': { template: '<div><slot /></div>' },
          'q-table': { template: '<div><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-spinner-dots': true
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Prenotazione Aule e Laboratori')

    // Teacher fills booking modal for 3A class, weekly recurrence
    wrapper.vm.bookingForm = {
      room_id: 'room-lab-1',
      class_id: 'c-3A',
      subject_id: 's-inf',
      booking_date: '2026-10-15',
      hour_index: 2,
      notes: 'Laboratorio coding',
      is_recurring: true,
      recurrence_pattern: 'weekly',
      recurring_until: '2026-11-15'
    }

    roomsService.createBooking.mockResolvedValue({
      data: {
        message: 'prenotazioni ricorrenti create con successo',
        total_booked: 4,
        occurrences: [
          { id: 'occ-1', booking_date: '2026-10-15', hour_index: 2, status: 'confirmed' },
          { id: 'occ-2', booking_date: '2026-10-22', hour_index: 2, status: 'confirmed' },
          { id: 'occ-3', booking_date: '2026-10-29', hour_index: 2, status: 'confirmed' },
          { id: 'occ-4', booking_date: '2026-11-05', hour_index: 2, status: 'confirmed' }
        ]
      }
    })

    await wrapper.vm.submitBooking()
    await flushPromises()

    expect(roomsService.createBooking).toHaveBeenCalledWith(
      expect.objectContaining({
        room_id: 'room-lab-1',
        class_id: 'c-3A',
        is_recurring: true,
        recurrence_pattern: 'weekly'
      })
    )
  })

  it('Step 3: Teacher cancels recurring booking series with cancel_series=true', async () => {
    const wrapper = mount(RoomBooking, {
      global: {
        plugins: [piniaTeacher],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-chip': { template: '<span><slot /></span>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-select': { template: '<div><slot /></div>' },
          'q-input': { template: '<div><slot /></div>' },
          'q-radio': { template: '<div><slot /></div>' },
          'q-toggle': { template: '<div><slot /></div>' },
          'q-table': { template: '<div><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-spinner-dots': true
        }
      }
    })

    await flushPromises()
    roomsService.cancelBooking.mockResolvedValue({ data: { message: 'prenotazione cancellata' } })

    const sampleRecurringBooking = {
      id: 'occ-1',
      room_name: 'Laboratorio Multimediale',
      booking_date: '2026-10-15',
      hour_index: 2,
      is_recurring: true
    }

    await wrapper.vm.confirmCancel(sampleRecurringBooking)
    await flushPromises()

    expect(roomsService.cancelBooking).toHaveBeenCalledWith('occ-1', true)
  })
})
