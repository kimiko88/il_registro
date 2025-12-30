import { ref } from 'vue'
import api from '@/services/api'
import { Notify } from 'quasar'

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
            Notify.create({ type: 'negative', message: 'Failed to fetch data' })
        } finally {
            loading.value = false
        }
    }

    return { data, error, loading, fetch }
}
