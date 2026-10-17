/**
 * useMenuBadges — Badge/contatori dinamici per le voci del menu laterale
 *
 * Ritorna un oggetto reattivo con i contatori per i badge.
 * Ogni chiave corrisponde al campo `badge` definito in useMenuItems.js.
 *
 * Il polling avviene ogni POLL_INTERVAL ms; si ferma automaticamente
 * quando il componente viene smontato (onUnmounted).
 */
import { ref, onUnmounted } from 'vue'
import api from '@/services/api'

const POLL_INTERVAL = 5 * 60 * 1000 // 5 minuti

// Singleton: un solo polling attivo per tutta l'app
let _globalBadges = null
let _pollTimer = null
let _refCount = 0

function createBadges() {
    return {
        // Comunicazioni / messaggi non letti
        unreadMessages: ref(0),
        // Voti da inserire (docente)
        pendingGrades: ref(0),
        // Colloqui da confermare
        pendingColloqui: ref(0),
        // Sostituzioni in attesa di copertura
        pendingSubstitutions: ref(0),
        // Compiti assegnati non ancora consegnati (studente)
        pendingHomework: ref(0),
        // Personale assente oggi
        absentStaff: ref(0),
    }
}

async function fetchBadges(badges, role) {
    try {
        const normRole = (role || '').toLowerCase()

        // Messaggi non letti — tutti i ruoli
        try {
            const res = await api.get('/notifications/unread-count')
            badges.unreadMessages.value = res?.data?.count ?? 0
        } catch { /* silenzioso */ }

        // Badge specifici per ruolo docente
        if (['teacher', 'docente', 'coordinator', 'vice_principal'].includes(normRole)) {
            try {
                const res = await api.get('/teacher/grades/pending-count')
                badges.pendingGrades.value = res?.data?.count ?? 0
            } catch { /* silenzioso */ }

            try {
                const res = await api.get('/teacher/colloqui/pending-count')
                badges.pendingColloqui.value = res?.data?.count ?? 0
            } catch { /* silenzioso */ }
        }

        // Badge specifici per studente
        if (['student', 'studente'].includes(normRole)) {
            try {
                const res = await api.get('/student/homework/pending-count')
                badges.pendingHomework.value = res?.data?.count ?? 0
            } catch { /* silenzioso */ }
        }

        // Badge specifici per genitore
        if (['parent', 'genitore'].includes(normRole)) {
            try {
                const res = await api.get('/parent/colloqui/pending-count')
                badges.pendingColloqui.value = res?.data?.count ?? 0
            } catch { /* silenzioso */ }
        }

        // Badge specifici per ruoli amministrativi/ATA
        const adminRoles = [
            'admin', 'secretary', 'principal', 'vice_principal',
            'dsga', 'collaboratore_ds', 'assistente_personale',
            'assistente_amministrativo'
        ]
        if (adminRoles.includes(normRole)) {
            try {
                const res = await api.get('/ata/substitutions/pending-count')
                badges.pendingSubstitutions.value = res?.data?.count ?? 0
            } catch { /* silenzioso */ }

            try {
                const res = await api.get('/ata/attendance/absent-count')
                badges.absentStaff.value = res?.data?.count ?? 0
            } catch { /* silenzioso */ }
        }
    } catch { /* silenzioso: non bloccare la UI per i badge */ }
}

/**
 * Composable principale.
 * @param {string} role - Ruolo corrente dell'utente
 * @returns {{ badges: Object, getBadge: Function }}
 */
export function useMenuBadges(role) {
    // Inizializza singleton al primo uso
    if (!_globalBadges) {
        _globalBadges = createBadges()
    }
    _refCount++

    // Primo fetch immediato
    fetchBadges(_globalBadges, role)

    // Polling ricorrente
    if (!_pollTimer) {
        _pollTimer = setInterval(() => {
            fetchBadges(_globalBadges, role)
        }, POLL_INTERVAL)
    }

    onUnmounted(() => {
        _refCount--
        if (_refCount <= 0) {
            clearInterval(_pollTimer)
            _pollTimer = null
            _globalBadges = null
            _refCount = 0
        }
    })

    /**
     * Restituisce il valore del badge per una chiave.
     * Ritorna undefined (nessun badge) se il contatore è 0.
     * @param {string} key
     * @returns {number|undefined}
     */
    function getBadge(key) {
        if (!key || !_globalBadges?.[key]) return undefined
        const val = _globalBadges[key].value
        return val > 0 ? val : undefined
    }

    return {
        badges: _globalBadges,
        getBadge,
    }
}

export default useMenuBadges
