import { describe, it, expect } from 'vitest'
import { useMenuItems } from '@/composables/useMenuItems'

describe('useMenuItems — Dynamic Role-Based Menu Generation', () => {
  it('returns superadmin menu items', () => {
    const items = useMenuItems('superadmin')
    expect(items.length).toBeGreaterThan(0)
    expect(items.some(i => i.path === '/admin/audit-logs')).toBe(true)
    expect(items.some(i => i.path === '/admin/monitoring')).toBe(true)
  })

  it('returns admin menu items', () => {
    const items = useMenuItems('admin')
    expect(items.length).toBeGreaterThan(0)
    expect(items.some(i => i.path === '/admin/schools')).toBe(true)
    expect(items.some(i => i.path === '/admin/school-settings')).toBe(true)
  })

  it('returns secretary menu structure with categories', () => {
    const items = useMenuItems('secretary')
    expect(items.length).toBeGreaterThan(0)
    expect(items.some(i => i.category === 'Anagrafiche & Classi')).toBe(true)
    expect(items.some(i => i.category === 'Atti & Certificati')).toBe(true)
    expect(items.some(i => i.category === 'Servizi & Report')).toBe(true)
  })

  it('returns teacher menu structure with categories and coordinator markers', () => {
    const items = useMenuItems('teacher')
    expect(items.length).toBeGreaterThan(0)
    const didattica = items.find(i => i.category === 'Didattica & Valutazione')
    expect(didattica).toBeDefined()
    expect(didattica.children.some(c => c.path === '/teacher/grades')).toBe(true)
    expect(didattica.children.some(c => c.coordinatorOnly === true)).toBe(true)
  })

  it('returns student menu items', () => {
    const items = useMenuItems('student')
    expect(items.length).toBeGreaterThan(0)
    const didattica = items.find(i => i.category === 'Didattica & Valutazione')
    expect(didattica.children.some(c => c.path === '/student/grades')).toBe(true)
    expect(didattica.children.some(c => c.path === '/student/attendance')).toBe(true)
  })

  it('returns parent menu items', () => {
    const items = useMenuItems('parent')
    expect(items.length).toBeGreaterThan(0)
    const didattica = items.find(i => i.category === 'Valutazione & Didattica')
    expect(didattica.children.some(c => c.path === '/parent/children')).toBe(true)
    expect(didattica.children.some(c => c.path === '/parent/grades')).toBe(true)
  })

  it('handles principal role mapped to secretary menu', () => {
    const principalItems = useMenuItems('principal')
    expect(principalItems.length).toBeGreaterThan(0)
    expect(principalItems.some(i => i.category === 'Anagrafiche & Classi')).toBe(true)
  })

  it('handles vice_principal with dedicated executive governance and teaching duties', () => {
    const vpItems = useMenuItems('vice_principal')
    expect(vpItems.length).toBeGreaterThan(0)
    expect(vpItems.some(i => i.category === 'Presidenza & Vicariato')).toBe(true)
    expect(vpItems.some(i => i.category === 'Didattica & Le Mie Classi')).toBe(true)
    const presidenza = vpItems.find(i => i.category === 'Presidenza & Vicariato')
    expect(presidenza.children.some(c => c.path === '/secretary/substitutions')).toBe(true)
    const didattica = vpItems.find(i => i.category === 'Didattica & Le Mie Classi')
    expect(didattica.children.some(c => c.path === '/teacher/classes')).toBe(true)
    expect(didattica.children.some(c => c.path === '/teacher/grades')).toBe(true)
  })

  it('handles coordinator role mapped to teacher menu', () => {
    const coordItems = useMenuItems('coordinator')
    expect(coordItems.length).toBeGreaterThan(0)
    expect(coordItems.some(i => i.category === 'Didattica & Valutazione')).toBe(true)
  })

  it('handles system_auditor mapped to superadmin menu', () => {
    const auditorItems = useMenuItems('system_auditor')
    expect(auditorItems.length).toBeGreaterThan(0)
    expect(auditorItems.some(i => i.path === '/admin/audit-logs')).toBe(true)
  })

  it('handles unknown or empty role gracefully', () => {
    expect(useMenuItems(null)).toEqual([])
    expect(useMenuItems('')).toEqual([])
    expect(useMenuItems('unknown_role_xyz')).toEqual([])
  })
})
