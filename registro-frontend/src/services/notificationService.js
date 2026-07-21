import api from './api'

export const notificationService = {
  async getNotifications(unreadOnly = false) {
    const response = await api.get('/notifications', { params: { unread_only: unreadOnly } })
    return response
  },
  async markAsRead(id) {
    return api.put(`/notifications/${id}/read`)
  },
  async markAllAsRead() {
    return api.put('/notifications/read-all')
  },
  async registerPushToken(token, platform = 'web') {
    return api.post('/notifications/push-tokens', { device_token: token, platform })
  },
  async unregisterPushToken(token) {
    return api.delete('/notifications/push-tokens', { params: { device_token: token } })
  }
}
