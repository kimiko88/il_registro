import { Notify } from 'quasar'
import { useErrorStore } from '@/stores/error'

/**
 * Custom composable for consistent, central error handling across Vue components.
 */
export function useErrorHandler() {
    const errorStore = useErrorStore()

    /**
     * Handles an error object by logging it, registering it in the error store,
     * and optionally triggering a Quasar toast notification.
     * 
     * @param {Error|Object} err - The error object to handle
     * @param {String} [fallbackMessage] - Optional human fallback message
     * @param {Boolean} [showToast=true] - Whether to display a toast notification
     */
    function handleError(err, fallbackMessage = 'Si è verificato un errore durante l\'operazione.', showToast = true) {
        const message = err?.userMessage || err?.response?.data?.error || err?.message || fallbackMessage

        console.error('[ErrorHandler]', err)

        errorStore.reportError(err, message)

        if (showToast) {
            Notify.create({
                message,
                color: 'negative',
                icon: 'error',
                position: 'top-right',
                timeout: 4000,
                attrs: { role: 'alert' }
            })
        }

        return message
    }

    return {
        handleError
    }
}
