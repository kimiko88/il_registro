import { describe, it, expect, vi, beforeEach } from 'vitest'
import { defineComponent } from 'vue'
import { mount } from '@vue/test-utils'
import { useMenuBadges } from '@/composables/useMenuBadges'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

describe('useMenuBadges composable', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('fetches notification count and role-specific badges for teacher', async () => {
    api.get.mockImplementation((url) => {
      if (url === '/notifications/unread-count') return Promise.resolve({ data: { count: 3 } })
      if (url === '/teacher/grades/pending-count') return Promise.resolve({ data: { count: 7 } })
      if (url === '/teacher/colloqui/pending-count') return Promise.resolve({ data: { count: 2 } })
      return Promise.resolve({ data: { count: 0 } })
    })

    const TestComponent = defineComponent({
      setup() {
        const { badges, getBadge } = useMenuBadges('teacher')
        return { badges, getBadge }
      },
      template: '<div></div>'
    })

    const wrapper = mount(TestComponent)
    await vi.waitFor(() => {
      expect(wrapper.vm.badges.unreadMessages.value).toBe(3)
    })

    expect(wrapper.vm.badges.pendingGrades.value).toBe(7)
    expect(wrapper.vm.badges.pendingColloqui.value).toBe(2)
    expect(wrapper.vm.getBadge('pendingGrades')).toBe(7)
    expect(wrapper.vm.getBadge('nonExistentKey')).toBeUndefined()
    expect(wrapper.vm.getBadge('absentStaff')).toBeUndefined() // is 0, so returns undefined

    wrapper.unmount()
  })

  it('fetches substitution and staff counts for secretary/admin', async () => {
    api.get.mockImplementation((url) => {
      if (url === '/notifications/unread-count') return Promise.resolve({ data: { count: 1 } })
      if (url === '/ata/substitutions/pending-count') return Promise.resolve({ data: { count: 4 } })
      if (url === '/ata/attendance/absent-count') return Promise.resolve({ data: { count: 6 } })
      return Promise.resolve({ data: { count: 0 } })
    })

    const TestComponent = defineComponent({
      setup() {
        const { badges, getBadge } = useMenuBadges('secretary')
        return { badges, getBadge }
      },
      template: '<div></div>'
    })

    const wrapper = mount(TestComponent)
    await vi.waitFor(() => {
      expect(wrapper.vm.badges.pendingSubstitutions.value).toBe(4)
    })

    expect(wrapper.vm.badges.absentStaff.value).toBe(6)
    expect(wrapper.vm.getBadge('pendingSubstitutions')).toBe(4)
    expect(wrapper.vm.getBadge('absentStaff')).toBe(6)

    wrapper.unmount()
  })

  it('silently handles network/API errors without throwing', async () => {
    api.get.mockRejectedValue(new Error('Network error'))

    const TestComponent = defineComponent({
      setup() {
        const { badges, getBadge } = useMenuBadges('student')
        return { badges, getBadge }
      },
      template: '<div></div>'
    })

    expect(() => {
      const wrapper = mount(TestComponent)
      wrapper.unmount()
    }).not.toThrow()
  })
})
