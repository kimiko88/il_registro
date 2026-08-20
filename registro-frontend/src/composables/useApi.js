import { ref } from 'vue'
import api from '@/services/api'
import { Notify } from 'quasar'
import { i18n } from '@/i18n'

export function useApi(url) {
    const data = ref(null)
    const error = ref(null)
    const loading = ref(false)

    async function fetch() {
        loading.value = true
        error.value = null
        try {
            const response = await api.get(url)
            data.value = response.data
        } catch (err) {
            error.value = err
            const t = i18n?.global?.t
            Notify.create({
                type: 'negative',
                message: err?.userMessage || (t ? t('common.error') : 'Failed to fetch data')
            })
        } finally {
            loading.value = false
        }
    }

    return { data, error, loading, fetch }
}
