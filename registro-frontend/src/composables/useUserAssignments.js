import { computed } from 'vue'
import { useAuthStore } from 'src/stores/auth'

/**
 * Italian School Governance Authority Matrix
 * Enforces who can assign which duty.
 */
export function canAssignDuty(actorRole, dutyType) {
    if (!actorRole || !dutyType) return false

    // SuperAdmin and Admin have full management delegation
    if (actorRole === 'superadmin' || actorRole === 'admin') {
        return true
    }

    // Dirigente Scolastico (Principal) assigns all pedagogical & institutional duties
    if (actorRole === 'principal' || actorRole === 'vice_principal') {
        return [
            'coordinatore_classe',
            'coordinator',
            'segretario_consiglio',
            'segretario_verbale',
            'collaboratore_ds',
            'referente_plesso',
            'funzione_strumentale',
            'referente_inclusione',
            'referente_progetto',
            'responsabile_dipartimento',
            'tutor_orientamento',
            'tutor_pcto',
            'animatore_digitale',
            'team_innovazione',
            'referente_bullismo',
            'rspp',
            'aspp',
            'addetto_antincendio',
            'addetto_primosoccorso',
            'responsabile_gestione_documentale',
            'responsabile_conservazione',
            'dpo'
        ].includes(dutyType)
    }

    // DSGA assigns Administrative/Technical/Auxiliary duties & ATA service managers
    if (actorRole === 'dsga') {
        return [
            'assistente_alunni',
            'assistente_personale',
            'assistente_contabilita',
            'assistente_protocollo',
            'assistente_sportello',
            'assistente_tecnico',
            'responsabile_servizio',
            'addetto_sicurezza',
            'addetto_antincendio',
            'addetto_primosoccorso'
        ].includes(dutyType)
    }

    // Secretary / Assistente Personale can record class coordinator or meeting secretary assignments
    if (actorRole === 'secretary' || actorRole === 'assistente_personale') {
        return [
            'coordinatore_classe',
            'coordinator',
            'segretario_consiglio',
            'segretario_verbale'
        ].includes(dutyType)
    }

    return false
}

export function useUserAssignments(userRef = null) {
    const authStore = useAuthStore()

    const targetUser = computed(() => {
        if (userRef && userRef.value) return userRef.value
        if (userRef && typeof userRef === 'object' && !('value' in userRef)) return userRef
        return authStore.user || {}
    })

    const assignments = computed(() => {
        return targetUser.value?.assignments || []
    })

    const activeAssignments = computed(() => {
        return assignments.value.filter(a => a.is_active !== false)
    })

    const hasAssignment = (type) => {
        return activeAssignments.value.some(a => a.assignment_type === type)
    }

    const isCoordinator = computed(() => {
        return hasAssignment('coordinatore_classe') || hasAssignment('coordinator')
    })

    const coordinatedClassIds = computed(() => {
        return activeAssignments.value
            .filter(a => (a.assignment_type === 'coordinatore_classe' || a.assignment_type === 'coordinator') && a.scope_id)
            .map(a => a.scope_id)
    })

    const isCouncilSecretary = computed(() => hasAssignment('segretario_verbale') || hasAssignment('segretario_consiglio'))
    const isInclusionLead = computed(() => hasAssignment('referente_inclusione'))
    const isProjectLead = computed(() => hasAssignment('referente_progetto'))
    const isDepartmentHead = computed(() => hasAssignment('responsabile_dipartimento'))
    const isGuidanceTutor = computed(() => hasAssignment('tutor_orientamento'))
    const isDigitalAnimator = computed(() => hasAssignment('animatore_digitale'))
    const isServiceManager = computed(() => hasAssignment('responsabile_servizio'))
    const isSafetyOfficer = computed(() => hasAssignment('rspp') || hasAssignment('aspp') || hasAssignment('addetto_antincendio') || hasAssignment('addetto_primosoccorso'))

    return {
        targetUser,
        assignments,
        activeAssignments,
        hasAssignment,
        isCoordinator,
        coordinatedClassIds,
        isCouncilSecretary,
        isInclusionLead,
        isProjectLead,
        isDepartmentHead,
        isGuidanceTutor,
        isDigitalAnimator,
        isServiceManager,
        isSafetyOfficer,
        canAssignDuty: (dutyType) => canAssignDuty(authStore.user?.role, dutyType)
    }
}

export default useUserAssignments
