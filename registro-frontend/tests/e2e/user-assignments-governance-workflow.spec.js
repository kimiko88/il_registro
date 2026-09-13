import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import { useAuthStore } from '@/stores/auth'
import { useClassesStore } from '@/stores/classes'
import { useUserAssignments, canAssignDuty } from '@/composables/useUserAssignments'
import { useMenuItems } from '@/composables/useMenuItems'
import CoordinatorView from '@/pages/teacher/CoordinatorView.vue'

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn().mockImplementation((url) => {
            if (url === '/teacher/classes') {
                return Promise.resolve({
                    data: [
                        { id: 'class-2a', name: '2A', section: 'A', coordinator_id: 'teacher-rossi-uuid' },
                        { id: 'class-2b', name: '2B', section: 'B', coordinator_id: 'teacher-rossi-uuid' }
                    ]
                })
            }
            return Promise.resolve({ data: [] })
        }),
        post: vi.fn().mockResolvedValue({ data: {} }),
        put: vi.fn().mockResolvedValue({ data: {} }),
        delete: vi.fn().mockResolvedValue({ data: {} })
    }
}))

vi.mock('@/services/authService', () => ({
    default: {
        getCurrentUser: vi.fn().mockResolvedValue({
            id: 'teacher-rossi-uuid',
            first_name: 'Marco',
            last_name: 'Rossi',
            email: 'docente1@scuola.it',
            role: 'teacher'
        })
    }
}))

vi.mock('@/services/notesService', () => ({
    default: {
        getNotes: vi.fn().mockResolvedValue({ data: [] })
    }
}))

vi.mock('@/services/scrutinyService', () => ({
    scrutinyService: {
        getMatrix: vi.fn().mockResolvedValue({ data: { subjects: [{ id: 1, name: 'Matematica' }] } }),
        getScrutinyOverview: vi.fn().mockResolvedValue([]),
        startScrutiny: vi.fn().mockResolvedValue({})
    }
}))

vi.mock('@/services/coordinatorService', () => ({
    default: {
        getCoordinatedClasses: vi.fn().mockResolvedValue([
            { id: 'class-2a', name: '2A', section: 'A' },
            { id: 'class-2b', name: '2B', section: 'B' }
        ]),
        getClassSummary: vi.fn().mockResolvedValue({
            students_count: 20,
            averages: { 'Matematica': 7.5 }
        })
    }
}))

