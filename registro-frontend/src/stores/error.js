import { defineStore } from 'pinia'
import { ref } from 'vue'
import { Notify } from 'quasar'
import { i18n } from '@/i18n'

export const useErrorStore = defineStore('error', () => {
    const globalError = ref(null)
    const errorHistory = ref([])

    function reportError(error, customMessage = null) {
        // i18n may not be available yet during early boot errors — use safe fallback
        const t = i18n?.global?.t
        const defaultMsg = t ? t('common.error') : 'Si è verificato un errore inaspettato.'
        const closeLabel = t ? t('common.close') : 'Chiudi'

        const message = customMessage || error?.userMessage || error?.response?.data?.error || error?.message || defaultMsg
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
                { label: closeLabel, color: 'white' }
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
