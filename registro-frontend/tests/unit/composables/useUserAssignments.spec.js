import { describe, it, expect, beforeEach } from 'vitest'
import { ref } from 'vue'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useUserAssignments, canAssignDuty } from '@/composables/useUserAssignments'

describe('useUserAssignments Composable Unit Tests', () => {
    let authStore

    beforeEach(() => {
        setActivePinia(createPinia())
        authStore = useAuthStore()
    })

    describe('Authority Matrix: canAssignDuty()', () => {
        it('allows superadmin and admin to assign any duty', () => {
            expect(canAssignDuty('superadmin', 'coordinatore_classe')).toBe(true)
            expect(canAssignDuty('superadmin', 'responsabile_servizio')).toBe(true)
            expect(canAssignDuty('admin', 'referente_inclusione')).toBe(true)
            expect(canAssignDuty('admin', 'animatore_digitale')).toBe(true)
        })

        it('allows Principal and Vice Principal to assign educational and institutional duties', () => {
            const educationalDuties = [
                'coordinatore_classe',
                'segretario_consiglio',
                'referente_inclusione',
                'referente_progetto',
                'responsabile_dipartimento',
                'tutor_orientamento',
                'animatore_digitale'
            ]
            for (const duty of educationalDuties) {
                expect(canAssignDuty('principal', duty)).toBe(true)
                expect(canAssignDuty('vice_principal', duty)).toBe(true)
            }
        })

        it('allows DSGA to assign ATA and service responsibilities, but denies educational duties', () => {
            expect(canAssignDuty('dsga', 'responsabile_servizio')).toBe(true)
            expect(canAssignDuty('dsga', 'assistente_tecnico')).toBe(true)
            expect(canAssignDuty('dsga', 'assistente_alunni')).toBe(true)
            expect(canAssignDuty('dsga', 'addetto_sicurezza')).toBe(true)

            expect(canAssignDuty('dsga', 'coordinatore_classe')).toBe(false)
            expect(canAssignDuty('dsga', 'referente_inclusione')).toBe(false)
            expect(canAssignDuty('dsga', 'animatore_digitale')).toBe(false)
        })

        it('allows Secretary to record class coordinators and meeting secretaries, but denies other duties', () => {
            expect(canAssignDuty('secretary', 'coordinatore_classe')).toBe(true)
            expect(canAssignDuty('secretary', 'segretario_consiglio')).toBe(true)

            expect(canAssignDuty('secretary', 'referente_inclusione')).toBe(false)
            expect(canAssignDuty('secretary', 'responsabile_servizio')).toBe(false)
        })

        it('denies duty assignment powers to Teachers, Students, Parents, and invalid roles', () => {
            const forbiddenRoles = ['teacher', 'student', 'parent', 'dpo', 'system_auditor', null, '']
            for (const role of forbiddenRoles) {
                expect(canAssignDuty(role, 'coordinatore_classe')).toBe(false)
                expect(canAssignDuty(role, 'referente_inclusione')).toBe(false)
                expect(canAssignDuty(role, 'responsabile_servizio')).toBe(false)
            }
        })
    })

    describe('Reactive Duty Flags and Multi-Class Coordination', () => {
        it('returns false/empty for users without assignments', () => {
            authStore.user = { id: 'u1', role: 'teacher', assignments: [] }
            const {
                isCoordinator,
                coordinatedClassIds,
                isInclusionLead,
                isDigitalAnimator,
                isServiceManager
            } = useUserAssignments()

            expect(isCoordinator.value).toBe(false)
            expect(coordinatedClassIds.value).toEqual([])
            expect(isInclusionLead.value).toBe(false)
            expect(isDigitalAnimator.value).toBe(false)
            expect(isServiceManager.value).toBe(false)
        })

        it('detects multiple coordinated classes and all educational duties', () => {
            const userWithDuties = ref({
                id: 'teacher-valerio',
                role: 'teacher',
                assignments: [
                    { id: '1', assignment_type: 'coordinatore_classe', scope_type: 'class', scope_id: 'class-2a', is_active: true },
                    { id: '2', assignment_type: 'coordinatore_classe', scope_type: 'class', scope_id: 'class-2b', is_active: true },
                    { id: '3', assignment_type: 'segretario_consiglio', is_active: true },
                    { id: '4', assignment_type: 'referente_inclusione', is_active: true },
                    { id: '5', assignment_type: 'referente_progetto', is_active: true },
                    { id: '6', assignment_type: 'responsabile_dipartimento', is_active: true },
                    { id: '7', assignment_type: 'tutor_orientamento', is_active: true },
                    { id: '8', assignment_type: 'animatore_digitale', is_active: true }
                ]
            })

            const {
                isCoordinator,
                coordinatedClassIds,
                isCouncilSecretary,
                isInclusionLead,
                isProjectLead,
                isDepartmentHead,
                isGuidanceTutor,
                isDigitalAnimator
            } = useUserAssignments(userWithDuties)

            expect(isCoordinator.value).toBe(true)
            expect(coordinatedClassIds.value).toEqual(['class-2a', 'class-2b'])
            expect(isCouncilSecretary.value).toBe(true)
            expect(isInclusionLead.value).toBe(true)
            expect(isProjectLead.value).toBe(true)
            expect(isDepartmentHead.value).toBe(true)
            expect(isGuidanceTutor.value).toBe(true)
            expect(isDigitalAnimator.value).toBe(true)
        })

        it('ignores inactive duties (is_active: false)', () => {
            const userWithInactive = ref({
                id: 't-inactive',
                role: 'teacher',
                assignments: [
                    { id: '1', assignment_type: 'coordinatore_classe', scope_id: 'class-2a', is_active: false },
                    { id: '2', assignment_type: 'referente_inclusione', is_active: false }
                ]
            })

            const { isCoordinator, coordinatedClassIds, isInclusionLead } = useUserAssignments(userWithInactive)

            expect(isCoordinator.value).toBe(false)
            expect(coordinatedClassIds.value).toEqual([])
            expect(isInclusionLead.value).toBe(false)
        })

        it('detects ATA service managers and safety officers', () => {
            const ataUser = ref({
                id: 'tec-1',
                role: 'assistente_tecnico',
                assignments: [
                    { id: '1', assignment_type: 'responsabile_servizio', scope_id: 'Lab IT', is_active: true },
                    { id: '2', assignment_type: 'addetto_antincendio', is_active: true }
                ]
            })

            const { isServiceManager, isSafetyOfficer } = useUserAssignments(ataUser)

            expect(isServiceManager.value).toBe(true)
            expect(isSafetyOfficer.value).toBe(true)
        })
    })
})
