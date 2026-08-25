import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../../../src/stores/auth'
import { useWebSocketStore } from '../../../src/stores/websocket'

describe('Security & Logic Audit Unit Tests', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()
        vi.restoreAllMocks()
    })

    it('authStore sanitizeUserData uses strict allowlist without sensitive PII', () => {
        const authStore = useAuthStore()

        const inputUser = {
            id: 'usr-123',
            email: 'teacher@school.it',
            first_name: 'Mario',
            last_name: 'Rossi',
            role: 'teacher',
            school_id: 'school-456',
            class_id: 'class-789',
            is_staff: true,
            password_hash: '$2a$12$secret_hash',
            mfa_secret: 'SUPER_SECRET_MFA',
            ssn: 'RSSMRA80A01H501U',
            tax_id: '12345678901'
        }

        authStore.updateUser(inputUser)

        expect(authStore.user).toEqual({
            id: 'usr-123',
            email: 'teacher@school.it',
            first_name: 'Mario',
            last_name: 'Rossi',
            role: 'teacher',
            school_id: 'school-456',
            class_id: 'class-789',
            is_staff: true
        })
        expect(authStore.user.password_hash).toBeUndefined()
        expect(authStore.user.mfa_secret).toBeUndefined()
        expect(authStore.user.ssn).toBeUndefined()
    })

    it('authStore initializes with isInitializing flag available for router guard', () => {
        const authStore = useAuthStore()
        expect(authStore.isInitializing).toBeDefined()
        expect(authStore.isInitializing).toBe(false)
    })

    it('websocketStore handles reconnect and disconnect states safely', () => {
        const wsStore = useWebSocketStore()
        expect(wsStore.isConnected).toBe(false)
        wsStore.disconnect(true)
        expect(wsStore.hasFailedPermanently).toBe(false)
    })
})
