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
  },

  // Curriculum Plans & Class Daily Limits
  getAcademicYears() {
    return api.get('/timetable/academic-years');
  },
  getClassesCurriculumPlans(params) {
    return api.get('/timetable/classes-plans', { params });
  },
  getClassCurriculumPlan(classId) {
    return api.get(`/timetable/classes/${classId}/plan`);
  },
  saveClassCurriculumPlan(classId, data) {
    return api.put(`/timetable/classes/${classId}/plan`, data);
  },
  inheritClassCurriculumPlan(classId, data) {
    return api.post(`/timetable/classes/${classId}/inherit`, data);
  },
  inheritAllClassesCurriculumPlans(data) {
    return api.post('/timetable/inherit-all-plans', data);
  },

  // Teacher Quick Preferences (Tabular representation)
  getTeachersQuickPreferences(params) {
    return api.get('/timetable/teachers-quick-preferences', { params });
  },
  saveTeachersQuickPreferences(data) {
    return api.post('/timetable/teachers-quick-preferences', data);
  },
  saveTeacherQuickPreference(teacherId, data, params) {
    return api.put(`/timetable/teachers-quick-preferences/${teacherId}`, data, { params });
  }
};

export default timetableGenService;
