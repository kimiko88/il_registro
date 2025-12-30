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
            primary: '#1976d2',
            secondary: '#26a69a',
            accent: '#9c27b0',
            dark: '#1d1d1d'
        }
    }
})

app.mount('#app')
