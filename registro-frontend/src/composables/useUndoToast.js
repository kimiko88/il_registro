import { useQuasar } from 'quasar'
import { ref } from 'vue'
import { i18n } from '@/i18n'

/**
 * useUndoToast — Composable for showing an undo-able action toast.
 *
 * Usage:
 *   const { notifyWithUndo } = useUndoToast()
 *
 *   // After marking an absence or saving a grade:
 *   const undone = await notifyWithUndo('Assenza segnata', async () => {
 *       await restoreAttendance(studentId)
 *   })
 *   if (!undone) {
 *       // proceed with permanent save
 *   }
 *
 * @param {number} timeoutMs  — How long the undo window lasts (default 5000ms)
 */
export function useUndoToast(timeoutMs = 6000) {
    const $q = useQuasar()
    // Track whether the user clicked Undo
    const undoClicked = ref(false)

    /**
     * Show a toast with an "Annulla" button.
     * Returns a Promise<boolean> that resolves to:
     *   true  → user clicked Undo (action should be reversed)
     *   false → toast expired (action confirmed)
     *
     * @param {string}   message  — Text to show in the toast
     * @param {Function} undoFn   — Async function to call on Undo
     * @param {object}   options  — Optional overrides { color, icon, position }
     */
    async function notifyWithUndo(message, undoFn, options = {}) {
        undoClicked.value = false
        const t = i18n?.global?.t

        const getCancelLabel = (sec) => t ? t('composables.undo.cancelWithSeconds', { remaining: sec }) : `Annulla (${sec}s)`
        const getActionCancelled = () => t ? t('composables.undo.actionCancelled') : 'Azione annullata'
        const getCancelError = () => t ? t('composables.undo.cancelError') : 'Errore durante l\'annullamento'

        return new Promise((resolve) => {
            let dismissed = false
            let timer = null

            const totalSeconds = Math.round(timeoutMs / 1000)
            let remaining = totalSeconds

            const handleUndo = async () => {
                if (dismissed) return
                dismissed = true
                if (timer) clearInterval(timer)
                undoClicked.value = true
                if (typeof notifHandle === 'function') notifHandle()
                try {
                    await undoFn()
                    $q.notify({
                        type: 'warning',
                        icon: 'undo',
                        message: getActionCancelled(),
                        position: options.position || 'bottom',
                        timeout: 2000
                    })
                } catch (err) {
                    $q.notify({
                        type: 'negative',
                        message: getCancelError(),
                        position: options.position || 'bottom',
                        timeout: 2500
                    })
                }
                resolve(true)
            }

            let notifHandle = $q.notify({
                type: options.color === 'negative' ? 'negative' : 'positive',
                color: options.color || 'dark',
                textColor: 'white',
                icon: options.icon || 'check_circle',
                message,
                position: options.position || 'bottom',
                timeout: timeoutMs + 200,
                actions: [
                    {
                        label: getCancelLabel(remaining),
                        color: 'yellow',
                        handler: handleUndo
                    }
                ]
            })

            // Update countdown label every second using notification update()
            timer = setInterval(() => {
                remaining--
                if (remaining <= 0 || dismissed) {
                    clearInterval(timer)
                    if (!dismissed) {
                        dismissed = true
                        if (typeof notifHandle === 'function') notifHandle()
                        resolve(false)
                    }
                    return
                }

                if (typeof notifHandle === 'function' && typeof notifHandle.update === 'function') {
                    notifHandle.update({
                        actions: [
                            {
                                label: getCancelLabel(remaining),
                                color: 'yellow',
                                handler: handleUndo
                            }
                        ]
                    })
                } else {
                    if (typeof notifHandle === 'function') notifHandle()
                    notifHandle = $q.notify({
                        type: options.color === 'negative' ? 'negative' : 'positive',
                        color: options.color || 'dark',
                        textColor: 'white',
                        icon: options.icon || 'check_circle',
                        message,
                        position: options.position || 'bottom',
                        timeout: (remaining * 1000) + 200,
                        actions: [
                            {
                                label: getCancelLabel(remaining),
                                color: 'yellow',
                                handler: handleUndo
                            }
                        ]
                    })
                }
            }, 1000)
        })
    }

    return { notifyWithUndo, undoClicked }
}
