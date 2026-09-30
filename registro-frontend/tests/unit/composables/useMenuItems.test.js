import { describe, it, expect } from 'vitest'
import { useMenuItems } from '@/composables/useMenuItems'

const hasPath = (items, path) =>
  items.some(i => i.path === path || (i.children && i.children.some(c => c.path === path)))

const hasCategory = (items, cat) =>
  items.some(i => i.category && (i.category === cat || i.category.toLowerCase().includes(cat.toLowerCase())))

describe('useMenuItems — Dynamic Role-Based Menu Generation', () => {
  it('returns superadmin menu items', () => {
    const items = useMenuItems('superadmin')
    expect(items.length).toBeGreaterThan(0)
    expect(hasPath(items, '/admin/audit-logs')).toBe(true)
    expect(hasPath(items, '/admin/monitoring')).toBe(true)
  })

  it('returns admin menu items', () => {
    const items = useMenuItems('admin')
    expect(items.length).toBeGreaterThan(0)
    expect(hasPath(items, '/admin/schools')).toBe(true)
    expect(hasPath(items, '/admin/school-settings')).toBe(true)
  })

  it('returns secretary menu structure with categories', () => {
    const items = useMenuItems('secretary')
    expect(items.length).toBeGreaterThan(0)
    expect(hasCategory(items, 'Anagrafic') || hasCategory(items, 'Classi')).toBe(true)
    expect(hasCategory(items, 'Atti & Certificati')).toBe(true)
    expect(hasCategory(items, 'Servizi & Report')).toBe(true)
  })

  it('returns teacher menu structure with categories and coordinator markers', () => {
    const items = useMenuItems('teacher')
    expect(items.length).toBeGreaterThan(0)
    const valut = items.find(i => i.category === 'Scrutinio & Valutazione' || i.category === 'Didattica & Valutazione')
    expect(valut).toBeDefined()
    expect(hasPath(items, '/teacher/grades')).toBe(true)
    expect(valut.children.some(c => c.coordinatorOnly === true)).toBe(true)
    expect(valut.children.some(c => c.label === 'Coordinamento')).toBe(true)
  })

  it('returns student menu items', () => {
    const items = useMenuItems('student')
    expect(items.length).toBeGreaterThan(0)
    expect(hasPath(items, '/student/grades')).toBe(true)
    expect(hasPath(items, '/student/attendance')).toBe(true)
  })

  it('returns parent menu items', () => {
    const items = useMenuItems('parent')
    expect(items.length).toBeGreaterThan(0)
    expect(hasPath(items, '/parent/children')).toBe(true)
    expect(hasPath(items, '/parent/grades')).toBe(true)
  })

  it('handles principal role mapped to secretary menu', () => {
    const principalItems = useMenuItems('principal')
    expect(principalItems.length).toBeGreaterThan(0)
    expect(hasCategory(principalItems, 'Anagrafic') || hasCategory(principalItems, 'Classi')).toBe(true)
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
    expect(didattica.children.some(c => c.path === '/teacher/grades' && c.label === 'Voti')).toBe(true)
    expect(didattica.children.some(c => c.path === '/teacher/scrutiny' && c.label === 'Scrutinio')).toBe(true)
    expect(didattica.children.some(c => c.label === 'Voti & Scrutinio')).toBe(false)
  })

  it('handles coordinator role mapped to teacher menu', () => {
    const coordItems = useMenuItems('coordinator')
    expect(coordItems.length).toBeGreaterThan(0)
    expect(coordItems.some(i => i.category === 'Scrutinio & Valutazione' || i.category === 'Registro & Classi' || i.category === 'Didattica & Valutazione')).toBe(true)
  })

  it('handles system_auditor mapped to superadmin menu', () => {
    const auditorItems = useMenuItems('system_auditor')
    expect(auditorItems.length).toBeGreaterThan(0)
    expect(auditorItems.some(i => i.path === '/admin/audit-logs' || (i.children && i.children.some(c => c.path === '/admin/audit-logs')))).toBe(true)
  })

  it('handles unknown or empty role gracefully', () => {
    expect(useMenuItems(null)).toEqual([])
    expect(useMenuItems('')).toEqual([])
    expect(useMenuItems('unknown_role_xyz')).toEqual([])
  })
})
