import { defineStore } from 'pinia'
import { ref } from 'vue'
import { Notify } from 'quasar'

export const useErrorStore = defineStore('error', () => {
    const globalError = ref(null)
    const errorHistory = ref([])

    function reportError(error, customMessage = null) {
        const message = customMessage || error?.userMessage || error?.response?.data?.error || error?.message || 'Si è verificato un errore inaspettato.'
        const errObj = {
            id: Date.now(),
            timestamp: new Date(),
            message,
            originalError: error,
            code: error?.response?.data?.code || error?.code || 'UNKNOWN_ERROR',
            status: error?.response?.status || null
        }

        globalError.value = errObj
        errorHistory.value.unshift(errObj)

        if (errorHistory.value.length > 50) {
            errorHistory.value.pop()
        }

        // Display Quasar toast notification
        Notify.create({
            message,
            color: 'negative',
            icon: 'error',
            position: 'top-right',
            timeout: 5000,
            actions: [
                { label: 'Chiudi', color: 'white' }
            ]
        })
    }

    function clearGlobalError() {
        globalError.value = null
    }

    return {
        globalError,
        errorHistory,
        reportError,
        clearGlobalError
    }
})
