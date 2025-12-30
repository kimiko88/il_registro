import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { Quasar, Notify, Dialog, Loading } from 'quasar'
import router from './router'
import App from './App.vue'

// Import Quasar css
import '@quasar/extras/material-icons/material-icons.css'
import 'quasar/src/css/index.sass'

// Global styles
import './assets/styles/globals.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Quasar, {
    plugins: {
        Notify,
        Dialog,
        Loading
    },
    config: {
        brand: {
            primary: '#4F46E5',  // Indigo 600
            secondary: '#06B6D4', // Cyan 500
            accent: '#F59E0B',   // Amber 500

            dark: '#1E293B',     // Slate 900
            'dark-page': '#0F172A', // Slate 950

            positive: '#10B981', // Emerald 500
            negative: '#EF4444', // Red 500
            info: '#3B82F6',     // Blue 500
            warning: '#F59E0B'   // Amber 500
        }
    }
})

app.mount('#app')
