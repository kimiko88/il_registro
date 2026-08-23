import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import TimelineActivityFeed from '@/components/Common/TimelineActivityFeed.vue'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

describe('TimelineActivityFeed Component Suite', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders timeline with combined grades, absences, and notes', async () => {
    api.get.mockImplementation((url) => {
      if (url === '/grades/my-grades') {
        return Promise.resolve({
          data: [
            { id: 101, subject_name: 'Matematica', date: '2026-05-10', grade_value: 9, description: 'Verifica scritta' }
          ]
        })
      }
      if (url === '/attendance') {
        return Promise.resolve({
          data: [
            { id: 201, status: 'absent', date: '2026-05-09', justification_reason: 'Febbre' }
          ]
        })
      }
      if (url === '/notes') {
        return Promise.resolve({
          data: [
            { id: 301, note_text: 'Dimenticato libro di testo', created_at: '2026-05-08' }
          ]
        })
      }
      return Promise.resolve({ data: [] })
    })

    const wrapper = mount(TimelineActivityFeed, {
      props: {
        studentId: 'stud-1'
      },
      global: {
        mocks: { t: (k) => k, locale: { value: 'it-IT' } },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-spinner': true,
          'q-icon': true,
          'q-timeline': { template: '<div class="timeline"><slot /></div>' },
          'q-timeline-entry': {
            props: ['title', 'subtitle', 'color', 'icon'],
            template: '<div class="timeline-entry" :data-title="title"><slot /></div>'
          },
          'q-chip': { template: '<span class="chip"><slot /></span>' }
        }
      }
    })

    await flushPromises()

    const entries = wrapper.findAll('.timeline-entry')
    expect(entries).toHaveLength(3)
    expect(wrapper.text()).toContain('Verifica scritta')
    expect(wrapper.text()).toContain('Febbre')
    expect(wrapper.text()).toContain('Dimenticato libro di testo')
  })

  it('renders empty placeholder when all streams return empty data', async () => {
    api.get.mockResolvedValue({ data: [] })

    const wrapper = mount(TimelineActivityFeed, {
      global: {
        mocks: { t: (k) => k, locale: { value: 'it-IT' } },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-badge': true,
          'q-spinner': true,
          'q-icon': true,
          'q-timeline': true
        }
      }
    })

    await flushPromises()

    expect(wrapper.text()).toContain('dashboardPage.noLessons')
  })

  it('resiliently handles network error on one or all streams without throwing', async () => {
    api.get.mockRejectedValue(new Error('Network error'))

    const wrapper = mount(TimelineActivityFeed, {
      global: {
        mocks: { t: (k) => k, locale: { value: 'it-IT' } },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-badge': true,
          'q-spinner': true,
          'q-icon': true,
          'q-timeline': true
        }
      }
    })

    await flushPromises()

    expect(wrapper.text()).toContain('dashboardPage.noLessons')
  })
})
