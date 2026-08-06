import { useQuasar } from 'quasar'
import { ref } from 'vue'

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
export function useUndoToast(timeoutMs = 15000) {
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

        return new Promise((resolve) => {
            let dismissed = false
            let timer = null

            // Build countdown label
            const totalSeconds = Math.round(timeoutMs / 1000)
            let remaining = totalSeconds

            let dismiss = $q.notify({
                type: options.color === 'negative' ? 'negative' : 'positive',
                color: options.color || 'dark',
                textColor: 'white',
                icon: options.icon || 'check_circle',
                message,
                position: options.position || 'bottom',
                timeout: timeoutMs + 200,     // slightly longer than our timer
                actions: [
                    {
                        label: `Annulla (${remaining}s)`,
                        color: 'yellow',
                        handler: async () => {
                            if (dismissed) return
                            dismissed = true
                            clearInterval(timer)
                            undoClicked.value = true
                            try {
                                await undoFn()
                                $q.notify({
                                    type: 'warning',
                                    icon: 'undo',
                                    message: 'Azione annullata',
                                    position: options.position || 'bottom',
                                    timeout: 2000
                                })
                            } catch (err) {
                                $q.notify({
                                    type: 'negative',
                                    message: "Errore durante l'annullamento",
                                    position: options.position || 'bottom',
                                    timeout: 2500
                                })
                            }
                            resolve(true)
                        }
                    }
                ]
            })

            // Update countdown label every second
            timer = setInterval(() => {
                remaining--
                if (remaining <= 0 || dismissed) {
                    clearInterval(timer)
                    if (!dismissed) {
                        dismissed = true
                        if (typeof dismiss === 'function') dismiss()
                        resolve(false)
                    }
                    return
                }
                // Re-issue notify with updated label (Quasar doesn't support reactive actions,
                // so we close and re-open the same notify to refresh the countdown)
                if (typeof dismiss === 'function') dismiss()
                dismiss = $q.notify({
                    type: options.color === 'negative' ? 'negative' : 'positive',
                    color: options.color || 'dark',
                    textColor: 'white',
                    icon: options.icon || 'check_circle',
                    message,
                    position: options.position || 'bottom',
                    timeout: (remaining * 1000) + 200,
                    actions: [
                        {
                            label: `Annulla (${remaining}s)`,
                            color: 'yellow',
                            handler: async () => {
                                if (dismissed) return
                                dismissed = true
                                clearInterval(timer)
                                undoClicked.value = true
                                try {
                                    await undoFn()
                                    $q.notify({
                                        type: 'warning',
                                        icon: 'undo',
                                        message: 'Azione annullata',
                                        position: options.position || 'bottom',
                                        timeout: 2000
                                    })
                                } catch {
                                    $q.notify({
                                        type: 'negative',
                                        message: "Errore durante l'annullamento",
                                        position: options.position || 'bottom',
                                        timeout: 2500
                                    })
                                }
                                resolve(true)
                            }
                        }
                    ]
                })
            }, 1000)
        })
    }

    return { notifyWithUndo, undoClicked }
}
