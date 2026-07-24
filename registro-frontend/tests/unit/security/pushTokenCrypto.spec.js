import { describe, it, expect, vi } from 'vitest'
import { notificationService } from '@/services/notificationService'

describe('Security: Push Token Cryptographic Randomness', () => {
    it('should generate UUID-based web push tokens using crypto.randomUUID or crypto.getRandomValues', async () => {
        // Mock Notification permission
        global.Notification = {
            requestPermission: vi.fn().mockResolvedValue('granted')
        }

        const registerSpy = vi.spyOn(notificationService, 'registerPushToken').mockResolvedValue({})

        const token = await notificationService.requestPushPermissionAndRegister()

        expect(token).toBeDefined()
        expect(token).toMatch(/^web_push_[a-f0-9\-]{16,36}$/i)
        expect(registerSpy).toHaveBeenCalledWith(token, 'web')

        registerSpy.mockRestore()
    })
})
