import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAgendaStore } from 'src/stores/agenda'
import api from 'src/services/api'

// Mock API
vi.mock('src/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn(),
        patch: vi.fn(),
        delete: vi.fn()
    }
}))

describe('Agenda Store', () => {
    let store

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useAgendaStore()
        vi.clearAllMocks()
    })

    it('fetches agenda events successfully', async () => {
        const mockEvents = [
            { id: 'ev-1', title: 'Verifica di Matematica', event_type: 'test', date: '2026-09-10' }
        ]
        api.get.mockResolvedValue({ data: mockEvents })

        const res = await store.fetchAgenda({ class_id: 'class-1' })

        expect(api.get).toHaveBeenCalledWith('/agenda', { params: { class_id: 'class-1' } })
        expect(res).toEqual(mockEvents)
        expect(store.events).toEqual(mockEvents)
        expect(store.loading).toBe(false)
        expect(store.error).toBeNull()
    })

    it('handles fetch agenda error', async () => {
        api.get.mockRejectedValue({ response: { data: { error: 'Impossibile caricare agenda' } } })

        await expect(store.fetchAgenda()).rejects.toBeDefined()
        expect(store.error).toBe('Impossibile caricare agenda')
        expect(store.loading).toBe(false)
    })

    it('creates agenda event successfully', async () => {
        const payload = { title: 'Compiti Fisica', event_type: 'homework', date: '2026-09-12' }
        const createdEvent = { id: 'ev-2', ...payload }
        api.post.mockResolvedValue({ data: createdEvent })

        const res = await store.createEvent(payload)

        expect(api.post).toHaveBeenCalledWith('/agenda', payload)
        expect(res).toEqual(createdEvent)
        expect(store.events).toContainEqual(createdEvent)
        expect(store.loading).toBe(false)
    })

    it('updates agenda event successfully', async () => {
        store.events = [
            { id: 'ev-1', title: 'Verifica Matematica', date: '2026-09-10' }
        ]
        const updatePayload = { title: 'Verifica Matematica Spostata' }
        api.patch.mockResolvedValue({ data: { id: 'ev-1', title: 'Verifica Matematica Spostata', date: '2026-09-10' } })

        const res = await store.updateEvent('ev-1', updatePayload)

        expect(api.patch).toHaveBeenCalledWith('/agenda/ev-1', updatePayload)
        expect(res.title).toBe('Verifica Matematica Spostata')
        expect(store.events[0].title).toBe('Verifica Matematica Spostata')
    })

    it('deletes agenda event successfully', async () => {
        store.events = [
            { id: 'ev-1', title: 'Compiti 1' },
            { id: 'ev-2', title: 'Compiti 2' }
        ]
        api.delete.mockResolvedValue({})

        await store.deleteEvent('ev-1')

        expect(api.delete).toHaveBeenCalledWith('/agenda/ev-1')
        expect(store.events).toHaveLength(1)
        expect(store.events[0].id).toBe('ev-2')
    })
})
