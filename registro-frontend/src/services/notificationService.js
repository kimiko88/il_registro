import api from './api'

export const notificationService = {
  getNotifications(unreadOnly = false) {
    return api.get('/notifications', { params: { unread_only: unreadOnly } })
  },
  markAsRead(id) {
    return api.put(`/notifications/${id}/read`)
  },
  markAllAsRead() {
    return api.put('/notifications/read-all')
  },
  registerPushToken(token, platform = 'web') {
    return api.post('/notifications/push-tokens', { device_token: token, platform })
  },
  unregisterPushToken(token) {
    return api.delete('/notifications/push-tokens', { params: { device_token: token } })
  },
  async requestPushPermissionAndRegister() {
    if (typeof window === 'undefined' || !('Notification' in window)) {
      console.warn('Web Push Notifications are not supported by this browser.')
      return null
    }
    if (Notification.permission === 'denied') {
      return null
    }
    const permission = await Notification.requestPermission()
    if (permission === 'granted') {
      const uuid = (typeof crypto !== 'undefined' && crypto.randomUUID)
        ? crypto.randomUUID()
        : Array.from(crypto.getRandomValues(new Uint8Array(16)), b => b.toString(16).padStart(2, '0')).join('')
      const token = 'web_push_' + uuid
      await this.registerPushToken(token, 'web')
      return token
    }
    return null
  }
}

export default notificationService

