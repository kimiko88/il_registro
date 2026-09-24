import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import ActiveStrikeNoticeBanner from '@/components/Common/ActiveStrikeNoticeBanner.vue'
import strikeService from '@/services/strikeService'

vi.mock('@/services/strikeService', () => ({
  default: {
    getStrikeNotices: vi.fn(),
    submitDeclaration: vi.fn()
  }
}))

describe('ActiveStrikeNoticeBanner.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  const createWrapper = async (role = 'teacher', mockNotices = []) => {
    strikeService.getStrikeNotices.mockResolvedValue(mockNotices)

    const wrapper = mount(ActiveStrikeNoticeBanner, {
      global: {
        plugins: [
          Quasar,
          createTestingPinia({
            initialState: {
              auth: {
                user: { id: 'user-1', name: 'Mario Rossi', role },
                token: null
              }
            },
            stubActions: false
          })
        ],
        mocks: {
          t: (key) => key
        },
        stubs: {
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-badge': { template: '<span class="q-badge"><slot /></span>' },
          'q-chip': { template: '<span class="q-chip"><slot /></span>' },
          'q-btn': { template: '<button class="q-btn"><slot /></button>' },
          'q-icon': true
        }
      }
    })

    await flushPromises()
    return wrapper
  }

  it('1. does not render anything when user is not staff (e.g. student)', async () => {
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 2)

    const wrapper = await createWrapper('student', [
      {
        id: 'notice-1',
        title: 'Sciopero Generale',
        strike_date: tomorrow.toISOString(),
        declaration_deadline: tomorrow.toISOString(),
        is_expired: false
      }
    ])
    expect(wrapper.find('.strike-banners-container').exists()).toBe(false)
  })

  it('2. renders banner for staff user with active strike notices', async () => {
    const futureDate = new Date()
    futureDate.setDate(futureDate.getDate() + 3)

    const mockNotices = [
      {
        id: 'notice-1',
        title: 'Sciopero Comparto Scuola',
        proclaimed_by: 'FLC CGIL',
        strike_date: futureDate.toISOString(),
        declaration_deadline: futureDate.toISOString(),
        is_expired: false,
        user_declaration: null
      }
    ]

    const wrapper = await createWrapper('teacher', mockNotices)
    expect(wrapper.find('.strike-banners-container').exists()).toBe(true)
    expect(wrapper.text()).toContain('Sciopero Comparto Scuola')
    expect(wrapper.text()).toContain('FLC CGIL')
  })

  it('3. submits voluntary intention declaration when action button clicked', async () => {
    const futureDate = new Date()
    futureDate.setDate(futureDate.getDate() + 2)

    const mockNotices = [
      {
        id: 'notice-123',
        title: 'Sciopero Nazionale',
        strike_date: futureDate.toISOString(),
        declaration_deadline: futureDate.toISOString(),
        is_expired: false,
        user_declaration: null
      }
    ]

    strikeService.submitDeclaration.mockResolvedValue({
      id: 'decl-1',
      intention: 'participates'
    })

    const wrapper = await createWrapper('teacher', mockNotices)

    const notice = wrapper.vm.notices[0]
    expect(notice).toBeDefined()
    await wrapper.vm.setDeclaration(notice, 'participates')

    expect(strikeService.submitDeclaration).toHaveBeenCalledWith('notice-123', 'participates')
    expect(notice.user_declaration).toEqual({ id: 'decl-1', intention: 'participates' })
  })

  it('4. displays expired state when notice deadline is locked', async () => {
    const futureDate = new Date()
    futureDate.setDate(futureDate.getDate() + 1)

    const mockNotices = [
      {
        id: 'notice-exp',
        title: 'Sciopero con Termine Scaduto',
        strike_date: futureDate.toISOString(),
        declaration_deadline: '2026-01-01T00:00:00Z',
        is_expired: true,
        user_declaration: { intention: 'not_participates' }
      }
    ]

    const wrapper = await createWrapper('collaboratore_scolastico', mockNotices)
    expect(wrapper.text()).toContain('Sciopero con Termine Scaduto')
    expect(wrapper.vm.currentUserIntention(wrapper.vm.notices[0])).toBe('not_participates')
  })
})
