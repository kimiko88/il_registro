import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSchoolYearStore } from '../../src/stores/schoolYear'
import { useMenuItems } from '../../src/composables/useMenuItems'
import authService from '../../src/services/authService'
import api from '../../src/services/api'

vi.mock('../../src/services/api', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn()
  }
}))

describe('Teacher Settings & School Year Logic', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  it('generates available school years from user registration year up to current year', () => {
    const store = useSchoolYearStore()
    
    // User registered in October 2023 -> 2023/2024 academic year
    store.initializeForUser({ created_at: '2023-10-15T10:00:00Z' })
    
    expect(store.availableSchoolYears).toBeDefined()
    expect(store.availableSchoolYears.length).toBeGreaterThanOrEqual(1)
    expect(store.availableSchoolYears).toContain('2023/2024')
    expect(store.availableSchoolYears[0]).toBe('2025/2026') // current active school year
  })

  it('includes Impostazioni in teacher menu items', () => {
    const items = useMenuItems('teacher')
    const comCat = items.find(cat => cat.category === 'Comunicazioni & Atti')
    expect(comCat).toBeDefined()
    const settingsItem = comCat.children.find(child => child.path === '/teacher/settings')
    expect(settingsItem).toBeDefined()
    expect(settingsItem.label).toBe('Impostazioni')
  })

  it('calls changePassword endpoint via authService', async () => {
    api.post.mockResolvedValueOnce({ data: { message: 'password changed successfully' } })

    const res = await authService.changePassword('OldPass123!', 'NewPass123!')
    expect(api.post).toHaveBeenCalledWith('/auth/change-password', {
      current_password: 'OldPass123!',
      new_password: 'NewPass123!'
    })
    expect(res.message).toBe('password changed successfully')
  })

  it('sets and persists selected school year in localStorage', () => {
    const store = useSchoolYearStore()
    store.setSchoolYear('2024/2025')
    
    expect(store.selectedSchoolYear).toBe('2024/2025')
    expect(localStorage.getItem('selected_school_year')).toBe('2024/2025')
  })
})
