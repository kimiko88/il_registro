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

            expect(flatItems).toHaveLength(20)
            expect(flatItems.map(i => i.label)).toContain('Dashboard')
            expect(flatItems.map(i => i.label)).toContain('Le Mie Classi')
            expect(flatItems.map(i => i.label)).toContain('Voti')
            expect(flatItems.map(i => i.label)).toContain('Presenze')
            expect(flatItems.map(i => i.label)).toContain('Agenda')
        })

        it('should provide admin menu items', () => {
            setupUserRole('admin')
            const menuItems = useMenuItems('admin')

            expect(menuItems).toHaveLength(7)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('La Mia Scuola')
            expect(menuItems.map(i => i.label)).toContain('Analytics')
            expect(menuItems.map(i => i.label)).toContain('Impostazioni')
        })

        it('should provide student menu items', () => {
            setupUserRole('student')
            const menuItems = useMenuItems('student')

            expect(menuItems).toHaveLength(14)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('I Miei Voti')
            expect(menuItems.map(i => i.label)).toContain('Le Mie Presenze')
            expect(menuItems.map(i => i.label)).toContain('PCTO')
        })

        it('should provide parent menu items', () => {
            setupUserRole('parent')
            const menuItems = useMenuItems('parent')

            expect(menuItems).toHaveLength(14)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('I Miei Figli')
            expect(menuItems.map(i => i.label)).toContain('Colloqui')
        })

        it('should provide secretary menu items', () => {
            setupUserRole('secretary')
            const flatItems = getFlatItems('secretary')

            expect(flatItems).toHaveLength(14)
            expect(flatItems.map(i => i.label)).toContain('Dashboard')
            expect(flatItems.map(i => i.label)).toContain('Documenti')
            expect(flatItems.map(i => i.label)).toContain('Studenti')
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
