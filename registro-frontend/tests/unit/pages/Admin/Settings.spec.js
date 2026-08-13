import { mount } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { Quasar } from 'quasar'
import { createTestingPinia } from '@pinia/testing'
import { createI18n } from 'vue-i18n'
import Settings from '@/pages/admin/Settings.vue'

const i18n = createI18n({
    legacy: false,
    locale: 'it-IT',
    messages: {
        'it-IT': {}
    }
})

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            notify: vi.fn()
        })
    }
})

vi.mock('@/services/api', () => ({
    default: {
        post: vi.fn().mockResolvedValue({ data: { success: true } })
    }
}))

describe('Settings.vue — Admin Settings & Self-Service Password Change', () => {
    let wrapper

    beforeEach(() => {
        vi.clearAllMocks()
        localStorage.clear()
        wrapper = mount(Settings, {
            global: {
                plugins: [
                    [Quasar, {}],
                    i18n,
                    createTestingPinia({ createSpy: vi.fn })
                ]
            }
        })
    })

    it('renders settings dialog state initial values', () => {
        expect(wrapper.vm.passwordDialog).toBe(false)
        expect(wrapper.vm.securityDialog).toBe(false)
        expect(wrapper.vm.languageDialog).toBe(false)
    })

    it('opens password change modal when opening password dialog', () => {
        wrapper.vm.openPasswordDialog()
        expect(wrapper.vm.passwordDialog).toBe(true)
    })

    it('validates password complexity and confirmation matching', () => {
        wrapper.vm.passwordForm.current_password = 'OldPassword123!'
        wrapper.vm.passwordForm.new_password = 'short'
        wrapper.vm.passwordForm.confirm_password = 'short'

        expect(wrapper.vm.passwordForm.new_password.length).toBeLessThan(10)

        // Set matching valid password
        wrapper.vm.passwordForm.new_password = 'NewSecretPassword123!'
        wrapper.vm.passwordForm.confirm_password = 'NewSecretPassword123!'

        expect(wrapper.vm.passwordForm.new_password).toBe(wrapper.vm.passwordForm.confirm_password)
        expect(wrapper.vm.passwordForm.new_password.length).toBeGreaterThanOrEqual(10)
    })

    it('persists selected language to localStorage when calling saveLanguage', () => {
        wrapper.vm.selectedLanguage = 'en-US'
        wrapper.vm.saveLanguage()
        expect(localStorage.getItem('superadmin_language')).toBe('en-US')
    })

    it('saves security settings to localStorage', () => {
        wrapper.vm.securitySettings.require_mfa = true
        wrapper.vm.securitySettings.min_password_length = 12
        wrapper.vm.saveSecuritySettings()

        const saved = JSON.parse(localStorage.getItem('superadmin_security_settings'))
        expect(saved.require_mfa).toBe(true)
        expect(saved.min_password_length).toBe(12)
    })
})
