import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import DashboardQuickActions from '@/components/Common/DashboardQuickActions.vue'

const mockPush = vi.fn().mockResolvedValue()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush,
    currentRoute: { value: { path: '/' } }
  })
}))

describe('DashboardQuickActions.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  const testActions = [
    { key: 'attendance', label: 'Segna Presenze', icon: 'how_to_reg', color: 'positive', path: '/teacher/attendance' },
    { key: 'grades', label: 'Inserisci Voto', icon: 'grade', color: 'primary', path: '/teacher/grades', badge: 5 },
    { key: 'substitutions', label: 'Sostituzioni', icon: 'swap_horiz', color: 'warning', path: '/secretary/substitutions', badge: 2, alertOnBadge: true },
    { key: 'customAction', label: 'Azione Personalizzata', icon: 'bolt', action: vi.fn() },
    { key: 'hiddenAction', label: 'Nascosta', icon: 'visibility_off', hidden: true, path: '/hidden' }
  ]

  const createWrapper = (props = {}) => {
    return mount(DashboardQuickActions, {
      props: {
        title: 'Azioni Rapide',
        actions: testActions,
        ...props
      },
      global: {
        stubs: {
          'q-card': {
            template: '<div class="quick-action-card" :class="$attrs.class" tabindex="0"><slot /></div>'
          },
          'q-card-section': {
            template: '<div class="q-card-section"><slot /></div>'
          },
          'q-badge': {
            props: ['label', 'color'],
            template: '<span class="q-badge" :class="color">{{ label }}<slot /></span>'
          },
          'q-icon': {
            template: '<i class="q-icon" />'
          }
        }
      }
    })
  }

  it('renders section title and visible action cards', () => {
    const wrapper = createWrapper()
    expect(wrapper.text()).toContain('Azioni Rapide')
    expect(wrapper.text()).toContain('Segna Presenze')
    expect(wrapper.text()).toContain('Inserisci Voto')
    expect(wrapper.text()).toContain('Sostituzioni')
    expect(wrapper.text()).not.toContain('Nascosta')
  })

  it('displays badge and applies alert class when alertOnBadge is true', () => {
    const wrapper = createWrapper()
    const badges = wrapper.findAll('.q-badge')
    expect(badges.length).toBeGreaterThanOrEqual(2)
    expect(wrapper.text()).toContain('5')
    expect(wrapper.text()).toContain('2')

    const alertCard = wrapper.find('.quick-action-card--alert')
    expect(alertCard.exists()).toBe(true)
    expect(alertCard.text()).toContain('Sostituzioni')
  })

  it('triggers router.push and emits action event on card click', async () => {
    const wrapper = createWrapper()
    const attendanceCard = wrapper.findAll('.quick-action-card')[0]

    await attendanceCard.trigger('click')

    expect(mockPush).toHaveBeenCalledWith('/teacher/attendance')
    expect(wrapper.emitted('action')).toBeTruthy()
    expect(wrapper.emitted('action')[0][0].key).toBe('attendance')
  })

  it('executes action callback if provided instead of navigating', async () => {
    const customCallback = vi.fn()
    const actions = [
      { key: 'custom', label: 'Bozza', icon: 'edit', action: customCallback }
    ]
    const wrapper = createWrapper({ actions })
    const card = wrapper.find('.quick-action-card')

    await card.trigger('click')

    expect(customCallback).toHaveBeenCalledTimes(1)
    expect(mockPush).not.toHaveBeenCalled()
    expect(wrapper.emitted('action')).toBeTruthy()
  })

  it('supports keyboard enter/space activation for accessibility', async () => {
    const wrapper = createWrapper()
    const card = wrapper.findAll('.quick-action-card')[0]

    await card.trigger('keydown.enter')
    expect(mockPush).toHaveBeenCalledWith('/teacher/attendance')

    await card.trigger('keydown.space')
    expect(mockPush).toHaveBeenCalledTimes(2)
  })

  it('respects maxVisible prop', () => {
    const wrapper = createWrapper({ maxVisible: 2 })
    const cards = wrapper.findAll('.quick-action-card')
    expect(cards).toHaveLength(2)
  })
})
