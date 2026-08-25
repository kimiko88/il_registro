import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuditLogStore } from '@/stores/auditLog'
import { useCertificatesStore } from '@/stores/certificates'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn(),
  }
}))

describe('AuditLog and Certificates Stores — Params & Blob Handling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('useAuditLogStore', () => {
    it('passes filter parameters cleanly in params config object', async () => {
      api.get.mockResolvedValueOnce({
        data: {
          data: [{ id: 1, action: 'LOGIN', actor_id: 'u1' }],
          total: 1,
          page: 1,
          limit: 20,
          total_pages: 1
        }
      })

      const store = useAuditLogStore()
      const res = await store.fetchLogs({ actor_id: 'u1', action: 'LOGIN' }, 1, 20)

      expect(api.get).toHaveBeenCalledWith('/audit-log', {
        params: {
          page: 1,
          limit: 20,
          actor_id: 'u1',
          action: 'LOGIN'
        }
      })
      expect(store.logs).toHaveLength(1)
      expect(store.total).toBe(1)
    })
  })

  describe('useCertificatesStore', () => {
    it('passes academic year and filter parameters in params object', async () => {
      api.get.mockResolvedValueOnce({
        data: [{ id: 'cert1', student_id: 's1', academic_year: '2025/2026' }]
      })

      const store = useCertificatesStore()
      const res = await store.fetchCertificates({
        student_id: 's1',
        academic_year: '2025/2026',
        type: 'enrollment'
      })

      expect(api.get).toHaveBeenCalledWith('/certificates', {
        params: {
          student_id: 's1',
          academic_year: '2025/2026',
          type: 'enrollment'
        }
      })
      expect(res).toHaveLength(1)
      expect(store.certificates).toHaveLength(1)
    })

    it('generates a certificate and prepends it to certificates array', async () => {
      api.post.mockResolvedValueOnce({
        data: { cert: { id: 'cert2', title: 'Certificato di Frequenza' } }
      })

      const store = useCertificatesStore()
      await store.generateCertificate({ type: 'attendance' })

      expect(api.post).toHaveBeenCalledWith('/certificates/generate', { type: 'attendance' })
      expect(store.certificates[0].id).toBe('cert2')
    })
  })
})
