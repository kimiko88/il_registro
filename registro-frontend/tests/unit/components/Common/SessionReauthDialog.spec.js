import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import SessionReauthDialog from '@/components/Common/SessionReauthDialog.vue'
import { useSessionReauth } from '@/composables/useSessionReauth'
import authService from '@/services/authService'
import { useAuthStore } from '@/stores/auth'

const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  })
}))


vi.mock('@/services/authService', () => ({
  default: {
    login: vi.fn()
  }
}))

describe('SessionReauthDialog.vue', () => {
  let reauth

  beforeEach(() => {
    vi.clearAllMocks()
    reauth = useSessionReauth()
    reauth.showDialog.value = false
  })

  it('renders reauth dialog content when triggered', async () => {
    const pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { token: null, user: { email: 'prof.test@scuola.it' } }
      }
    })

    const wrapper = mount(SessionReauthDialog, {
      global: {
        plugins: [pinia],
        stubs: {
          'q-dialog': {
            props: ['modelValue'],
            template: '<div v-if="modelValue" class="q-dialog"><slot /></div>'
          },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-input': {
            props: ['modelValue', 'label', 'type', 'errorMessage'],
            methods: {
              focus() {}
            },
            template: `
              <div class="q-input">
                <label>{{ label }}</label>
                <span v-if="modelValue" class="input-val">{{ modelValue }}</span>
                <input :value="modelValue" :type="type" @input="$emit('update:modelValue', $event.target.value)" />
                <span v-if="errorMessage" class="error-msg">{{ errorMessage }}</span>
              </div>
            `
          },
          'q-btn': {
            props: ['label'],
            template: '<button class="q-btn" type="button">{{ label }}<slot /></button>'
          },
          'q-icon': { template: '<i class="q-icon" />' }
        }
      }
    })

    // Initially not visible
    expect(wrapper.find('.q-dialog').exists()).toBe(false)

    // Trigger reauth with user email
    reauth.triggerReauth('prof.test@scuola.it')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.q-dialog').exists()).toBe(true)
    expect(wrapper.text()).toContain('prof.test@scuola.it')
  })

  it('successfully logs in and resolves reauth promise on submit', async () => {
    const pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { token: null, user: { email: 'prof.test@scuola.it' } }
      }
    })

    authService.login.mockResolvedValueOnce({
      access_token: 'new-valid-jwt-token',
      user: { id: 'u1', email: 'prof.test@scuola.it', role: 'teacher' }
    })

    const wrapper = mount(SessionReauthDialog, {
      global: {
        plugins: [pinia],
        stubs: {
          'q-dialog': {
            props: ['modelValue'],
            template: '<div v-if="modelValue" class="q-dialog"><slot /></div>'
          },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-input': {
            props: ['modelValue', 'label', 'type'],
            methods: { focus() {} },
            template: `
              <div class="q-input">
                <input :value="modelValue" :type="type" @input="$emit('update:modelValue', $event.target.value)" />
              </div>
            `
          },
          'q-btn': {
            props: ['label'],
            template: '<button class="q-btn" type="button">{{ label }}<slot /></button>'
          },
          'q-icon': { template: '<i class="q-icon" />' }
        }
      }
    })

    const reauthPromise = reauth.triggerReauth('prof.test@scuola.it')
    await wrapper.vm.$nextTick()

    // Type password
    const passwordInput = wrapper.findAll('input')[1]
    await passwordInput.setValue('CorrectSecretPass123!')

    // Find and click the login button
    const buttons = wrapper.findAll('.q-btn')
    const loginBtn = buttons[buttons.length - 1] // second button is Accedi / login
    await loginBtn.trigger('click')

    expect(authService.login).toHaveBeenCalledWith('prof.test@scuola.it', 'CorrectSecretPass123!')

    // Reauth promise should resolve with the new token
    const resolvedToken = await reauthPromise
    expect(resolvedToken).toBe('new-valid-jwt-token')

    // Dialog should be closed
    expect(reauth.showDialog.value).toBe(false)
  })

  it('cancels reauth, logs out, and redirects to /login on cancel click', async () => {
    const pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: { token: null, user: { email: 'prof.test@scuola.it' } }
      }
    })

    const wrapper = mount(SessionReauthDialog, {
      global: {
        plugins: [pinia],
        stubs: {
          'q-dialog': {
            props: ['modelValue'],
            template: '<div v-if="modelValue" class="q-dialog"><slot /></div>'
          },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-input': {
            methods: { focus() {} },
            template: '<div class="q-input"><input /></div>'
          },
          'q-btn': {
            props: ['label'],
            template: '<button class="q-btn" type="button">{{ label }}<slot /></button>'
          },
          'q-icon': { template: '<i class="q-icon" />' }
        }
      }
    })

    let rejectedError = null
    const reauthPromise = reauth.triggerReauth('prof.test@scuola.it').catch((err) => {
      rejectedError = err
      return err
    })
    await wrapper.vm.$nextTick()

    // Click cancel button (first button in actions: "Esci")
    const cancelBtn = wrapper.find('.q-btn')
    await cancelBtn.trigger('click')

    await reauthPromise
    expect(rejectedError).not.toBeNull()
    expect(rejectedError.message).toBe('reauth_cancelled')
    const authStore = useAuthStore()
    expect(authStore.logout).toHaveBeenCalled()
    expect(mockPush).toHaveBeenCalledWith('/login')
  })
})
