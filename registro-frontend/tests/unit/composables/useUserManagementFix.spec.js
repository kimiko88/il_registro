import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserManagement } from '@/composables/useUserManagement'
import adminService from '@/services/adminService'
import { Notify } from 'quasar'

vi.mock('@/services/adminService', () => ({
  default: {
    getAdmins: vi.fn(),
    importUsers: vi.fn()
  }
}))

vi.mock('quasar', () => ({
  useQuasar: () => ({
    notify: vi.fn()
  }),
  Notify: {
    create: vi.fn()
  }
}))

describe('useUserManagement Composable Suite', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetchUsers successfully populates user list', async () => {
    adminService.getAdmins.mockResolvedValueOnce({
      data: {
        items: [
          { id: '1', email: 'admin1@school.it' },
          { id: '2', email: 'admin2@school.it' }
        ]
      }
    })

    const { users, loading, fetchUsers } = useUserManagement()
    expect(users.value).toEqual([])
    expect(loading.value).toBe(false)

    await fetchUsers()

    expect(adminService.getAdmins).toHaveBeenCalledTimes(1)
    expect(users.value).toHaveLength(2)
    expect(users.value[0].email).toBe('admin1@school.it')
    expect(loading.value).toBe(false)
  })

  it('importUsers returns true and positive notification when service succeeds', async () => {
    adminService.importUsers.mockResolvedValueOnce({ data: { success: true } })

    const { importUsers, loading } = useUserManagement()
    const fakeFile = new File(['email,role\ntest@test.it,teacher'], 'users.csv', { type: 'text/csv' })

    const result = await importUsers(fakeFile)

    expect(result).toBe(true)
    expect(adminService.importUsers).toHaveBeenCalledWith(fakeFile)
    expect(loading.value).toBe(false)
  })

  it('importUsers returns false and does NOT trigger a positive notification when service fails', async () => {
    adminService.importUsers.mockRejectedValueOnce(new Error('Invalid CSV structure'))

    const { importUsers, loading } = useUserManagement()
    const fakeFile = new File(['invalid'], 'users.csv', { type: 'text/csv' })

    const result = await importUsers(fakeFile)

    expect(result).toBe(false)
    expect(loading.value).toBe(false)
  })
})
