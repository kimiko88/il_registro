import api from './api';

export default {
  getSettings() {
    return api.get('/school-settings');
  },
  updateSettings(data) {
    return api.put('/school-settings', data);
  }
};
