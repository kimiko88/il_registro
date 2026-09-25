import api from './api';

export const timetableGenService = {
  // Automatic Timetable Generation Job
  startGeneration(data) {
    return api.post('/timetable/generate', data);
  },
  getJobStatus(jobId) {
    return api.get(`/timetable/generate/${jobId}`);
  },
  publishSchedule(jobId) {
    return api.post(`/timetable/generate/${jobId}/publish`);
  },

  // Teacher Schedule Preferences (Desiderata)
  getPreferences(params) {
    return api.get('/timetable/preferences', { params });
  },
  savePreferences(data, params) {
    return params ? api.post('/timetable/preferences', data, { params }) : api.post('/timetable/preferences', data);
  },
  getDesiderataWindow() {
    return api.get('/timetable/preferences/window');
  },
  setDesiderataWindow(isOpen) {
    return api.post('/timetable/preferences/window', { is_open: isOpen });
  },
  adjustSchedule(jobId, data) {
    return api.post(`/timetable/generate/${jobId}/adjust`, data);
  },

  // Room Requirements per Subject
  getRoomRequirements() {
    return api.get('/timetable/room-requirements');
  },
  saveRoomRequirement(data) {
    return api.post('/timetable/room-requirements', data);
  },
  deleteRoomRequirement(id) {
    return api.delete(`/timetable/room-requirements/${id}`);
  },

  // Configurable Constraints
  getConstraints() {
    return api.get('/timetable/constraints');
  },
  saveConstraint(data) {
    return api.post('/timetable/constraints', data);
  },
  deleteConstraint(id) {
    return api.delete(`/timetable/constraints/${id}`);
  }
};

export default timetableGenService;
