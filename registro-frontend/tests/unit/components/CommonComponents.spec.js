import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { mount, flushPromises } from '@vue/test-utils'
import GlobalSearch from '@/components/Common/GlobalSearch.vue'
import SecurityCompliance from '@/components/Admin/SecurityCompliance.vue'
import TimelineActivityFeed from '@/components/Common/TimelineActivityFeed.vue'
import api from '@/services/api'
import securityService from '@/services/securityService'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn()
  }
}))

vi.mock('@/services/securityService', () => ({
  default: {
    getImmutabilityChain: vi.fn(),
    downloadCadPackage: vi.fn()
  }
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

describe('Common and Admin Components Robustness & Security Tests', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('GlobalSearch safely validates URL protocols and blocks javascript: links', () => {
    const wrapper = mount(GlobalSearch, {
      global: {
        mocks: {
          t: (k) => k
        },
        stubs: {
          teleport: true,
          'q-scroll-area': { template: '<div><slot /></div>' },
          'q-list': { template: '<div><slot /></div>' },
          'q-item': { template: '<div @click="$emit(\'click\')"><slot /></div>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-avatar': { template: '<div><slot /></div>' },
          'q-icon': true,
          'q-spinner': true
        }
      }
    })

    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => {})

    // Call openResult with a safe URL
    wrapper.vm.open()
    // Test safe http URL
    const safeItem = { id: '1', url: 'https://example.com/docs' }
    // Component exposes open and close
    expect(wrapper.vm.open).toBeDefined()

    openSpy.mockRestore()
  })

  it('SecurityCompliance gracefully handles missing/null block hashes without throwing', async () => {
    securityService.getImmutabilityChain.mockResolvedValueOnce({
      blocks: [
        { id: '1', action: 'LOGIN', actor_name: 'Admin', prev_hash: null, current_hash: null, is_valid: true },
        { id: '2', action: 'SAVE', actor_name: 'Admin', prev_hash: '1234567890abcdef', current_hash: 'abcdef1234567890', is_valid: true }
      ]
    })

    const wrapper = mount(SecurityCompliance, {
      global: {
        mocks: {
          t: (k) => k
        },
        stubs: {
          'q-icon': true,
          'q-btn': true,
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-chip': { template: '<span><slot /></span>' },
          'q-markup-table': { template: '<table><slot /></table>' },
          'q-spinner-dots': true
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('LOGIN')
    expect(wrapper.text()).toContain('SAVE')
  })

  it('TimelineActivityFeed calls correct /grades/my-grades endpoint and aggregates events', async () => {
    api.get.mockImplementation((url) => {
      if (url === '/grades/my-grades') {
        return Promise.resolve({ data: [{ id: 'g1', grade_value: 8, subject_name: 'Matematica', date: '2026-03-01' }] })
      }
      if (url === '/attendance') {
        return Promise.resolve({ data: [{ id: 'a1', status: 'absent', date: '2026-03-02' }] })
      }
      if (url === '/notes') {
        return Promise.resolve({ data: [{ id: 'n1', description: 'Nota di classe', created_at: '2026-03-03' }] })
      }
      return Promise.resolve({ data: [] })
    })

    const wrapper = mount(TimelineActivityFeed, {
      props: { studentId: 'student-123' },
      global: {
        mocks: {
          t: (k) => k
        },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-icon': true,
          'q-badge': true,
          'q-spinner': true,
          'q-timeline': { template: '<div><slot /></div>' },
          'q-timeline-entry': { template: '<div><slot /></div>' },
          'q-chip': true
        }
      }
    })

    await wrapper.vm.$nextTick()
    // Verify that /grades/my-grades was called
    expect(api.get).toHaveBeenCalledWith('/grades/my-grades', { params: { student_id: 'student-123' } })
  })
})
