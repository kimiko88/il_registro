import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'
import { useGradesStore } from './grades'
import { useAttendanceStore } from './attendance'
import { useCommunicationsStore } from './communications'
import { useScrutinyStore } from './scrutiny'
import { Notify } from 'quasar'

const escapeHtml = (str) => {
    if (!str) return ''
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;')
}

export const useWebSocketStore = defineStore('websocket', () => {
    const socket = ref(null)
    const isConnected = ref(false)
    const reconnectTimer = ref(null)
    const reconnectAttempts = ref(0)
    const hasFailedPermanently = ref(false)
    const lastError = ref(null)
    const authStore = useAuthStore()

    const heartbeatTimer = ref(null)

    function startHeartbeat() {
        stopHeartbeat()
        heartbeatTimer.value = setInterval(() => {
            if (socket.value && socket.value.readyState === WebSocket.OPEN) {
                try {
                    socket.value.send(JSON.stringify({ type: 'PING' }))
                } catch (e) {
                    console.error('WebSocket: Heartbeat send error, closing socket', e)
                    socket.value.close()
                }
            } else if (socket.value && socket.value.readyState !== WebSocket.CONNECTING) {
                console.warn('WebSocket: Half-open / silent dead connection detected, closing socket')
                socket.value.close()
            }
        }, 25000)
    }

    function stopHeartbeat() {
        if (heartbeatTimer.value) {
            clearInterval(heartbeatTimer.value)
            heartbeatTimer.value = null
        }
    }

    function connect() {
        if (socket.value && (socket.value.readyState === WebSocket.OPEN || socket.value.readyState === WebSocket.CONNECTING)) {
            return
        }

        if (!authStore.isAuthenticated) {
            console.warn('WebSocket: Token expired or missing, skipping connection')
            return
        }

        const rawUrl = import.meta.env.VITE_API_URL
        let baseUrl = rawUrl || `${window.location.protocol}//${window.location.host}/api/v1`
        if (rawUrl && !rawUrl.endsWith('/api/v1') && !rawUrl.endsWith('/api/v1/')) {
            baseUrl = rawUrl.endsWith('/') ? `${rawUrl}api/v1` : `${rawUrl}/api/v1`
        }

        let wsUrl = ''
        if (baseUrl.startsWith('http://') || baseUrl.startsWith('https://')) {
            wsUrl = `${baseUrl.replace(/^http/, 'ws')}/ws`
        } else if (baseUrl.startsWith('/')) {
            const host = window.location.host
            const wsScheme = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
            wsUrl = `${wsScheme}//${host}${baseUrl}/ws`
        } else {
            const wsScheme = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
            wsUrl = `${wsScheme}//${baseUrl}/ws`
        }

        const token = authStore.token
        try {
            socket.value = new WebSocket(wsUrl, ['access_token', token])
        } catch (e) {
            console.error('WebSocket connection error:', e)
            attemptReconnect()
            return
        }

        socket.value.onopen = () => {
            console.log('WebSocket: Connected')
            isConnected.value = true
            reconnectAttempts.value = 0
            if (reconnectTimer.value) {
                clearTimeout(reconnectTimer.value)
                reconnectTimer.value = null
            }
            startHeartbeat()
        }

        socket.value.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data)
                if (message && message.type === 'PONG') return
                handleMessage(message)
            } catch (e) {
                console.error('WebSocket: Failed to parse message', e)
            }
        }

        socket.value.onclose = (event) => {
            console.log('WebSocket: Closed', event)
            isConnected.value = false
            socket.value = null
            stopHeartbeat()
            attemptReconnect()
        }

        socket.value.onerror = (error) => {
            console.error('WebSocket: Error', error)
            stopHeartbeat()
            if (socket.value) {
                socket.value.close()
            }
        }
    }

    function disconnect() {
        stopHeartbeat()
        if (socket.value) {
            socket.value.close()
            socket.value = null
        }
        isConnected.value = false
        reconnectAttempts.value = 0
        hasFailedPermanently.value = false
        lastError.value = null
        if (reconnectTimer.value) {
            clearTimeout(reconnectTimer.value)
            reconnectTimer.value = null
        }
    }

    function attemptReconnect() {
        if (!authStore.isAuthenticated) {
            disconnect()
            return
        }

        if (reconnectAttempts.value >= 30) {
            console.warn('WebSocket: Reached max reconnect attempts (30), stopping automatic reconnection')
            hasFailedPermanently.value = true
            lastError.value = 'Connessione WebSocket non disponibile dopo tentativi ripetuti.'
            disconnect()
            return
        }

        if (reconnectTimer.value) return

        // Exponential backoff with random jitter (1s, 2s, 4s, 8s... up to max 30s)
        const baseDelay = 1000 * Math.pow(2, Math.min(reconnectAttempts.value, 5))
        const maxDelay = 30000
        const jitter = Math.random() * 1000
        const delay = Math.min(baseDelay + jitter, maxDelay)

        reconnectAttempts.value++
        console.log(`WebSocket: Attempting reconnect in ${(delay / 1000).toFixed(1)}s (attempt ${reconnectAttempts.value})`)

        reconnectTimer.value = setTimeout(() => {
            reconnectTimer.value = null
            if (authStore.isAuthenticated) {
                connect()
            } else {
                disconnect()
            }
        }, delay)
    }

    function handleMessage(message) {
        if (!message || !message.type) return

        const payload = message.payload || {}

        switch (message.type) {
            case 'GRADE_ADDED':
            case 'GRADE_UPDATED':
                try {
                    const gradesStore = useGradesStore()
                    if (payload.class_id) gradesStore.fetchGrades(payload.class_id, payload.subject_id, true)
                } catch (err) {
                    console.debug('Failed to refresh grades store:', err)
                }
                Notify.create({
                    message: `Aggiornamento voto: ${escapeHtml(payload.grade_value || '')} (${escapeHtml(payload.subject_name || 'Materia')})`,
                    color: 'info',
                    icon: 'school',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            case 'ATTENDANCE_LATE':
            case 'ATTENDANCE_ABSENT':
            case 'ATTENDANCE_PRESENT':
                try {
                    const attendanceStore = useAttendanceStore()
                    attendanceStore.fetchMyAttendance()
                } catch (err) {
                    console.debug('Failed to refresh attendance store:', err)
                }
                Notify.create({
                    message: `Aggiornamento presenze: ${escapeHtml(payload.status || 'Presenza registrata')}`,
                    color: 'warning',
                    icon: 'warning',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            case 'JUSTIFICATION_APPROVED':
                try {
                    const attendanceStore = useAttendanceStore()
                    attendanceStore.fetchMyAttendance()
                } catch (err) {
                    console.debug('Failed to refresh attendance store:', err)
                }
                Notify.create({
                    message: `Giustifica approvata: ${escapeHtml(payload.reason || '')}`,
                    color: 'positive',
                    icon: 'check_circle',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            case 'JUSTIFICATION_REJECTED':
                try {
                    const attendanceStore = useAttendanceStore()
                    attendanceStore.fetchMyAttendance()
                } catch (err) {
                    console.debug('Failed to refresh attendance store:', err)
                }
                Notify.create({
                    message: `Giustifica non approvata: ${escapeHtml(payload.reason || '')}`,
                    color: 'negative',
                    icon: 'cancel',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            case 'NEW_COMMUNICATION':
            case 'COMMUNICATION_PUBLISHED':
                try {
                    const commsStore = useCommunicationsStore()
                    commsStore.fetchCommunications()
                } catch (err) {
                    console.debug('Failed to refresh communications store:', err)
                }
                Notify.create({
                    message: `Nuova comunicazione: ${escapeHtml(payload.title || 'Circolare scolastica')}`,
                    color: 'primary',
                    icon: 'mail',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            case 'NOTE_ADDED':
                Notify.create({
                    message: `Nuova nota disciplinare registrata: ${escapeHtml(payload.title)}`,
                    color: 'negative',
                    icon: 'report_problem',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            case 'SCRUTINY_PUBLISHED':
                try {
                    const scrutinyStore = useScrutinyStore()
                    scrutinyStore.fetchOverview()
                } catch (err) {
                    console.debug('Failed to refresh scrutiny store:', err)
                }
                Notify.create({
                    message: `Esito scrutinio pubblicato per ${escapeHtml(payload.student_name || 'lo studente')}`,
                    color: 'positive',
                    icon: 'assignment_turned_in',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            case 'GOAL_UPDATED':
                Notify.create({
                    message: `Obiettivo aggiornato: ${escapeHtml(payload.title)}`,
                    color: 'secondary',
                    icon: 'star',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            case 'SLOT_BOOKED':
            case 'SLOT_CANCELLED':
                Notify.create({
                    message: `Aggiornamento colloquio: ${escapeHtml(payload.message || message.type)}`,
                    color: 'accent',
                    icon: 'event',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            default:
                if (payload.title || payload.body) {
                    Notify.create({
                        message: payload.title ? `${escapeHtml(payload.title)}: ${escapeHtml(payload.body)}` : escapeHtml(payload.body),
                        color: 'info',
                        icon: 'notifications',
                        position: 'top-right',
                        attrs: { role: 'alert' }
                    })
                }
                break
        }
    }

    return {
        connect,
        disconnect,
        isConnected,
        reconnectAttempts,
        hasFailedPermanently,
        lastError
    }
})
