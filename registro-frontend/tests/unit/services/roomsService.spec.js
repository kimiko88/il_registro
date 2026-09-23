import { describe, it, expect, vi, beforeEach } from 'vitest'
import roomsService from '@/services/roomsService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn()
  }
}))

describe('roomsService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Buildings (Plessi)', () => {
    it('getBuildings calls GET /rooms/buildings', () => {
      roomsService.getBuildings()
      expect(api.get).toHaveBeenCalledWith('/rooms/buildings')
    })

    it('createBuilding calls POST /rooms/buildings with payload', () => {
      const payload = { name: 'Plesso Centrale', code: 'CENTRAL' }
      roomsService.createBuilding(payload)
      expect(api.post).toHaveBeenCalledWith('/rooms/buildings', payload)
    })

    it('updateBuilding calls PUT /rooms/buildings/:id', () => {
      const payload = { name: 'Plesso Succursale' }
      roomsService.updateBuilding('bld-1', payload)
      expect(api.put).toHaveBeenCalledWith('/rooms/buildings/bld-1', payload)
    })

    it('deleteBuilding calls DELETE /rooms/buildings/:id', () => {
      roomsService.deleteBuilding('bld-1')
      expect(api.delete).toHaveBeenCalledWith('/rooms/buildings/bld-1')
    })
  })

  describe('Rooms (Aule)', () => {
    it('getRooms calls GET /rooms with query params', () => {
      const params = { building_id: 'b-1', room_type: 'computer_lab' }
      roomsService.getRooms(params)
      expect(api.get).toHaveBeenCalledWith('/rooms', { params })
    })

    it('createRoom calls POST /rooms with payload', () => {
      const payload = { name: 'Laboratorio 1', capacity: 30 }
      roomsService.createRoom(payload)
      expect(api.post).toHaveBeenCalledWith('/rooms', payload)
    })

    it('getRoomAvailability calls GET /rooms/:id/availability with dates', () => {
      roomsService.getRoomAvailability('r-1', '2026-10-01', '2026-10-07')
      expect(api.get).toHaveBeenCalledWith('/rooms/r-1/availability', {
        params: { from: '2026-10-01', to: '2026-10-07' }
      })
    })
  })

  describe('Bookings (Prenotazioni)', () => {
    it('createBooking calls POST /rooms/bookings with payload', () => {
      const payload = { room_id: 'r-1', class_id: 'c-1', booking_date: '2026-10-05', day_of_week: 1, hour_slot: 2 }
      roomsService.createBooking(payload)
      expect(api.post).toHaveBeenCalledWith('/rooms/bookings', payload)
    })

    it('cancelBooking calls DELETE /rooms/bookings/:id', () => {
      roomsService.cancelBooking('bk-1', true)
      expect(api.delete).toHaveBeenCalledWith('/rooms/bookings/bk-1', { params: { cancel_series: true } })
    })

    it('getBookings calls GET /rooms/bookings', () => {
      roomsService.getBookings({ room_id: 'r-1' })
      expect(api.get).toHaveBeenCalledWith('/rooms/bookings', { params: { room_id: 'r-1' } })
    })
  })
})
