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
  },
  async requestPushPermissionAndRegister() {
    if (!('Notification' in window)) {
      console.warn('Web Push Notifications are not supported by this browser.')
      return null
    }
    const permission = await Notification.requestPermission()
    if (permission === 'granted') {
      const token = 'web_push_' + Math.random().toString(36).substring(2, 15)
      await this.registerPushToken(token, 'web')
      return token
    }
    return null
  }
}
