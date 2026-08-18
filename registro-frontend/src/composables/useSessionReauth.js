/**
 * useSessionReauth — Singleton per la re-autenticazione in-page.
 *
 * Quando il token scade e il refresh fallisce, invece di navigare a /login
 * e perdere lo stato della pagina corrente, questo composable mostra un
 * dialog modale di re-autenticazione. L'utente inserisce solo la password
 * (l'email è già nota dal profilo). Al successo il token viene aggiornato e
 * la request originale può essere ritentata automaticamente.
 *
 * Pattern: ref globali condivisi tra tutte le istanze (singleton).
 */
import { ref } from 'vue'

// ── Stato condiviso (singleton) ──────────────────────────────────────────────
const showDialog = ref(false)
const userEmail = ref('')

// Resolve/reject della Promise corrente — impostati da triggerReauth()
let _resolve = null
let _reject = null

// ── Composable ────────────────────────────────────────────────────────────────
export function useSessionReauth() {
    /**
     * Apre il dialog di re-autenticazione e restituisce una Promise.
     * La Promise si risolve con il nuovo access_token al successo del login,
     * oppure viene rigettata se l'utente chiude il dialog senza autenticarsi.
     *
     * @param {string} email - Email dell'utente loggato (precompilata nel form)
     * @returns {Promise<string>} - Risolve col nuovo token o rigetta
     */
    function triggerReauth(email) {
        userEmail.value = email || ''
        showDialog.value = true

        return new Promise((resolve, reject) => {
            _resolve = resolve
            _reject = reject
        })
    }

    /**
     * Chiamato dal SessionReauthDialog al successo del login.
     * Risolve la Promise pendente col nuovo token.
     * @param {string} newToken
     */
    function resolveReauth(newToken) {
        showDialog.value = false
        if (_resolve) {
            _resolve(newToken)
            _resolve = null
            _reject = null
        }
    }

    /**
     * Chiamato dal SessionReauthDialog quando l'utente annulla / chiude.
     * Rigetta la Promise pendente.
     */
    function cancelReauth() {
        showDialog.value = false
        if (_reject) {
            _reject(new Error('reauth_cancelled'))
            _resolve = null
            _reject = null
        }
    }

    return {
        showDialog,
        userEmail,
        triggerReauth,
        resolveReauth,
        cancelReauth
    }
}
