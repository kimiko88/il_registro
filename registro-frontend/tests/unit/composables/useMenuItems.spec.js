import { useMenuItems } from 'src/composables/useMenuItems'
import { describe, it, expect } from 'vitest'

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

describe('useMenuItems', () => {
    describe('admin role', () => {
        it('should return admin menu items', () => {
            const menuItems = useMenuItems('admin')

            expect(menuItems).toHaveLength(8)
            expect(menuItems[0].label).toBe('Dashboard')
            expect(menuItems[1].label).toBe('La Mia Scuola')
            expect(menuItems[2].label).toBe('Gestione Utenti')
            expect(menuItems[3].label).toBe('Gestione Sostituzioni')
            expect(menuItems[4].label).toBe('Feature Flags & Istituto')
        })

        it('should have correct paths for admin', () => {
            const menuItems = useMenuItems('admin')

            expect(menuItems[0].path).toBe('/')
            expect(menuItems[1].path).toBe('/admin/schools')
            expect(menuItems[2].path).toBe('/admin/users')
            expect(menuItems[3].path).toBe('/secretary/substitutions')
            expect(menuItems[4].path).toBe('/admin/school-settings')
        })

        it('should have exact flag for dashboard', () => {
            const menuItems = useMenuItems('admin')

            expect(menuItems[0].exact).toBe(true)
        })
    })

    describe('secretary role', () => {
        it('should return secretary menu items', () => {
            const flatItems = getFlatItems('secretary')

            expect(flatItems).toHaveLength(17)
            expect(flatItems.map(item => item.label)).toContain('Documenti')
            expect(flatItems.map(item => item.label)).toContain('Studenti')
            expect(flatItems.map(item => item.label)).toContain('Flussi SIDI')
            expect(flatItems.map(item => item.label)).toContain('Report')
        })
    })

    describe('teacher role', () => {
        it('should return teacher menu items', () => {
            const flatItems = getFlatItems('teacher')

            expect(flatItems.length).toBeGreaterThan(15)
            expect(flatItems.map(item => item.label)).toContain('Le Mie Classi')
            expect(flatItems.map(item => item.label)).toContain('Voti')
            expect(flatItems.map(item => item.label)).toContain('Presenze')
            expect(flatItems.map(item => item.label)).toContain('Colloqui')
            expect(flatItems.map(item => item.label)).toContain('Agenda')
        })

        it('should have correct paths for teacher', () => {
            const flatItems = getFlatItems('teacher')
            const paths = flatItems.map(item => item.path)

            expect(paths).toContain('/teacher/classes')
            expect(paths).toContain('/teacher/grades')
            expect(paths).toContain('/teacher/attendance')
        })
    })

    describe('student role', () => {
        it('should return student menu items', () => {
            const flatItems = getFlatItems('student')

            expect(flatItems.length).toBeGreaterThan(10)
            expect(flatItems.map(item => item.label)).toContain('I Miei Voti')
            expect(flatItems.map(item => item.label)).toContain('Le Mie Presenze')
            expect(flatItems.map(item => item.label)).toContain('PCTO')
            expect(flatItems.map(item => item.label)).toContain('Orientamento')
        })
    })

    describe('parent role', () => {
        it('should return parent menu items', () => {
            const flatItems = getFlatItems('parent')

            expect(flatItems.length).toBeGreaterThan(10)
            expect(flatItems.map(item => item.label)).toContain('I Miei Figli')
            expect(flatItems.map(item => item.label)).toContain('Colloqui')
            expect(flatItems.map(item => item.label)).toContain('Supporto')
        })
    })

    describe('menu icons', () => {
        it('should have icons for all menu items', () => {
            const roles = ['admin', 'teacher', 'student', 'parent', 'secretary']

            roles.forEach(role => {
                const flatItems = getFlatItems(role)
                flatItems.forEach(item => {
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
                const flatItems = getFlatItems(role)
                flatItems.forEach(item => {
                    expect(item).toHaveProperty('label')
                    expect(item).toHaveProperty('icon')
                    expect(item).toHaveProperty('path')
                })
            })
        })

        it('dashboard should be first item for all roles', () => {
            const roles = ['admin', 'teacher', 'student', 'parent', 'secretary']

            roles.forEach(role => {
                const flatItems = getFlatItems(role)
                expect(flatItems[0].label).toBe('Dashboard')
                if (role === 'parent') {
                    expect(flatItems[0].path).toBe('/parent')
                } else {
                    expect(flatItems[0].path).toBe('/')
                }
            })
        })
    })
})
