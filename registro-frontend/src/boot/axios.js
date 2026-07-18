import axios from 'axios'
import { boot } from 'quasar/wrappers'

const getBaseURL = () => {
    const rawUrl = import.meta.env.VITE_API_URL;
    if (!rawUrl) {
        return 'http://localhost:8080/api/v1';
    }
    if (rawUrl.endsWith('/api/v1') || rawUrl.endsWith('/api/v1/')) {
        return rawUrl;
    }
    return rawUrl.endsWith('/') ? `${rawUrl}api/v1` : `${rawUrl}/api/v1`;
};

const api = axios.create({ baseURL: getBaseURL() });

export default boot(({ app }) => {
    app.config.globalProperties.$axios = axios
    app.config.globalProperties.$api = api
})

export { api }
