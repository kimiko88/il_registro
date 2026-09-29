import { describe, it, expect, vi, beforeEach } from 'vitest'
import strikeService from '@/services/strikeService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn()
  }
}))

describe('strikeService & Strike Management Unit Suite', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('1. getStrikeNotices fetches strike notices from backend', async () => {
    const mockData = [
      { id: 'notice-1', title: 'Sciopero Generale', strike_date: '2026-10-15' }
    ]
    api.get.mockResolvedValue({ data: mockData })

    const res = await strikeService.getStrikeNotices()
    expect(api.get).toHaveBeenCalledWith('/strike-notices')
    expect(res).toEqual(mockData)
  })

  it('2. getStrikeNotice fetches single notice detail', async () => {
    const mockDetail = { id: 'notice-42', title: 'Sciopero Comparto Istruzione' }
    api.get.mockResolvedValue({ data: mockDetail })

    const res = await strikeService.getStrikeNotice('notice-42')
    expect(api.get).toHaveBeenCalledWith('/strike-notices/notice-42')
    expect(res).toEqual(mockDetail)
  })

  it('3. createStrikeNotice posts valid notice payload to backend', async () => {
    const payload = {
      title: 'Sciopero Nazionale Scuola',
      proclaimed_by: 'Sindacati Uniti',
      strike_date: '2026-11-04',
      declaration_deadline: '2026-11-02T12:00:00Z',
      content: 'Sciopero intera giornata',
      notes: 'Contingenti minimi garantiti',
      publish_to_bacheca: true
    }
    api.post.mockResolvedValue({ data: { id: 'created-id', ...payload } })

    const res = await strikeService.createStrikeNotice(payload)
    expect(api.post).toHaveBeenCalledWith('/strike-notices', payload)
    expect(res.id).toBe('created-id')
  })

  it('4. deleteStrikeNotice sends DELETE to /strike-notices/:id', async () => {
    api.delete.mockResolvedValue({ data: { message: 'avviso eliminato' } })

    const res = await strikeService.deleteStrikeNotice('notice-del-1')
    expect(api.delete).toHaveBeenCalledWith('/strike-notices/notice-del-1')
    expect(res.message).toBe('avviso eliminato')
  })

  it('5. submitDeclaration submits employee voluntary declaration', async () => {
    api.post.mockResolvedValue({
      data: { id: 'decl-10', intention: 'participates' }
    })

    const res = await strikeService.submitDeclaration('notice-99', 'participates')
    expect(api.post).toHaveBeenCalledWith('/strike-notices/notice-99/declare', { intention: 'participates' })
    expect(res.intention).toBe('participates')
  })

  it('6. getNoticeSummary fetches summary calculations and participation metrics', async () => {
    const mockSummary = {
      total_staff: 50,
      participates_count: 15,
      not_participates_count: 25,
      undecided_count: 5,
      unanswered_count: 5,
      participates_percent: 30.0
    }
    api.get.mockResolvedValue({ data: mockSummary })

    const res = await strikeService.getNoticeSummary('notice-summ-1')
    expect(api.get).toHaveBeenCalledWith('/strike-notices/notice-summ-1/summary')
    expect(res.participates_percent).toBe(30.0)
    expect(res.total_staff).toBe(50)
  })

  it('7. calculates participation percentage accurately', () => {
    const calculatePercent = (participates, total) => {
      if (!total || total <= 0) return 0
      return Math.round((participates / total) * 1000) / 10
    }

    expect(calculatePercent(12, 48)).toBe(25.0)
    expect(calculatePercent(0, 50)).toBe(0)
    expect(calculatePercent(50, 50)).toBe(100.0)
    expect(calculatePercent(10, 0)).toBe(0)
  })

  it('8. categorizes notices into active vs past based on date threshold', () => {
    const today = new Date('2026-09-24T00:00:00Z')
    const yesterday = new Date(today)
    yesterday.setDate(yesterday.getDate() - 1)

    const notices = [
      { id: 'past', strike_date: '2026-09-01' },
      { id: 'today', strike_date: '2026-09-24' },
      { id: 'future', strike_date: '2026-10-15' }
    ]

    const active = notices.filter(n => new Date(n.strike_date) >= yesterday)
    const past = notices.filter(n => new Date(n.strike_date) < yesterday)

    expect(active.map(n => n.id)).toEqual(['today', 'future'])
    expect(past.map(n => n.id)).toEqual(['past'])
  })
})
