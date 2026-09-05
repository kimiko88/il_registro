import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const DEFAULT_TIMEOUT_SECONDS = 30 * 60 // 30 minutes
const WARNING_BEFORE_SECONDS = 60 // Show warning 60s before logout

export function useInactivityTimer(timeoutSeconds = DEFAULT_TIMEOUT_SECONDS) {
    const authStore = useAuthStore()
    const router = useRouter()

    const showWarningDialog = ref(false)
    const secondsRemaining = ref(WARNING_BEFORE_SECONDS)

    let lastActivityTimestamp = Date.now()
    let checkInterval = null
    let countdownInterval = null

    function resetActivity() {
        lastActivityTimestamp = Date.now()
        if (showWarningDialog.value) {
            stayLoggedIn()
        }
    }

    function stayLoggedIn() {
        showWarningDialog.value = false
        lastActivityTimestamp = Date.now()
        secondsRemaining.value = WARNING_BEFORE_SECONDS
        if (countdownInterval) {
            clearInterval(countdownInterval)
            countdownInterval = null
        }
    }

    function logoutNow() {
        showWarningDialog.value = false
        if (checkInterval) clearInterval(checkInterval)
        if (countdownInterval) clearInterval(countdownInterval)
        authStore.logout()
        if (router && typeof router.push === 'function') {
            router.push('/login?reason=session_timeout')
        }
    }

    function checkInactivity() {
        if (!authStore.isAuthenticated) return

        const elapsedSeconds = Math.floor((Date.now() - lastActivityTimestamp) / 1000)
        const timeUntilTimeout = timeoutSeconds - elapsedSeconds

        if (timeUntilTimeout <= 0) {
            logoutNow()
        } else if (timeUntilTimeout <= WARNING_BEFORE_SECONDS && !showWarningDialog.value) {
            showWarningDialog.value = true
            secondsRemaining.value = timeUntilTimeout
            startCountdown()
        }
    }

    function startCountdown() {
        if (countdownInterval) clearInterval(countdownInterval)
        countdownInterval = setInterval(() => {
            const elapsedSeconds = Math.floor((Date.now() - lastActivityTimestamp) / 1000)
            const remaining = timeoutSeconds - elapsedSeconds
            if (remaining <= 0) {
                logoutNow()
            } else {
                secondsRemaining.value = remaining
            }
        }, 1000)
    }

    const activityEvents = ['mousemove', 'mousedown', 'keydown', 'touchstart', 'scroll']

    function startTracking() {
        if (typeof window === 'undefined') return
        activityEvents.forEach((ev) => window.addEventListener(ev, resetActivity, { passive: true }))
        if (checkInterval) clearInterval(checkInterval)
        checkInterval = setInterval(checkInactivity, 5000)
    }

    function stopTracking() {
        if (typeof window === 'undefined') return
        activityEvents.forEach((ev) => window.removeEventListener(ev, resetActivity))
        if (checkInterval) clearInterval(checkInterval)
        if (countdownInterval) clearInterval(countdownInterval)
    }

    onMounted(() => {
        startTracking()
    })

    onUnmounted(() => {
        stopTracking()
    })

    return {
        showWarningDialog,
        secondsRemaining,
        stayLoggedIn,
        logoutNow,
        resetActivity
    }
}
