import { useMenuItems } from 'src/composables/useMenuItems'
import { describe, it, expect } from 'vitest'

describe('useMenuItems', () => {
    describe('admin role', () => {
        it('should return admin menu items', () => {
            const menuItems = useMenuItems('admin')

            expect(menuItems).toHaveLength(4)
            expect(menuItems[0].label).toBe('Dashboard')
            expect(menuItems[1].label).toBe('La Mia Scuola')
            expect(menuItems[2].label).toBe('Analytics')
            expect(menuItems[3].label).toBe('Impostazioni')
        })

        it('should have correct paths for admin', () => {
            const menuItems = useMenuItems('admin')

            expect(menuItems[0].path).toBe('/')
            expect(menuItems[1].path).toBe('/admin/schools')
            expect(menuItems[2].path).toBe('/admin/analytics')
        })

        it('should have exact flag for dashboard', () => {
            const menuItems = useMenuItems('admin')

            expect(menuItems[0].exact).toBe(true)
        })
    })

    describe('secretary role', () => {
        it('should return secretary menu items', () => {
            const menuItems = useMenuItems('secretary')

            expect(menuItems).toHaveLength(8)
            expect(menuItems.map(item => item.label)).toContain('Documenti')
            expect(menuItems.map(item => item.label)).toContain('Studenti')
            expect(menuItems.map(item => item.label)).toContain('Report')
        })
    })

    describe('teacher role', () => {
        it('should return teacher menu items', () => {
            const menuItems = useMenuItems('teacher')

            expect(menuItems).toHaveLength(7)
            expect(menuItems.map(item => item.label)).toContain('Le Mie Classi')
            expect(menuItems.map(item => item.label)).toContain('Voti')
            expect(menuItems.map(item => item.label)).toContain('Presenze')
            expect(menuItems.map(item => item.label)).toContain('Colloqui')
        })

        it('should have correct paths for teacher', () => {
            const menuItems = useMenuItems('teacher')
            const paths = menuItems.map(item => item.path)

            expect(paths).toContain('/teacher/classes')
            expect(paths).toContain('/teacher/grades')
            expect(paths).toContain('/teacher/attendance')
        })
    })

    describe('student role', () => {
        it('should return student menu items', () => {
            const menuItems = useMenuItems('student')

            expect(menuItems).toHaveLength(8)
            expect(menuItems.map(item => item.label)).toContain('I Miei Voti')
            expect(menuItems.map(item => item.label)).toContain('Le Mie Presenze')
            expect(menuItems.map(item => item.label)).toContain('PCTO')
            expect(menuItems.map(item => item.label)).toContain('Orientamento')
        })
    })

    describe('parent role', () => {
        it('should return parent menu items', () => {
            const menuItems = useMenuItems('parent')

            expect(menuItems).toHaveLength(9)
            expect(menuItems.map(item => item.label)).toContain('I Miei Figli')
            expect(menuItems.map(item => item.label)).toContain('Colloqui')
            expect(menuItems.map(item => item.label)).toContain('Supporto')
        })
    })

    describe('menu icons', () => {
        it('should have icons for all menu items', () => {
            const roles = ['admin', 'teacher', 'student', 'parent', 'secretary']

            roles.forEach(role => {
                const menuItems = useMenuItems(role)
                menuItems.forEach(item => {
                    expect(item.icon).toBeDefined()
                    expect(typeof item.icon).toBe('string')
                })
            })
        })
    })

    describe('unknown role', () => {
        it('should return empty array for unknown role', () => {
            const menuItems = useMenuItems('unknown')

            expect(menuItems).toEqual([])
        })

        it('should return empty array for null role', () => {
            const menuItems = useMenuItems(null)

            expect(menuItems).toEqual([])
        })
    })

    describe('menu structure', () => {
        it('all menu items should have required properties', () => {
            const roles = ['admin', 'teacher', 'student', 'parent', 'secretary']

            roles.forEach(role => {
                const menuItems = useMenuItems(role)
                menuItems.forEach(item => {
                    expect(item).toHaveProperty('label')
                    expect(item).toHaveProperty('icon')
                    expect(item).toHaveProperty('path')
                })
            })
        })

        it('dashboard should be first item for all roles', () => {
            const roles = ['admin', 'teacher', 'student', 'parent', 'secretary']

            roles.forEach(role => {
                const menuItems = useMenuItems(role)
                expect(menuItems[0].label).toBe('Dashboard')
                if (role === 'parent') {
                    expect(menuItems[0].path).toBe('/parent')
                } else {
                    expect(menuItems[0].path).toBe('/')
                }
            })
        })
    })
})
