import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import { Quasar } from 'quasar'
import InactivityDialog from '@/components/Common/InactivityDialog.vue'

const mockShowWarning = ref(false)
const mockSecondsRemaining = ref(60)
const mockStayLoggedIn = vi.fn()
const mockLogoutNow = vi.fn()

vi.mock('@/composables/useInactivityTimer', () => ({
  useInactivityTimer: () => ({
    showWarningDialog: mockShowWarning,
    secondsRemaining: mockSecondsRemaining,
    stayLoggedIn: mockStayLoggedIn,
    logoutNow: mockLogoutNow
  })
}))

describe('InactivityDialog.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockShowWarning.value = false
    mockSecondsRemaining.value = 60
  })

  const createWrapper = () => {
    return mount(InactivityDialog, {
      global: {
        plugins: [Quasar],
        mocks: {
          t: (k) => k
        },
        stubs: {
          'q-dialog': {
            props: ['modelValue'],
            template: '<div v-if="modelValue" class="q-dialog" role="alertdialog" aria-modal="true"><slot /></div>'
          },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-circular-progress': {
            props: ['value'],
            template: '<div class="q-circular-progress" :data-value="value"><slot /></div>'
          },
          'q-btn': {
            props: ['label'],
            template: '<button class="q-btn">{{ label }}<slot /></button>'
          },
          'q-icon': true
        }
      }
    })
  }

  it('1. does not render dialog content when warning flag is false', () => {
    mockShowWarning.value = false
    const wrapper = createWrapper()
    expect(wrapper.find('.q-dialog').exists()).toBe(false)
  })

  it('2. renders alert dialog with countdown when showWarningDialog is true', async () => {
    mockShowWarning.value = true
    mockSecondsRemaining.value = 45
    const wrapper = createWrapper()
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.q-dialog').exists()).toBe(true)
    expect(wrapper.find('[role="alertdialog"]').exists()).toBe(true)
    expect(wrapper.find('[aria-modal="true"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('45s')
    expect(wrapper.text()).toContain('inactivity.warningTitle')
  })

  it('3. calculates progress value accurately based on seconds remaining', async () => {
    mockShowWarning.value = true
    mockSecondsRemaining.value = 30 // 30s out of 60s is 50%
    const wrapper = createWrapper()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.progressValue).toBe(50)
  })

  it('4. calls stayLoggedIn when clicking stay connected button', async () => {
    mockShowWarning.value = true
    const wrapper = createWrapper()
    await wrapper.vm.$nextTick()

    // Find all buttons
    const buttons = wrapper.findAll('.q-btn')
    expect(buttons.length).toBeGreaterThanOrEqual(2)

    // Second button is stayLoggedIn
    const stayBtn = buttons[1]
    await stayBtn.trigger('click')

    expect(mockStayLoggedIn).toHaveBeenCalledTimes(1)
  })

  it('5. calls logoutNow when clicking logout button', async () => {
    mockShowWarning.value = true
    const wrapper = createWrapper()
    await wrapper.vm.$nextTick()

    const logoutBtn = wrapper.findAll('.q-btn')[0]
    await logoutBtn.trigger('click')

    expect(mockLogoutNow).toHaveBeenCalledTimes(1)
  })
})
