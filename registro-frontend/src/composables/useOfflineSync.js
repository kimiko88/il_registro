import { watch, getCurrentInstance, onMounted } from 'vue'
import { Notify } from 'quasar'
import api from '@/services/api'
import { useOutboxStore } from '@/stores/outbox'
import { useNetworkStatus } from '@/composables/useNetworkStatus'
import { i18n } from '@/i18n'

// Singleton flag: the reactive watcher is attached once globally the first time
// any component calls useOfflineSync() — subsequent calls share the same watch.
let _watcherInstalled = false

export function useOfflineSync() {
    const outboxStore = useOutboxStore()
    const { isOnline } = useNetworkStatus()

    async function handleAutoSync() {
        if (!outboxStore.hasPending || !isOnline.value) return

        const res = await outboxStore.syncQueue()
        if (res && res.synced > 0) {
            const msg = i18n?.global?.t
                ? i18n.global.t('offlineBanner.syncedSuccess', { count: res.synced })
                : `Connessione ripristinata: ${res.synced} operazioni sincronizzate con successo.`

            Notify.create({
                type: 'positive',
                icon: 'cloud_done',
                message: msg,
                timeout: 4000
            })
        }
    }

    function installWatcher() {
        if (_watcherInstalled) return
        _watcherInstalled = true

        // Watch the reactive isOnline signal: triggers sync whenever connectivity
        // is restored. This replaces the module-level window.addEventListener('online')
        // approach which was never cleaned up and could fire multiple times.
        watch(isOnline, (online, wasOnline) => {
            if (online && !wasOnline) {
                handleAutoSync()
            }
        })
    }

    if (getCurrentInstance()) {
        onMounted(installWatcher)
    } else {
        installWatcher()
    }

    async function executeWithOfflineQueue(requestConfig, options = {}) {
        const title = options.title || `${requestConfig.method?.toUpperCase() || 'POST'} ${requestConfig.url}`

        // 1. If currently offline, enqueue immediately.
        if (typeof navigator !== 'undefined' && navigator.onLine === false) {
            const id = await outboxStore.enqueue({
                url: requestConfig.url,
                method: requestConfig.method || 'post',
                data: requestConfig.data,
                params: requestConfig.params,
                title
            })

            const offlineMsg = i18n?.global?.t
                ? i18n.global.t('offlineBanner.enqueuedOffline')
                : 'Nessuna connessione. Operazione salvata in locale: verrà inviata appena tornerai online.'

            Notify.create({
                type: 'warning',
                icon: 'wifi_off',
                message: offlineMsg,
                timeout: 4500
            })

            return { offline: true, id, enqueued: true }
        }

        // 2. Try online execution.
        try {
            let response
            const method = (requestConfig.method || 'post').toLowerCase()
            if (typeof api[method] === 'function') {
                if (['post', 'put', 'patch'].includes(method)) {
                    const config = (requestConfig.params || requestConfig.headers)
                        ? { params: requestConfig.params, headers: requestConfig.headers }
                        : undefined
                    response = config
                        ? await api[method](requestConfig.url, requestConfig.data, config)
                        : await api[method](requestConfig.url, requestConfig.data)
                } else {
                    response = await api[method](requestConfig.url, requestConfig)
                }
            } else if (typeof api.request === 'function') {
                response = await api.request(requestConfig)
            } else {
                throw new Error(`api.${method} is not a function`)
            }
            return response
        } catch (err) {
            const isNetworkError =
                !err.response ||
                err.code === 'ERR_NETWORK' ||
                (err.message && /network|fetch|timeout/i.test(err.message))

            if (isNetworkError) {
                const id = await outboxStore.enqueue({
                    url: requestConfig.url,
                    method: requestConfig.method || 'post',
                    data: requestConfig.data,
                    params: requestConfig.params,
                    title
                })

                const savedMsg = i18n?.global?.t
                    ? i18n.global.t('offlineBanner.enqueuedOffline')
                    : 'Errore di connessione. Operazione salvata in locale: verrà sincronizzata automaticamente.'

                Notify.create({
                    type: 'warning',
                    icon: 'cloud_queue',
                    message: savedMsg,
                    timeout: 4500
                })

                return { offline: true, id, enqueued: true }
            }

            // Regular API error (4xx/5xx): re-throw for the component to handle.
            throw err
        }
    }

    return {
        outboxStore,
        isOnline,
        executeWithOfflineQueue,
        triggerSync: handleAutoSync
    }
}
