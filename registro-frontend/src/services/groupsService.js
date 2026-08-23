import api from './api';

export const groupsService = {
  getGroups(params) {
    return api.get('/groups', { params });
  },
  getGroup(id) {
    return api.get(`/groups/${id}`);
  },
  createGroup(data) {
    return api.post('/groups', data);
  },
  updateGroup(id, data) {
    return api.put(`/groups/${id}`, data);
  },
  deleteGroup(id) {
    return api.delete(`/groups/${id}`);
  },
  addStudents(id, studentIds) {
    return api.post(`/groups/${id}/students`, { student_ids: studentIds });
  },
  removeStudent(id, studentId) {
    return api.delete(`/groups/${id}/students/${studentId}`);
  }
};

export default groupsService;

