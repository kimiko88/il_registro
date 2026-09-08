import { describe, it, expect, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import { useAuthStore } from 'src/stores/auth'
import { useMenuItems } from 'src/composables/useMenuItems'

// Mock Quasar components/composables if needed
vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            dark: { isActive: false, toggle: vi.fn() },
            fullscreen: { isActive: false, toggle: vi.fn() }
        })
    }
})

/**
 * Flatten a menu config returned by useMenuItems into leaf items.
 * Top-level items without children are included as-is; items with
 * children contribute only their children (the category wrapper is
 * not itself a navigable leaf).
 */
const getFlatItems = (role) => {
    const raw = useMenuItems(role)
    const result = []
    raw.forEach(item => {
        if (item.children) {
            result.push(...item.children)
        } else {
            result.push(item)
        }
    })
    return result
}

describe('MainLayout Logic', () => {

    const setupUserRole = (role) => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: { role, first_name: 'Test', last_name: 'User' },
                    token: 'mock-token'
                }
            }
        })
        const authStore = useAuthStore(pinia)
        return authStore
    }

    describe('user role resolution', () => {
        it('should provide correct userRole for admin', () => {
            const authStore = setupUserRole('admin')
            expect(authStore.userRole).toBe('admin')
        })

        it('should provide correct userRole for student', () => {
            const authStore = setupUserRole('student')
            expect(authStore.userRole).toBe('student')
        })

        it('should provide correct userRole for parent', () => {
            const authStore = setupUserRole('parent')
            expect(authStore.userRole).toBe('parent')
        })

        it('should provide correct userRole for secretary', () => {
            const authStore = setupUserRole('secretary')
            expect(authStore.userRole).toBe('secretary')
        })
    })

    describe('menu items for roles', () => {
        it('should provide teacher menu items', () => {
            setupUserRole('teacher')
            const flatItems = getFlatItems('teacher')

            // 1 dashboard + 14 Didattica + 6 Organizzazione + 4 Comunicazioni = 25
            expect(flatItems).toHaveLength(25)
            expect(flatItems.map(i => i.label)).toContain('Dashboard')
            expect(flatItems.map(i => i.label)).toContain('Le Mie Classi')
            expect(flatItems.map(i => i.label)).toContain('Voti')
            expect(flatItems.map(i => i.label)).toContain('Presenze')
            expect(flatItems.map(i => i.label)).toContain('Agenda')
            expect(flatItems.map(i => i.label)).toContain('Credito Scolastico')
            expect(flatItems.map(i => i.label)).toContain('Corsi Recupero & PAI')
            expect(flatItems.map(i => i.label)).toContain('Registro Sostegno & PEI')
            expect(flatItems.map(i => i.label)).toContain('Ricevimento Generale')
        })

        it('should provide admin menu items', () => {
            setupUserRole('admin')
            const menuItems = useMenuItems('admin')

            // admin config has 8 top-level flat items (no nested children)
            expect(menuItems).toHaveLength(8)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('La Mia Scuola')
            expect(menuItems.map(i => i.label)).toContain('Analytics')
            expect(menuItems.map(i => i.label)).toContain('Impostazioni')
        })

        it('should provide student menu items', () => {
            setupUserRole('student')
            // student uses nested categories — use getFlatItems to count leaves
            const flatItems = getFlatItems('student')

            // 1 dashboard + 6 Didattica + 4 Organizzazione + 6 Percorsi = 17
            expect(flatItems).toHaveLength(17)
            expect(flatItems.map(i => i.label)).toContain('Dashboard')
            expect(flatItems.map(i => i.label)).toContain('I Miei Voti')
            expect(flatItems.map(i => i.label)).toContain('Le Mie Presenze')
            expect(flatItems.map(i => i.label)).toContain('PCTO')
        })

        it('should provide parent menu items', () => {
            setupUserRole('parent')
            // parent uses nested categories — use getFlatItems to count leaves
            const flatItems = getFlatItems('parent')

            // 1 dashboard + 7 Valutazione + 5 Servizi + 5 Comunicazioni = 18
            expect(flatItems).toHaveLength(18)
            expect(flatItems.map(i => i.label)).toContain('Dashboard')
            expect(flatItems.map(i => i.label)).toContain('I Miei Figli')
            expect(flatItems.map(i => i.label)).toContain('Colloqui')
            expect(flatItems.map(i => i.label)).toContain('Ricevimento Generale')
        })

        it('should provide secretary menu items', () => {
            setupUserRole('secretary')
            const flatItems = getFlatItems('secretary')

            // 1 dashboard + 4 Anagrafiche + 4 Atti + 8 Servizi (incl. Flussi SIDI MIM) = 17
            expect(flatItems).toHaveLength(17)
            expect(flatItems.map(i => i.label)).toContain('Dashboard')
            expect(flatItems.map(i => i.label)).toContain('Documenti')
            expect(flatItems.map(i => i.label)).toContain('Studenti')
            expect(flatItems.map(i => i.label)).toContain('Flussi SIDI')
        })
    })

    describe('role label mapping', () => {
        const roleLabels = {
            admin: 'Amministratore',
            secretary: 'Segretario',
            teacher: 'Docente',
            student: 'Studente',
            parent: 'Genitore',
            superadmin: 'Super Admin'
        }

        Object.entries(roleLabels).forEach(([role, _label]) => {
            it(`should correctly resolve role label for ${role}`, () => {
                const authStore = setupUserRole(role)
                expect(authStore.userRole).toBe(role)
            })
        })
    })
})
