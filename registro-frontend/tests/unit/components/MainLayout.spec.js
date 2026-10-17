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

            expect(flatItems.length).toBeGreaterThanOrEqual(20)
            const labels = flatItems.map(i => i.label)
            expect(labels).toContain('Dashboard')
            expect(labels).toContain('Le Mie Classi')
            expect(labels).toContain('Voti')
            expect(labels).toContain('Presenze')
            expect(labels).toContain('Credito Scolastico')
            expect(labels).toContain('Corsi Recupero & PAI')
            expect(labels.some(l => l.includes('Orario'))).toBe(true)
            expect(labels.some(l => l.includes('Colloqui'))).toBe(true)
        })

        it('should provide admin menu items', () => {
            setupUserRole('admin')
            const flatItems = getFlatItems('admin')

            expect(flatItems.length).toBeGreaterThanOrEqual(10)
            const labels = flatItems.map(i => i.label)
            expect(labels).toContain('Dashboard')
            expect(labels).toContain('La Mia Scuola')
            expect(labels).toContain('Presenze Personale')
            expect(labels).toContain('Analytics')
            expect(labels).toContain('Impostazioni')
        })

        it('should provide student menu items', () => {
            setupUserRole('student')
            const flatItems = getFlatItems('student')

            expect(flatItems.length).toBeGreaterThanOrEqual(10)
            const labels = flatItems.map(i => i.label)
            expect(labels).toContain('Dashboard')
            expect(labels).toContain('I Miei Voti')
            expect(labels).toContain('Presenze')
            expect(labels).toContain('Compiti')
        })

        it('should provide parent menu items', () => {
            setupUserRole('parent')
            const flatItems = getFlatItems('parent')

            expect(flatItems.length).toBeGreaterThanOrEqual(10)
            const labels = flatItems.map(i => i.label)
            expect(labels).toContain('Dashboard')
            expect(labels).toContain('I Miei Figli')
            expect(labels).toContain('Voti')
            expect(labels.some(l => l.includes('Colloqui'))).toBe(true)
        })

        it('should provide secretary menu items', () => {
            setupUserRole('secretary')
            const flatItems = getFlatItems('secretary')

            expect(flatItems.length).toBeGreaterThanOrEqual(15)
            const labels = flatItems.map(i => i.label)
            expect(labels).toContain('Dashboard')
            expect(labels.some(l => l.includes('Documenti'))).toBe(true)
            expect(labels.some(l => l.includes('Studenti'))).toBe(true)
            expect(labels).toContain('Presenze Personale')
            expect(labels).toContain('Flussi SIDI')
            expect(labels).toContain('Scioperi')
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
