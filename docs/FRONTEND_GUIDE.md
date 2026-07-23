# Guida allo sviluppo Frontend

Questa guida è rivolta agli sviluppatori che contribuiscono al frontend di RegistroV2. Copre le convenzioni di codice, i pattern architetturali, la gestione dello stato e le best practice per Vue 3 + Quasar.

---

## Indice

- [Principi generali](#principi-generali)
- [Convenzioni di naming](#convenzioni-di-naming)
- [Struttura dei componenti Vue](#struttura-dei-componenti-vue)
- [Composables](#composables)
- [Pinia Stores](#pinia-stores)
- [API Service Layer](#api-service-layer)
- [Routing e protezione route](#routing-e-protezione-route)
- [Gestione errori](#gestione-errori)
- [Testing](#testing)
- [Internazionalizzazione](#internazionalizzazione)
- [Performance](#performance)

---

## Principi generali

1. **Composition API over Options API**: tutti i nuovi componenti usano `<script setup>`.
2. **Un componente, una responsabilità**: componenti piccoli e focalizzati. Se supera ~200 righe, considera di spezzarlo.
3. **Logica nei composables, non nei componenti**: la logica riutilizzabile va in `src/composables/`.
4. **Stato globale in Pinia**: non passare dati tra componenti non collegati via props/events — usa uno store.
5. **Nessuna chiamata HTTP diretta nei componenti**: usa sempre i service in `src/services/`.

---

## Convenzioni di naming

| Elemento | Convenzione | Esempio |
|---|---|---|
| Componenti `.vue` | PascalCase | `GradeTable.vue`, `MyGradesList.vue` |
| Composables | camelCase con `use` prefix | `useGradeEntry.js`, `useMyGrades.js` |
| Pinia stores | camelCase con `use` prefix | `useUserStore.js`, `useGradesStore.js` |
| Services | camelCase con `Service` suffix | `gradeService.js`, `authService.js` |
| Props | camelCase | `studentId`, `classGrades` |
| Emits | kebab-case | `grade-saved`, `form-submit` |
| Route paths | kebab-case | `/teacher/class-grades` |

---

## Struttura dei componenti Vue

Ordine canonico nelle sezioni di un SFC (`Single File Component`):

```vue
<script setup>
// 1. Import (Vue, Quasar, composables, stores, services)
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useGradesStore } from 'stores/useGradesStore'
import { gradeService } from 'services/gradeService'

// 2. Props e Emits
const props = defineProps({
  classId: { type: String, required: true },
  readonly: { type: Boolean, default: false },
})
const emit = defineEmits(['grade-saved', 'grade-deleted'])

// 3. Store e composables
const $q = useQuasar()
const gradesStore = useGradesStore()

// 4. State locale
const loading = ref(false)
const grades = ref([])

// 5. Computed
const hasGrades = computed(() => grades.value.length > 0)

// 6. Methods
async function loadGrades() {
  loading.value = true
  try {
    const { data } = await gradeService.getClassGrades(props.classId)
    grades.value = data
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel caricamento dei voti' })
  } finally {
    loading.value = false
  }
}

// 7. Lifecycle
onMounted(loadGrades)
</script>

<template>
  <!-- Template qui -->
</template>

<style scoped lang="scss">
/* Stili scoped qui */
</style>
```

---

## Composables

I composables in `src/composables/` incapsulano logica riutilizzabile. Ogni composable:
- Ha nome con prefisso `use`
- Restituisce lo stato e le funzioni necessarie
- Non importa direttamente store (li riceve come parametro o li usa internamente solo se necessario)

### Esempio: `useGradeEntry.js`

```js
import { ref, computed } from 'vue'

export function useGradeEntry(initialValue = null) {
  const value = ref(initialValue)
  const comment = ref('')

  // Validazione
  const isValid = computed(() => {
    const n = parseFloat(value.value)
    return !isNaN(n) && n >= 1 && n <= 10
  })

  // Colore visuale (rosso < 6, arancio 6-7, verde > 7)
  const color = computed(() => {
    const n = parseFloat(value.value)
    if (isNaN(n)) return 'grey'
    if (n < 6) return 'negative'
    if (n < 7) return 'warning'
    return 'positive'
  })

  function reset() {
    value.value = null
    comment.value = ''
  }

  return { value, comment, isValid, color, reset }
}
```

---

## Pinia Stores

Gli store in `src/stores/` sono definiti con `defineStore` e usano la **Setup Store syntax** (più vicina alla Composition API):

```js
// stores/useUserStore.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from 'services/authService'

export const useUserStore = defineStore('user', () => {
  // State
  const user = ref(null)
  const accessToken = ref(null)
  const refreshToken = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!accessToken.value)
  const userRole = computed(() => user.value?.role ?? null)
  const hasRole = (role) => userRole.value === role ||
    userRole.value === 'superadmin' ||
    (userRole.value === 'admin' && role !== 'superadmin')

  // Actions
  async function login(credentials) {
    const { data } = await authService.login(credentials)
    user.value = data.user
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token
  }

  async function refreshTokens() {
    const { data } = await authService.refreshToken(refreshToken.value)
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
  }

  return { user, accessToken, refreshToken, isAuthenticated, userRole,
           hasRole, login, refreshTokens, logout }
})
```

> ⚠️ Non usare `localStorage` o `sessionStorage`: il sito è servito in ambienti sandbox che li bloccano. I token vivono solo in memoria durante la sessione.

---

## API Service Layer

### Configurazione Axios (`boot/axios.js`)

```js
import axios from 'axios'
import { useUserStore } from 'stores/useUserStore'

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 15000,
})

// Request interceptor: aggiunge token JWT
api.interceptors.request.use((config) => {
  const store = useUserStore()
  if (store.accessToken) {
    config.headers.Authorization = `Bearer ${store.accessToken}`
  }
  return config
})

// Response interceptor: gestisce scadenza token
let isRefreshing = false
let queue = []

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve) => queue.push(resolve))
          .then(() => api(originalRequest))
      }
      originalRequest._retry = true
      isRefreshing = true
      try {
        await useUserStore().refreshTokens()
        queue.forEach((cb) => cb())
        queue = []
        return api(originalRequest)
      } catch {
        useUserStore().logout()
        return Promise.reject(error)
      } finally {
        isRefreshing = false
      }
    }
    return Promise.reject(error)
  }
)
```

---

## Routing e protezione route

### Definizione route con meta

```js
// router/routes.js
export const routes = [
  {
    path: '/login',
    component: () => import('layouts/AuthLayout.vue'),
    children: [
      { path: '', component: () => import('pages/auth/LoginPage.vue') },
      { path: 'mfa', component: () => import('pages/auth/MfaPage.vue') },
    ],
  },
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: 'teacher/grades',
        component: () => import('pages/teacher/GradesPage.vue'),
        meta: { roles: ['teacher', 'admin', 'superadmin'] },
      },
      {
        path: 'student/my-grades',
        component: () => import('pages/student/MyGradesPage.vue'),
        meta: { roles: ['student'] },
      },
    ],
  },
]
```

### Navigation guard

```js
// router/index.js
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isAuthenticated) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  if (to.meta.roles && !to.meta.roles.includes(userStore.userRole)) {
    return next({ path: '/unauthorized' })
  }

  next()
})
```

---

## Gestione errori

Usa `$q.notify()` di Quasar per i messaggi all'utente. Non mostrare mai messaggi tecnici o stack trace.

```js
import { useQuasar } from 'quasar'
const $q = useQuasar()

try {
  await gradeService.createGrade(payload)
  $q.notify({ type: 'positive', message: 'Voto inserito correttamente' })
} catch (err) {
  const msg = err.response?.data?.error ?? 'Errore durante il salvataggio'
  $q.notify({ type: 'negative', message: msg })
}
```

Codici HTTP backend e messaggio suggerito:

| Status | Messaggio UI |
|---|---|
| `400` | Dati non validi — mostra il campo specifico |
| `401` | Sessione scaduta — redirect a login |
| `403` | Non hai i permessi per questa azione |
| `404` | Risorsa non trovata |
| `409` | Conflitto — es. voto già presente |
| `422` | Errore di validazione — mostra i dettagli dal backend |
| `500` | Errore del server — contatta l'amministratore |

---

## Testing

Usa **Vitest** + **@vue/test-utils** + **happy-dom**.

### Test di un componente

```js
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import GradeTable from 'components/Teacher/GradeTable.vue'

describe('GradeTable', () => {
  it('mostra i voti ricevuti come prop', () => {
    const wrapper = mount(GradeTable, {
      props: { classId: 'class-1' },
      global: {
        plugins: [createTestingPinia()],
      },
    })
    expect(wrapper.find('table').exists()).toBe(true)
  })
})
```

### Test di un composable

```js
import { useGradeEntry } from 'composables/useGradeEntry'

describe('useGradeEntry', () => {
  it('valida correttamente voti nel range 1-10', () => {
    const { value, isValid } = useGradeEntry()
    value.value = 7.5
    expect(isValid.value).toBe(true)
    value.value = 11
    expect(isValid.value).toBe(false)
  })
})
```

### Test di uno store Pinia

```js
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from 'stores/useUserStore'
import { vi } from 'vitest'

beforeEach(() => setActivePinia(createPinia()))

it('isAuthenticated è false dopo logout', () => {
  const store = useUserStore()
  store.accessToken = 'fake-token'
  store.logout()
  expect(store.isAuthenticated).toBe(false)
})
```

---

## Internazionalizzazione

La struttura i18n è già predisposta in `src/boot/i18n.js`. Per aggiungere una lingua:

```
src/
└── i18n/
    ├── it/          # Italiano (default)
    │   └── index.js
    └── en/          # English
        └── index.js
```

Uso nel template:
```vue
<template>
  <span>{{ $t('grades.average') }}</span>
</template>
```

---

## Performance

- **Lazy loading delle route**: tutti i `component: () => import(...)` sono già lazy. Non importare mai pagine direttamente nel router.
- **`v-memo`**: per liste lunghe di voti con aggiornamenti frequenti, usa `v-memo` per evitare re-render inutili.
- **`defineAsyncComponent`**: per componenti pesanti caricati condizionalmente.
- **Chart.js**: passa i dataset come `readonly` e usa `chart.update()` invece di ricreare il chart ad ogni aggiornamento.
- **Quasar Virtual Scroll**: usa `<q-virtual-scroll>` per liste con >100 elementi (registro voti, presenze).
