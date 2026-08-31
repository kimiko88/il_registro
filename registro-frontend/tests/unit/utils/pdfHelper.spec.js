import { describe, it, expect, vi, beforeEach } from 'vitest'
import { handleAsyncPdfDownload } from '@/utils/pdfHelper'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn()
    }
}))

describe('pdfHelper - handleAsyncPdfDownload', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        window.URL.createObjectURL = vi.fn(() => 'blob:http://localhost/mock-uuid')
        window.URL.revokeObjectURL = vi.fn()
    })

    it('downloads direct PDF blob synchronously', async () => {
        const mockBlob = new Blob(['PDF data'], { type: 'application/pdf' })
        const fetchFn = vi.fn().mockResolvedValue({
            status: 200,
            headers: { 'content-type': 'application/pdf' },
            data: mockBlob
        })

        const result = await handleAsyncPdfDownload(fetchFn, 'test.pdf')
        expect(result).toBe(true)
        expect(fetchFn).toHaveBeenCalledTimes(1)
        expect(window.URL.createObjectURL).toHaveBeenCalledWith(expect.any(Blob))
    })

    it('handles async 202 Accepted response and polls until completed', async () => {
        const jobResponse = {
            job_id: 'job-123',
            status: 'pending',
            status_url: '/api/v1/scrutiny/pdf-jobs/job-123'
        }
        const jsonBlob = new Blob([JSON.stringify(jobResponse)], { type: 'application/json' })
        const fetchFn = vi.fn().mockResolvedValue({
            status: 202,
            headers: { 'content-type': 'application/json' },
            data: jsonBlob
        })

        const pdfBlob = new Blob(['Real PDF Content'], { type: 'application/pdf' })

        api.get.mockImplementation(url => {
            if (url === '/api/v1/scrutiny/pdf-jobs/job-123') {
                return Promise.resolve({
                    data: {
                        job_id: 'job-123',
                        status: 'completed',
                        download_url: '/api/v1/scrutiny/export/s1/pdf'
                    }
                })
            }
            if (url.includes('/api/v1/scrutiny/export/s1/pdf')) {
                return Promise.resolve({
                    status: 200,
                    data: pdfBlob
                })
            }
            return Promise.reject(new Error('Unknown URL: ' + url))
        })

        const result = await handleAsyncPdfDownload(fetchFn, 'pagella.pdf', { pollIntervalMs: 10, maxPollAttempts: 5 })
        expect(result).toBe(true)
        expect(api.get).toHaveBeenCalledWith('/api/v1/scrutiny/pdf-jobs/job-123')
        expect(api.get).toHaveBeenCalledWith('/api/v1/scrutiny/export/s1/pdf?sync=true', { responseType: 'blob', timeout: 60000 })
    })

    it('throws error when job fails', async () => {
        const jobResponse = {
            job_id: 'job-err',
            status: 'pending',
            status_url: '/api/v1/scrutiny/pdf-jobs/job-err'
        }
        const jsonBlob = new Blob([JSON.stringify(jobResponse)], { type: 'application/json' })
        const fetchFn = vi.fn().mockResolvedValue({
            status: 202,
            headers: { 'content-type': 'application/json' },
            data: jsonBlob
        })

        api.get.mockResolvedValue({
            data: {
                job_id: 'job-err',
                status: 'failed',
                error: 'Rendering engine failed'
            }
        })

        await expect(
            handleAsyncPdfDownload(fetchFn, 'failed.pdf', { pollIntervalMs: 10, maxPollAttempts: 5 })
        ).rejects.toThrow('Rendering engine failed')
    })
})
