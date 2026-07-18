import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'
import { Notify } from 'quasar'

export const useWebSocketStore = defineStore('websocket', () => {
    const socket = ref(null)
    const isConnected = ref(false)
    const reconnectInterval = ref(null)
    const authStore = useAuthStore()

    function connect() {
        if (socket.value && (socket.value.readyState === WebSocket.OPEN || socket.value.readyState === WebSocket.CONNECTING)) {
            return
        }

        const token = authStore.token
        if (!token) {
            console.warn('WebSocket: No token available, skipping connection')
            return
        }

        // Calculate base API URL safely
        const rawUrl = import.meta.env.VITE_API_URL;
        let baseUrl = rawUrl || `${window.location.protocol}//${window.location.host}/api/v1`;
        if (rawUrl && !rawUrl.endsWith('/api/v1') && !rawUrl.endsWith('/api/v1/')) {
            baseUrl = rawUrl.endsWith('/') ? `${rawUrl}api/v1` : `${rawUrl}/api/v1`;
        }
        const wsUrl = `${baseUrl.replace('http', 'ws')}/ws?token=${token}`

        socket.value = new WebSocket(wsUrl)

        socket.value.onopen = () => {
            console.log('WebSocket: Connected')
            isConnected.value = true
            if (reconnectInterval.value) {
                clearInterval(reconnectInterval.value)
                reconnectInterval.value = null
            }
        }

        socket.value.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data)
                handleMessage(message)
            } catch (e) {
                console.error('WebSocket: Failed to parse message', e)
            }
        }

        socket.value.onclose = (event) => {
            console.log('WebSocket: Closed', event)
            isConnected.value = false
            socket.value = null
            attemptReconnect()
        }

        socket.value.onerror = (error) => {
            console.error('WebSocket: Error', error)
            socket.value.close()
        }
    }

    function disconnect() {
        if (socket.value) {
            socket.value.close()
            socket.value = null
        }
        isConnected.value = false
        if (reconnectInterval.value) {
            clearInterval(reconnectInterval.value)
            reconnectInterval.value = null
        }
    }

    function attemptReconnect() {
        if (reconnectInterval.value) return

        reconnectInterval.value = setInterval(() => {
            console.log('WebSocket: Attempting reconnect...')
            if (authStore.isAuthenticated) {
                connect()
            } else {
                disconnect()
            }
        }, 5000)
    }

    function handleMessage(message) {
        // Global handler for notification types
        console.log('WS Message:', message)

        switch (message.type) {
            case 'GRADE_ADDED':
                Notify.create({
                    message: `New Grade Posted! Value: ${message.payload.grade_value}`,
                    color: 'info',
                    icon: 'school',
                    position: 'top-right'
                })
                break
            case 'ATTENDANCE_LATE':
            case 'ATTENDANCE_ABSENT':
                Notify.create({
                    message: `Attendance Update: ${message.payload.status}`,
                    color: 'warning',
                    icon: 'warning',
                    position: 'top-right'
                })
                break
            default:
                // Handle generic
                break
        }
    }

    return {
        connect,
        disconnect,
        isConnected
    }
})
