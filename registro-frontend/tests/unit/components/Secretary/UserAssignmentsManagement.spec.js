import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useUserAssignments, canAssignDuty } from 'src/composables/useUserAssignments'
import { useMenuItems } from 'src/composables/useMenuItems'

describe('User Assignments & Governance Authority Tests', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    describe('Authority Matrix (canAssignDuty)', () => {
        it('allows Dirigente Scolastico (principal) to assign pedagogical and institutional duties', () => {
            expect(canAssignDuty('principal', 'coordinatore_classe')).toBe(true)
            expect(canAssignDuty('principal', 'segretario_verbale')).toBe(true)
            expect(canAssignDuty('principal', 'referente_inclusione')).toBe(true)
            expect(canAssignDuty('principal', 'referente_progetto')).toBe(true)
            expect(canAssignDuty('principal', 'responsabile_dipartimento')).toBe(true)
            expect(canAssignDuty('principal', 'tutor_orientamento')).toBe(true)
            expect(canAssignDuty('principal', 'animatore_digitale')).toBe(true)
            expect(canAssignDuty('principal', 'dpo')).toBe(true)
            expect(canAssignDuty('principal', 'responsabile_gestione_documentale')).toBe(true)
            expect(canAssignDuty('principal', 'responsabile_conservazione')).toBe(true)
        })

        it('allows DSGA to assign ATA branches and service responsibilities', () => {
            expect(canAssignDuty('dsga', 'responsabile_servizio')).toBe(true)
            expect(canAssignDuty('dsga', 'assistente_alunni')).toBe(true)
            expect(canAssignDuty('dsga', 'assistente_personale')).toBe(true)
            expect(canAssignDuty('dsga', 'assistente_contabilita')).toBe(true)
            expect(canAssignDuty('dsga', 'assistente_protocollo')).toBe(true)
            expect(canAssignDuty('dsga', 'assistente_sportello')).toBe(true)
            expect(canAssignDuty('dsga', 'assistente_tecnico')).toBe(true)

            // DSGA must NOT assign teacher pedagogical functions
            expect(canAssignDuty('dsga', 'referente_inclusione')).toBe(false)
            expect(canAssignDuty('dsga', 'responsabile_dipartimento')).toBe(false)
            expect(canAssignDuty('dsga', 'animatore_digitale')).toBe(false)
        })

        it('allows Admin and SuperAdmin full management delegation', () => {
            expect(canAssignDuty('admin', 'coordinatore_classe')).toBe(true)
            expect(canAssignDuty('admin', 'referente_inclusione')).toBe(true)
            expect(canAssignDuty('admin', 'responsabile_servizio')).toBe(true)
            expect(canAssignDuty('superadmin', 'dpo')).toBe(true)
        })

        it('allows Secretary to record class coordinators and meeting secretaries', () => {
            expect(canAssignDuty('secretary', 'coordinatore_classe')).toBe(true)
            expect(canAssignDuty('secretary', 'segretario_verbale')).toBe(true)
            expect(canAssignDuty('secretary', 'referente_inclusione')).toBe(false)
            expect(canAssignDuty('secretary', 'responsabile_servizio')).toBe(false)
        })

        it('denies teachers, students, parents, and dpo from assigning duties', () => {
            expect(canAssignDuty('teacher', 'coordinatore_classe')).toBe(false)
            expect(canAssignDuty('student', 'coordinatore_classe')).toBe(false)
            expect(canAssignDuty('parent', 'coordinatore_classe')).toBe(false)
            expect(canAssignDuty('dpo', 'coordinatore_classe')).toBe(false)
        })
    })

    describe('useUserAssignments Composable', () => {
        it('correctly derives active duty flags for teacher profile', () => {
            const authStore = useAuthStore()
            authStore.user = {
                id: 'teacher-1',
                role: 'teacher',
                assignments: [
                    { id: 'a1', assignment_type: 'coordinatore_classe', scope_id: 'class-1', title: 'Coord 1A', is_active: true },
                    { id: 'a2', assignment_type: 'coordinatore_classe', scope_id: 'class-2', title: 'Coord 2A', is_active: true },
                    { id: 'a3', assignment_type: 'referente_inclusione', title: 'Inclusione', is_active: true },
                    { id: 'a4', assignment_type: 'animatore_digitale', title: 'PNSD', is_active: true },
                    { id: 'a5', assignment_type: 'responsabile_dipartimento', title: 'Dip. Lettere', is_active: false } // inactive
                ]
            }

            const {
                isCoordinator,
                coordinatedClassIds,
                isInclusionLead,
                isDigitalAnimator,
                isDepartmentHead
            } = useUserAssignments()

            expect(isCoordinator.value).toBe(true)
            expect(coordinatedClassIds.value).toEqual(['class-1', 'class-2'])
            expect(isInclusionLead.value).toBe(true)
            expect(isDigitalAnimator.value).toBe(true)
            expect(isDepartmentHead.value).toBe(false) // Inactive duty is ignored
        })
    })

    describe('Dynamic Menu Generation with Incarichi', () => {
        it('enables coordinator items and injects guidance tutor and digital animator', () => {
            const assignments = [
                { assignment_type: 'coordinatore_classe', scope_id: 'c1', is_active: true },
                { assignment_type: 'tutor_orientamento', is_active: true },
                { assignment_type: 'animatore_digitale', is_active: true }
            ]

            const items = useMenuItems('teacher', assignments)
            const flat = []
            items.forEach(cat => {
                if (cat.children) flat.push(...cat.children)
                else flat.push(cat)
            })

            const labels = flat.map(c => c.label)
            expect(labels).toContain('Coordinamento')
            expect(labels).toContain('Tutor Orientamento')
            expect(labels).toContain('Team Digitale & E-Learning')
        })
    })
})