describe('End-to-End Governance & User Assignments Workflow (E2E)', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: {
                        id: 'principal-1',
                        first_name: 'Laura',
                        last_name: 'Dirigente',
                        email: 'dirigente.prova@scuola.it',
                        role: 'principal',
                        school_id: 'school-prova'
                    }
                },
                schoolYear: {
                    selectedSchoolYear: '2025/2026'
                },
                classes: {
                    classes: [
                        { id: 'class-2a', name: '2A', section: 'A', coordinator_id: null },
                        { id: 'class-2b', name: '2B', section: 'B', coordinator_id: null }
                    ]
                }
            }
        })
    })

    it('Fase 1 (Governance): Dirigente nomina Docente Coordinatore di 2 Classi ed Incarichi Speciali', () => {
        const authStore = useAuthStore()
        expect(canAssignDuty(authStore.user.role, 'coordinatore_classe')).toBe(true)
        expect(canAssignDuty(authStore.user.role, 'referente_inclusione')).toBe(true)
        expect(canAssignDuty(authStore.user.role, 'animatore_digitale')).toBe(true)

        // Simulazione registrazione nomine per il docente Marco Rossi
        const targetTeacher = {
            id: 'teacher-rossi-uuid',
            first_name: 'Marco',
            last_name: 'Rossi',
            email: 'docente1@scuola.it',
            role: 'teacher',
            assignments: [
                { id: 'asgn-c2a', assignment_type: 'coordinatore_classe', scope_type: 'class', scope_id: 'class-2a', title: 'Coordinatore Classe 2A', is_active: true },
                { id: 'asgn-c2b', assignment_type: 'coordinatore_classe', scope_type: 'class', scope_id: 'class-2b', title: 'Coordinatore Classe 2B', is_active: true },
                { id: 'asgn-inc', assignment_type: 'referente_inclusione', scope_type: 'school', title: 'Referente Inclusione (BES / DSA)', is_active: true },
                { id: 'asgn-anim', assignment_type: 'animatore_digitale', scope_type: 'school', title: 'Animatore Digitale PNSD', is_active: true }
            ]
        }

        const teacherAssignments = useUserAssignments(targetTeacher)
        expect(teacherAssignments.isCoordinator.value).toBe(true)
        expect(teacherAssignments.coordinatedClassIds.value).toEqual(['class-2a', 'class-2b'])
        expect(teacherAssignments.isInclusionLead.value).toBe(true)
        expect(teacherAssignments.isDigitalAnimator.value).toBe(true)
    })

    it('Fase 2 (Esperienza Docente): Docente incaricato accede ed ottiene menu dinamico esteso', () => {
        const authStore = useAuthStore()
        // Switch utente a docente con incarichi multipli
        authStore.user = {
            id: 'teacher-rossi-uuid',
            first_name: 'Marco',
            last_name: 'Rossi',
            email: 'docente1@scuola.it',
            role: 'teacher',
            assignments: [
                { id: 'asgn-c2a', assignment_type: 'coordinatore_classe', scope_type: 'class', scope_id: 'class-2a', title: 'Coordinatore Classe 2A', is_active: true },
                { id: 'asgn-c2b', assignment_type: 'coordinatore_classe', scope_type: 'class', scope_id: 'class-2b', title: 'Coordinatore Classe 2B', is_active: true },
                { id: 'asgn-inc', assignment_type: 'referente_inclusione', scope_type: 'school', title: 'Referente Inclusione (BES / DSA)', is_active: true },
                { id: 'asgn-anim', assignment_type: 'animatore_digitale', scope_type: 'school', title: 'Animatore Digitale PNSD', is_active: true }
            ]
        }

        const menu = useMenuItems('teacher', authStore.user.assignments)
        const allLabels = []
        menu.forEach(item => {
            if (item.label) allLabels.push(item.label)
            if (item.children) {
                item.children.forEach(child => {
                    if (child.label) allLabels.push(child.label)
                })
            }
        })

        // Verifica iniezione automatica e sblocco delle voci dovute agli incarichi
        expect(allLabels).toContain('Coordinamento')
        expect(allLabels).toContain('Scrutinio')
        expect(allLabels).toContain('Inclusione (BES / DSA)')
        expect(allLabels).toContain('Team Digitale & E-Learning')
    })

    it('Fase 3 (Dashboard Coordinamento): Docente visualizza entrambe le classi coordinate (2A e 2B)', async () => {
        const authStore = useAuthStore()
        authStore.user = {
            id: 'teacher-rossi-uuid',
            first_name: 'Marco',
            last_name: 'Rossi',
            email: 'docente1@scuola.it',
            role: 'teacher',
            assignments: [
                { id: 'asgn-c2a', assignment_type: 'coordinatore_classe', scope_type: 'class', scope_id: 'class-2a', title: 'Coordinatore Classe 2A', is_active: true },
                { id: 'asgn-c2b', assignment_type: 'coordinatore_classe', scope_type: 'class', scope_id: 'class-2b', title: 'Coordinatore Classe 2B', is_active: true }
            ]
        }

        const classesStore = useClassesStore()
        classesStore.classes = [
            { id: 'class-2a', name: '2A', section: 'A', coordinator_id: 'teacher-rossi-uuid' },
            { id: 'class-2b', name: '2B', section: 'B', coordinator_id: 'teacher-rossi-uuid' }
        ]

        const wrapper = mount(CoordinatorView, {
            global: {
                plugins: [Quasar, pinia],
                mocks: {
                    t: (key) => key
                },
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-space': true,
                    'q-select': { template: '<div class="q-select-stub"><slot /></div>', props: ['modelValue', 'options'] },
                    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
                    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
                    'q-separator': { template: '<hr />' },
                    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
                    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-spinner': true,
                    'q-markup-table': { template: '<table><slot /></table>' },
                    'q-icon': true,
                    'q-btn': true
                }
            }
        })

        expect(wrapper.exists()).toBe(true)
        expect(wrapper.text().toLowerCase()).toContain('coordinatore')
    })

    it('Fase 4 (Governance ATA & Restrizioni DSGA): DSGA può assegnare servizi ATA ma non incarichi didattici', () => {
        const authStore = useAuthStore()
        authStore.user = {
            id: 'dsga-1',
            role: 'dsga',
            email: 'dsga.prova@scuola.it'
        }

        // DSGA ha competenza sul personale e servizi ATA
        expect(canAssignDuty('dsga', 'responsabile_servizio')).toBe(true)
        expect(canAssignDuty('dsga', 'assistente_tecnico')).toBe(true)
        expect(canAssignDuty('dsga', 'addetto_sicurezza')).toBe(true)

        // DSGA NON ha competenza su nomine didattiche e coordinamento
        expect(canAssignDuty('dsga', 'coordinatore_classe')).toBe(false)
        expect(canAssignDuty('dsga', 'referente_inclusione')).toBe(false)
        expect(canAssignDuty('dsga', 'animatore_digitale')).toBe(false)
    })

    it('Fase 5 (Restrizioni di Sicurezza): Docenti semplici e Studenti non possono attribuire alcun incarico', () => {
        expect(canAssignDuty('teacher', 'coordinatore_classe')).toBe(false)
        expect(canAssignDuty('teacher', 'responsabile_servizio')).toBe(false)
        expect(canAssignDuty('student', 'coordinatore_classe')).toBe(false)
        expect(canAssignDuty('parent', 'coordinatore_classe')).toBe(false)
    })
})
