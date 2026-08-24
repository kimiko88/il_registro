import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'
import { useGradesStore } from './grades'
import { useAttendanceStore } from './attendance'
import { useCommunicationsStore } from './communications'
import { useScrutinyStore } from './scrutiny'
import { Notify } from 'quasar'
import { i18n } from '@/i18n'

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
    const heartbeatTimer = ref(null)
    const isReconnecting = ref(false)
    const authStore = useAuthStore()

    const debounceTimers = {}
    function debouncedFetch(key, fetchFn, delayMs = 300) {
        if (debounceTimers[key]) clearTimeout(debounceTimers[key])
        debounceTimers[key] = setTimeout(() => {
            fetchFn()
            delete debounceTimers[key]
        }, delayMs)
    }

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

    async function connect() {
        if (reconnectTimer.value) {
            clearTimeout(reconnectTimer.value)
            reconnectTimer.value = null
        }

        if (socket.value && (socket.value.readyState === WebSocket.OPEN || socket.value.readyState === WebSocket.CONNECTING)) {
            return
        }

        if (!authStore.isAuthenticated) {
            console.warn('WebSocket: Token expired or missing, skipping connection')
            return
        }

        // 1. Acquire single-use WS ticket via REST API (Bearer token in Authorization header)
        const rawUrl = import.meta.env.VITE_API_URL
        let baseUrl = rawUrl ? rawUrl.replace(/\/+$/, '') : `${window.location.protocol}//${window.location.host}/api/v1`

        let ticket = null
        try {
            const controller = new AbortController()
            const timeoutId = setTimeout(() => controller.abort(), 5000)
            const res = await fetch(`${baseUrl}/auth/ws-ticket`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${authStore.token}`,
                    'Content-Type': 'application/json'
                },
                signal: controller.signal
            })
            clearTimeout(timeoutId)
            if (!res.ok) throw new Error(`ws-ticket status ${res.status}`)
            const data = await res.json()
            ticket = data.ticket
        } catch (e) {
            console.error('WebSocket: Failed to acquire ws ticket', e)
            attemptReconnect()
            return
        }

        if (!ticket) {
            attemptReconnect()
            return
        }

        // 2. Build WebSocket URL with opaque ticket parameter
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

        wsUrl += `?ticket=${encodeURIComponent(ticket)}`

        try {
            socket.value = new WebSocket(wsUrl)
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
            const lines = (event.data || '').split('\n').filter(line => line.trim().length > 0)
            for (const line of lines) {
                try {
                    const message = JSON.parse(line)
                    if (message && (message.type === 'PONG' || message.type === 'AUTH_ACK')) continue
                    handleMessage(message)
                } catch (e) {
                    console.error('WebSocket: Failed to parse message line', e)
                }
            }
        }

        socket.value.onclose = (event) => {
            console.log('WebSocket: Closed', event)
            isConnected.value = false
            socket.value = null
            stopHeartbeat()
            if (!event.wasClean && event.code !== 1000 && !hasFailedPermanently.value) {
                attemptReconnect()
            }
        }

        socket.value.onerror = (error) => {
            console.error('WebSocket: Error', error)
            stopHeartbeat()
            if (socket.value) {
                socket.value.close()
            }
        }
    }

    function disconnect(resetPermanentFlag = false) {
        stopHeartbeat()
        if (socket.value) {
            socket.value.close()
            socket.value = null
        }
        isConnected.value = false
        isReconnecting.value = false
        reconnectAttempts.value = 0
        if (resetPermanentFlag) {
            hasFailedPermanently.value = false
            lastError.value = null
        }
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

        if (isReconnecting.value || reconnectTimer.value) return

        reconnectAttempts.value++

        if (reconnectAttempts.value > 8) {
            console.warn('WebSocket: Reached max reconnect attempts (8), stopping automatic reconnection')
            hasFailedPermanently.value = true
            const t = i18n?.global?.t
            lastError.value = t ? t('notifications.wsConnectionFailed') : 'Connessione WebSocket non disponibile dopo tentativi ripetuti.'
            disconnect()
            return
        }

        isReconnecting.value = true

        // Exponential backoff with random jitter (1s, 2s, 4s, 8s... up to max 30s)
        const baseDelay = 1000 * Math.pow(2, Math.min(reconnectAttempts.value - 1, 5))
        const maxDelay = 30000
        const jitter = Math.random() * 1000
        const delay = Math.min(baseDelay + jitter, maxDelay)

        reconnectTimer.value = setTimeout(async () => {
            reconnectTimer.value = null
            try {
                if (authStore.isAuthenticated && authStore.token) {
                    console.log(`WebSocket: Executing reconnect attempt ${reconnectAttempts.value}`)
                    await connect()
                } else {
                    disconnect()
                }
            } finally {
                isReconnecting.value = false
            }
        }, delay)
    }

    function handleMessage(message) {
        if (!message || !message.type) return

        const payload = message.payload || {}

        const t = i18n?.global?.t

        switch (message.type) {
            case 'GRADE_ADDED':
            case 'GRADE_UPDATED':
            case 'GRADE_DELETED': {
                try {
                    const gradesStore = useGradesStore()
                    const targetClassId = gradesStore.currentClassId || payload.class_id
                    if (targetClassId) {
                        debouncedFetch(`grades_${targetClassId}`, () => {
                            gradesStore.fetchGrades(targetClassId, payload.subject_id, true)
                        })
                    }
                } catch (err) {
                    console.debug('Failed to refresh grades store:', err)
                }
                const subjectFallback = t ? t('gradesPage.student') : 'Materia'
                Notify.create({
                    message: message.type === 'GRADE_DELETED'
                        ? (t ? t('notifications.wsGradeDeleted', { subject: escapeHtml(payload.subject_name || subjectFallback) }) : `Voto eliminato per ${escapeHtml(payload.subject_name || subjectFallback)}`)
                        : (t ? t('notifications.wsGradeUpdated', { value: escapeHtml(payload.grade_value || ''), subject: escapeHtml(payload.subject_name || subjectFallback) }) : `Aggiornamento voto: ${escapeHtml(payload.grade_value || '')} (${escapeHtml(payload.subject_name || subjectFallback)})`),
                    color: 'info',
                    icon: 'school',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            case 'ATTENDANCE_LATE':
            case 'ATTENDANCE_ABSENT':
            case 'ATTENDANCE_PRESENT': {
                try {
                    const attendanceStore = useAttendanceStore()
                    debouncedFetch('my_attendance', () => {
                        attendanceStore.fetchMyAttendance()
                    })
                } catch (err) {
                    console.debug('Failed to refresh attendance store:', err)
                }
                Notify.create({
                    message: t ? t('notifications.wsAttendanceUpdated', { status: escapeHtml(payload.status || '') }) : `Aggiornamento presenze: ${escapeHtml(payload.status || 'Presenza registrata')}`,
                    color: 'warning',
                    icon: 'warning',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            case 'JUSTIFICATION_APPROVED': {
                try {
                    const attendanceStore = useAttendanceStore()
                    attendanceStore.fetchMyAttendance()
                } catch (err) {
                    console.debug('Failed to refresh attendance store:', err)
                }
                Notify.create({
                    message: t ? t('notifications.wsJustificationApproved', { reason: escapeHtml(payload.reason || '') }) : `Giustifica approvata: ${escapeHtml(payload.reason || '')}`,
                    color: 'positive',
                    icon: 'check_circle',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            case 'JUSTIFICATION_REJECTED': {
                try {
                    const attendanceStore = useAttendanceStore()
                    attendanceStore.fetchMyAttendance()
                } catch (err) {
                    console.debug('Failed to refresh attendance store:', err)
                }
                Notify.create({
                    message: t ? t('notifications.wsJustificationRejected', { reason: escapeHtml(payload.reason || '') }) : `Giustifica non approvata: ${escapeHtml(payload.reason || '')}`,
                    color: 'negative',
                    icon: 'cancel',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            case 'NEW_COMMUNICATION':
            case 'COMMUNICATION_PUBLISHED': {
                try {
                    const commsStore = useCommunicationsStore()
                    commsStore.fetchCommunications()
                } catch (err) {
                    console.debug('Failed to refresh communications store:', err)
                }
                Notify.create({
                    message: t ? t('notifications.wsNewCommunication', { title: escapeHtml(payload.title || '') }) : `Nuova comunicazione: ${escapeHtml(payload.title || 'Circolare scolastica')}`,
                    color: 'primary',
                    icon: 'mail',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            case 'NOTE_ADDED': {
                Notify.create({
                    message: t ? t('notifications.wsNoteAdded', { title: escapeHtml(payload.title || '') }) : `Nuova nota disciplinare registrata: ${escapeHtml(payload.title || '')}`,
                    color: 'negative',
                    icon: 'report_problem',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            case 'SCRUTINY_PUBLISHED': {
                try {
                    const scrutinyStore = useScrutinyStore()
                    scrutinyStore.fetchOverview()
                } catch (err) {
                    console.debug('Failed to refresh scrutiny store:', err)
                }
                Notify.create({
                    message: t ? t('notifications.wsScrutinyPublished', { student: escapeHtml(payload.student_name || '') }) : `Esito scrutinio pubblicato per ${escapeHtml(payload.student_name || 'lo studente')}`,
                    color: 'positive',
                    icon: 'assignment_turned_in',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            case 'GOAL_UPDATED': {
                Notify.create({
                    message: t ? t('notifications.wsGoalUpdated', { title: escapeHtml(payload.title || '') }) : `Obiettivo aggiornato: ${escapeHtml(payload.title)}`,
                    color: 'secondary',
                    icon: 'star',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            case 'SLOT_BOOKED':
            case 'SLOT_CANCELLED': {
                Notify.create({
                    message: t ? t('notifications.wsSlotUpdated', { msg: escapeHtml(payload.message || '') }) : `Aggiornamento colloquio: ${escapeHtml(payload.message || message.type)}`,
                    color: 'accent',
                    icon: 'event',
                    position: 'top-right',
                    attrs: { role: 'alert' }
                })
                break
            }
            default:
                if ((message.type === 'NOTIFICATION' || message.type === 'SYSTEM_ALERT') && (payload.title || payload.body)) {
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
