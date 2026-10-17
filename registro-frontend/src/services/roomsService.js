import api from './api';

export const roomsService = {
  // Buildings (Plessi)
  getBuildings() {
    return api.get('/rooms/buildings');
  },
  getBuilding(id) {
    return api.get(`/rooms/buildings/${id}`);
  },
  createBuilding(data) {
    return api.post('/rooms/buildings', data);
  },
  updateBuilding(id, data) {
    return api.put(`/rooms/buildings/${id}`, data);
  },
  deleteBuilding(id) {
    return api.delete(`/rooms/buildings/${id}`);
  },

  // Rooms (Aule)
  getRooms(params) {
    return api.get('/rooms', { params });
  },
  getRoom(id) {
    return api.get(`/rooms/${id}`);
  },
  getRoomAvailability(id, from, to) {
    return api.get(`/rooms/${id}/availability`, {
      params: { from, to }
    });
  },
  createRoom(data) {
    return api.post('/rooms', data);
  },
  updateRoom(id, data) {
    return api.put(`/rooms/${id}`, data);
  },
  deleteRoom(id) {
    return api.delete(`/rooms/${id}`);
  },

  // Bookings (Prenotazioni)
  getBookings(params) {
    return api.get('/rooms/bookings', { params });
  },
  createBooking(data) {
    return api.post('/rooms/bookings', data);
  },
  cancelBooking(id, cancelSeries = false) {
    return api.delete(`/rooms/bookings/${id}`, {
      params: { cancel_series: cancelSeries }
    });
  }
};

export default roomsService;
