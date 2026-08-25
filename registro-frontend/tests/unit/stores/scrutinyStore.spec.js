import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useScrutinyStore } from 'src/stores/scrutiny'
import api from 'src/services/api'

// Mock API
vi.mock('src/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn()
    }
}))

describe('Scrutiny Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useScrutinyStore()
        vi.clearAllMocks()
    })

    describe('getters', () => {
        it('calculates classesWithRisk correctly', () => {
            store.overview = [
                { id: '1A', name: '1A', risk_count: 2, status: 'pending' },
                { id: '2B', name: '2B', risk_count: 0, status: 'finalized' },
                { id: '3C', name: '3C', has_risk: true, status: 'pending' }
            ]

            expect(store.classesWithRisk).toHaveLength(2)
            expect(store.classesWithRisk.map(c => c.id)).toEqual(['1A', '3C'])
        })

        it('calculates promotedCount correctly', () => {
            store.currentReport = {
                students: [
                    { id: 's1', name: 'Mario Rossi', promoted: 'SÌ' },
                    { id: 's2', name: 'Luca Bianchi', promoted: true },
                    { id: 's3', name: 'Giuseppe Verdi', promoted: 'NO' }
                ]
            }

            expect(store.promotedCount).toBe(2)
        })

        it('calculates pendingCount correctly', () => {
            store.overview = [
                { id: '1A', status: 'pending' },
                { id: '2B', status: 'finalized' },
                { id: '3C', status: 'pending' }
            ]

            expect(store.pendingCount).toBe(2)
        })
    })

    describe('actions', () => {
        it('fetches scrutiny overview successfully', async () => {
            const mockOverview = [
                { id: '1A', name: '1A', status: 'pending' }
            ]
            api.get.mockResolvedValue({ data: mockOverview })

            const res = await store.fetchOverview('school-1')

            expect(api.get).toHaveBeenCalledWith('/scrutiny/overview', { params: { school_id: 'school-1' } })
            expect(res).toEqual(mockOverview)
            expect(store.overview).toEqual(mockOverview)
            expect(store.loading).toBe(false)
            expect(store.error).toBeNull()
        })

        it('fetches class report successfully', async () => {
            const mockReport = {
                class_id: 'class-1A',
                students: [{ id: 's1', promoted: 'SÌ' }]
            }
            api.get.mockResolvedValue({ data: mockReport })

            const res = await store.fetchClassReport('class-1A')

            expect(api.get).toHaveBeenCalledWith('/scrutiny/class/class-1A/report')
            expect(res).toEqual(mockReport)
            expect(store.currentReport).toEqual(mockReport)
            expect(store.loading).toBe(false)
        })

        it('finalizes scrutiny for a class and re-fetches overview', async () => {
            api.post.mockResolvedValue({ data: { message: 'Scrutinio finalizzato' } })
            api.get.mockResolvedValue({ data: [{ id: '1A', status: 'finalized' }] })

            const res = await store.finalizeScrutiny('class-1A', 'school-1')

            expect(api.post).toHaveBeenCalledWith('/scrutiny/class/class-1A/finalize')
            expect(api.get).toHaveBeenCalledWith('/scrutiny/overview', { params: { school_id: 'school-1' } })
            expect(res).toEqual({ message: 'Scrutinio finalizzato' })
            expect(store.loading).toBe(false)
        })

        it('handles errors during overview fetch', async () => {
            api.get.mockRejectedValue({ response: { data: { error: 'Errore server' } } })

            await expect(store.fetchOverview()).rejects.toBeDefined()
            expect(store.error).toBe('Errore server')
            expect(store.loading).toBe(false)
        })
    })
})
