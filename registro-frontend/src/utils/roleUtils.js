/**
 * Shared role-to-dashboard mapping utility.
 * Returns the default landing path for a given user role.
 * Used by both the router guard (guards.js) and login redirect (useAuth.js)
 * to ensure a single source of truth for role-based routing.
 *
 * @param {string | null | undefined} role
 * @returns {string} path to redirect to
 */
export function getUserDashboard(role) {
    if (!role) return '/login'

    const r = role.toLowerCase().trim()

    if (r === 'admin' || r === 'superadmin' || r === 'system_auditor' || r === 'auditor') {
        return '/admin/dashboard'
    }
    if (r === 'teacher' || r === 'coordinator' || r === 'docente' || r === 'vice_principal' || r === 'collaboratore_vicario') {
        return '/teacher'
    }
    if (r === 'student' || r === 'studente') return '/student'
    if (r === 'parent' || r === 'genitore') return '/parent'
    if (r === 'secretary' || r === 'principal' || r === 'dirigente_scolastico' || r === 'staff') {
        return '/secretary'
    }
    if ([
        'dsga',
        'assistente_amministrativo',
        'collaboratore_ds',
        'collaboratore_scolastico',
        'assistente_alunni',
        'assistente_personale',
        'assistente_contabilita',
        'assistente_protocollo',
        'assistente_sportello',
        'assistente_tecnico',
        'responsabile_servizio'
    ].includes(r)) {
        return '/ata'
    }
    if (r === 'responsabile_gestione_documentale' || r === 'responsabile_conservazione') {
        return '/secretary/documents'
    }
    if (r === 'dpo') return '/admin/audit-logs'

    // Fallback: authenticated but unknown role → root dashboard
    return '/'
}
