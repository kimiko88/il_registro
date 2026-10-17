import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  RELIGION_JUDGMENT_OPTIONS,
  RELIGION_JUDGMENT_MAP,
  isReligionJudgment,
  religionJudgmentToNumeric,
  getGradeColor,
  formatGrade
} from '@/utils/gradeUtils'
import religionService from '@/services/religionService'
import { api } from '@/boot/axios'

vi.mock('@/boot/axios', () => ({
  api: {
    get: vi.fn(),
    put: vi.fn(),
    post: vi.fn()
  }
}))

describe('Religion (IRC) Grading & Service Unit Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('RELIGION_JUDGMENT_OPTIONS & RELIGION_JUDGMENT_MAP', () => {
    it('contains the 6 official ministerial judgments', () => {
      expect(RELIGION_JUDGMENT_OPTIONS).toEqual([
        'Ottimo',
        'Distinto',
        'Buono',
        'Sufficiente',
        'Insufficiente',
        'Non classificabile'
      ])
    })

    it('maps all 6 judgments to their expected numeric equivalence', () => {
      expect(RELIGION_JUDGMENT_MAP['Ottimo']).toBe(10)
      expect(RELIGION_JUDGMENT_MAP['Distinto']).toBe(8)
      expect(RELIGION_JUDGMENT_MAP['Buono']).toBe(7)
      expect(RELIGION_JUDGMENT_MAP['Sufficiente']).toBe(6)
      expect(RELIGION_JUDGMENT_MAP['Insufficiente']).toBe(4)
      expect(RELIGION_JUDGMENT_MAP['Non classificabile']).toBe(0)
    })
  })

  describe('isReligionJudgment', () => {
    it('recognizes valid religion judgments regardless of case and whitespace', () => {
      expect(isReligionJudgment('ottimo')).toBe(true)
      expect(isReligionJudgment('  DISTINTO  ')).toBe(true)
      expect(isReligionJudgment('buono')).toBe(true)
      expect(isReligionJudgment('Sufficiente')).toBe(true)
      expect(isReligionJudgment('insufficiente')).toBe(true)
      expect(isReligionJudgment('Non classificabile')).toBe(true)
    })

    it('rejects numerical grades and non-religion judgments', () => {
      expect(isReligionJudgment('8')).toBe(false)
      expect(isReligionJudgment('10')).toBe(false)
      expect(isReligionJudgment('Discreto')).toBe(false)
      expect(isReligionJudgment('Eccellente')).toBe(false)
      expect(isReligionJudgment(null)).toBe(false)
      expect(isReligionJudgment('')).toBe(false)
    })
  })

  describe('religionJudgmentToNumeric', () => {
    it('converts judgment strings to numeric values correctly', () => {
      expect(religionJudgmentToNumeric('Ottimo')).toBe(10)
      expect(religionJudgmentToNumeric('distinto')).toBe(8)
      expect(religionJudgmentToNumeric('  Buono ')).toBe(7)
      expect(religionJudgmentToNumeric('SUFFICIENTE')).toBe(6)
      expect(religionJudgmentToNumeric('insufficiente')).toBe(4)
      expect(religionJudgmentToNumeric('Non classificabile')).toBe(0)
    })

    it('returns null for unknown judgments', () => {
      expect(religionJudgmentToNumeric('Random')).toBeNull()
      expect(religionJudgmentToNumeric(null)).toBeNull()
      expect(religionJudgmentToNumeric('')).toBeNull()
    })
  })

  describe('getGradeColor & formatGrade for IRC', () => {
    it('assigns appropriate palette colors to IRC judgments', () => {
      expect(getGradeColor('Ottimo')).toBe('green-2')
      expect(getGradeColor('Distinto')).toBe('teal-2')
      expect(getGradeColor('Buono')).toBe('light-green-2')
      expect(getGradeColor('Sufficiente')).toBe('amber-2')
      expect(getGradeColor('Insufficiente')).toBe('red-2')
      expect(getGradeColor('Non classificabile')).toBe('grey-4')
      expect(getGradeColor(0)).toBe('grey-4')
    })

    it('formats 0 as Non class. and preserves text judgments', () => {
      expect(formatGrade(0)).toBe('Non class.')
      expect(formatGrade('Ottimo')).toBe('Ottimo')
      expect(formatGrade('Buono')).toBe('Buono')
    })
  })

  describe('religionService API Calls', () => {
    it('getStudentChoice calls GET /students/:id/religion-choice', async () => {
      api.get.mockResolvedValueOnce({ data: { choice: 'avvalente' } })
      const res = await religionService.getStudentChoice('stud-123')
      expect(api.get).toHaveBeenCalledWith('/students/stud-123/religion-choice')
      expect(res.data.choice).toBe('avvalente')
    })

    it('setStudentChoice calls PUT /students/:id/religion-choice', async () => {
      api.put.mockResolvedValueOnce({ data: { success: true } })
      const res = await religionService.setStudentChoice('stud-123', 'non_avvalente')
      expect(api.put).toHaveBeenCalledWith('/students/stud-123/religion-choice', {
        choice: 'non_avvalente'
      })
      expect(res.data.success).toBe(true)
    })

    it('listChoices calls GET /secretary/religion-choices with query params', async () => {
      api.get.mockResolvedValueOnce({ data: [{ student_id: 's1', choice: 'avvalente' }] })
      await religionService.listChoices({ class_id: 'c1' })
      expect(api.get).toHaveBeenCalledWith('/secretary/religion-choices', { params: { class_id: 'c1' } })
    })

    it('batchSetChoices calls POST /secretary/religion-choices/batch', async () => {
      api.post.mockResolvedValueOnce({ data: { updated: 2 } })
      const res = await religionService.batchSetChoices(['s1', 's2'], 'attivita_alternativa')
      expect(api.post).toHaveBeenCalledWith('/secretary/religion-choices/batch', {
        student_ids: ['s1', 's2'],
        choice: 'attivita_alternativa'
      })
      expect(res.data.updated).toBe(2)
    })

    it('getClassChoices calls GET /classes/:id/religion-choices', async () => {
      api.get.mockResolvedValueOnce({ data: [] })
      await religionService.getClassChoices('class-99')
      expect(api.get).toHaveBeenCalledWith('/classes/class-99/religion-choices')
    })
  })
})
