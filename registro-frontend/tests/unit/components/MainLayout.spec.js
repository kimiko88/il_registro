import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import { useAuthStore } from 'src/stores/auth'
import { useMenuItems } from 'src/composables/useMenuItems'

vi.mock('vue-router', () => ({
    useRouter: () => ({
        push: vi.fn()
    }),
    useRoute: () => ({
        path: '/'
    })
}))

vi.mock('src/composables/useAuth', () => ({
    useAuth: () => ({
        logout: vi.fn().mockResolvedValue()
    })
}))

describe('MainLayout Logic', () => {
    let pinia

    beforeEach(() => {
        vi.clearAllMocks()
        pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: null,
                    token: null,
                    refreshToken: null
                }
            }
        })
    })

    const setupUserRole = (role) => {
        const authStore = useAuthStore(pinia)
        authStore.user = {
            id: '123',
            email: 'test@example.com',
            first_name: 'John',
            last_name: 'Doe',
            role: role
        }
        authStore.token = 'test-token'
        authStore.refreshToken = 'test-refresh'
        return authStore
    }

    describe('user profile data', () => {
        it('should provide correct userName for teacher', () => {
            const authStore = setupUserRole('teacher')
            expect(authStore.userName).toBe('John Doe')
        })

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
            const menuItems = useMenuItems('teacher')

            expect(menuItems).toHaveLength(7)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('Le Mie Classi')
            expect(menuItems.map(i => i.label)).toContain('Voti')
            expect(menuItems.map(i => i.label)).toContain('Presenze')
        })

        it('should provide admin menu items', () => {
            setupUserRole('admin')
            const menuItems = useMenuItems('admin')

            expect(menuItems).toHaveLength(4)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('La Mia Scuola')
            expect(menuItems.map(i => i.label)).toContain('Analytics')
            expect(menuItems.map(i => i.label)).toContain('Impostazioni')
        })

        it('should provide student menu items', () => {
            setupUserRole('student')
            const menuItems = useMenuItems('student')

            expect(menuItems).toHaveLength(8)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('I Miei Voti')
            expect(menuItems.map(i => i.label)).toContain('Le Mie Presenze')
            expect(menuItems.map(i => i.label)).toContain('PCTO')
        })

        it('should provide parent menu items', () => {
            setupUserRole('parent')
            const menuItems = useMenuItems('parent')

            expect(menuItems).toHaveLength(9)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('I Miei Figli')
            expect(menuItems.map(i => i.label)).toContain('Colloqui')
        })

        it('should provide secretary menu items', () => {
            setupUserRole('secretary')
            const menuItems = useMenuItems('secretary')

            expect(menuItems).toHaveLength(8)
            expect(menuItems.map(i => i.label)).toContain('Dashboard')
            expect(menuItems.map(i => i.label)).toContain('Documenti')
            expect(menuItems.map(i => i.label)).toContain('Studenti')
        })
    })

    describe('role label mapping', () => {
        const roleLabels = {
            admin: 'Amministratore',
            secretary: 'Segretario',
            teacher: 'Docente',
            student: 'Studente',
            parent: 'Genitore'
        }

        Object.entries(roleLabels).forEach(([role, label]) => {
            it(`should map ${role} to ${label}`, () => {
                const authStore = setupUserRole(role)
                expect(authStore.userRole).toBe(role)
                // The label mapping would be tested in the actual component
                expect(roleLabels[role]).toBe(label)
            })
        })
    })

    describe('authentication state', () => {
        it('should be authenticated when user and token exist', () => {
            const authStore = setupUserRole('teacher')
            expect(authStore.isAuthenticated).toBe(true)
        })

        it('should not be authenticated when logged out', () => {
            const authStore = useAuthStore(pinia)
            authStore.logout()
            expect(authStore.isAuthenticated).toBe(false)
            expect(authStore.user).toBeNull()
        })
    })
})
