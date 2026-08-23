import api from './api';

export const schoolSettingsService = {
  getSettings() {
    return api.get('/school-settings');
  },
  updateSettings(data) {
    return api.put('/school-settings', data);
  }
};

export default schoolSettingsService;

