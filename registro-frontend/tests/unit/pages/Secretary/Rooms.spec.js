import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import Rooms from '@/pages/secretary/Rooms.vue'
import roomsService from '@/services/roomsService'

vi.mock('@/services/roomsService', () => ({
  default: {
    getBuildings: vi.fn(),
    getRooms: vi.fn(),
    getBookings: vi.fn(),
    createBuilding: vi.fn(),
    createRoom: vi.fn(),
    cancelBooking: vi.fn()
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

describe('Secretary Rooms.vue — Bookable Rooms & Buildings Management', () => {
  let wrapper

  const mockBuildings = [
    { id: 'bld-1', name: 'Sede Centrale', address: 'Via Roma 1' },
    { id: 'bld-2', name: 'Succursale Est', address: 'Via Milano 20' }
  ]

  const mockRooms = [
    {
      id: 'r-1',
      name: 'Lab Informatica 1',
      room_type: 'lab_informatica',
      building_id: 'bld-1',
      building_name: 'Sede Centrale',
      capacity: 25,
      equipment: ['PC', 'LIM'],
      is_active: true
    },
    {
      id: 'r-2',
      name: 'Palestra Principale',
      room_type: 'palestra',
      building_id: 'bld-1',
      capacity: 50,
      equipment: ['Spogliatoi', 'Canestri'],
      is_active: true
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    roomsService.getBuildings.mockResolvedValue({ data: mockBuildings })
    roomsService.getRooms.mockResolvedValue({ data: mockRooms })
    roomsService.getBookings.mockResolvedValue({ data: [] })

    wrapper = mount(Rooms, {
      global: {
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
          'q-tab': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-icon': true,
          'q-chip': { template: '<span class="q-chip"><slot /></span>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-select': true,
          'q-input': true,
          'q-toggle': true,
          'q-table': { template: '<div class="q-table"><slot /></div>' },
          'q-dialog': { template: '<div><slot /></div>' }
        }
      }
    })
  })

  it('renders header and calls getBuildings and getRooms on mount', async () => {
    await flushPromises()
    expect(wrapper.text()).toContain('Aule Prenotabili & Plessi')
    expect(roomsService.getBuildings).toHaveBeenCalled()
    expect(roomsService.getRooms).toHaveBeenCalled()
  })

  it('displays the list of loaded rooms and their equipment tags', async () => {
    await flushPromises()
    expect(wrapper.text()).toContain('Lab Informatica 1')
    expect(wrapper.text()).toContain('Palestra Principale')
    expect(wrapper.text()).toContain('PC')
    expect(wrapper.text()).toContain('LIM')
  })

  it('switches to Plessi tab and displays buildings list', async () => {
    await flushPromises()
    wrapper.vm.activeTab = 'buildings'
    await flushPromises()
    expect(wrapper.text()).toContain('Sede Centrale')
    expect(wrapper.text()).toContain('Succursale Est')
  })
})
