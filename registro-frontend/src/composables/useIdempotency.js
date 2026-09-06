/**
 * useIdempotency.js
 *
 * Composable per la generazione e gestione di chiavi di idempotenza per
 * operazioni mutanti critiche (inserimento voti, presenze, firme, pagamenti).
 *
 * Una chiave di idempotenza è un UUID v4 univoco per ogni intenzione di scrittura:
 * inviandola come header `Idempotency-Key` al backend, il server può riconoscere
 * richieste duplicate (es. doppio click, retry di rete) e restituire la risposta
 * già elaborata senza rieseguire la logica di business.
 *
 * Ciclo di vita:
 *  1. `generateKey()` — crea e memorizza una nuova chiave (invalida la precedente)
 *  2. `currentKey` — la chiave attiva da passare all'API call
 *  3. Dopo una risposta di successo, chiamare `rotateKey()` per preparare
 *     la chiave del prossimo submit (previene riuso accidentale).
 *
 * @example
 *   const { generateKey, currentKey, rotateKey } = useIdempotency('grade-add')
 *
 *   async function saveGrade() {
 *     generateKey()
 *     await gradeService.saveGrade(data, currentKey.value)
 *     rotateKey()
 *   }
 */
import { ref, readonly } from 'vue'

/**
 * Genera un UUID v4 compatibile con tutti i browser moderni.
 * Usa `crypto.randomUUID()` se disponibile (Chrome 92+, Firefox 95+, Safari 15.4+),
 * con fallback manuale per ambienti più datati.
 */
function generateUUID() {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  // Fallback RFC 4122 v4
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0
    const v = c === 'x' ? r : (r & 0x3) | 0x8
    return v.toString(16)
  })
}

/**
 * Prefisso opzionale per debugging: permette di identificare in Network DevTools
 * a quale operazione appartiene la chiave (es. "grade-add", "attendance-save").
 *
 * @param {string} [prefix='op'] - Prefisso leggibile
 * @returns {{ currentKey: Readonly<Ref<string>>, generateKey: Function, rotateKey: Function }}
 */
export function useIdempotency(prefix = 'op') {
  const _key = ref(_buildKey(prefix))

  function _buildKey(pfx) {
    return `${pfx}:${generateUUID()}`
  }

  /**
   * Genera una nuova chiave e la rende disponibile in `currentKey`.
   * Chiamare **prima** di ogni submit.
   */
  function generateKey() {
    _key.value = _buildKey(prefix)
  }

  /**
   * Ruota la chiave dopo un submit andato a buon fine.
   * Equivalente a `generateKey()` ma semanticamente chiarisce l'intenzione.
   */
  function rotateKey() {
    _key.value = _buildKey(prefix)
  }

  return {
    /** La chiave di idempotenza corrente (readonly). */
    currentKey: readonly(_key),
    generateKey,
    rotateKey,
  }
}

export default useIdempotency
