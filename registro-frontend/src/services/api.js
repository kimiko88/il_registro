import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { Notify } from 'quasar'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || '/api/v1',
    timeout: 10000
})

api.interceptors.request.use(config => {
    const authStore = useAuthStore()
    if (authStore.token) {
        config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
})

api.interceptors.response.use(
    response => response,
    error => {
        const status = error.response ? error.response.status : null

        if (status === 401) {
            const authStore = useAuthStore()
            authStore.logout()
            window.location.reload()
        } else if (status >= 500) {
            Notify.create({
                type: 'negative',
                message: 'Server Error. Please try again later.'
            })
        }

        return Promise.reject(error)
    }
)

export default api
