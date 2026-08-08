import { describe, it, expect, beforeEach } from 'vitest'
import { useGradeFormatter } from '@/composables/useGradeFormatter'

describe('useGradeFormatter composable', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('formats decimals with comma by default in Italian', () => {
    const { formatGrade, formatDecimal } = useGradeFormatter()
    expect(formatGrade(7.5)).toBe('7,5')
    expect(formatDecimal(8.4166, 2)).toBe('8,42')
  })

  it('formats decimals with dot when user_decimal_separator is set to dot', () => {
    localStorage.setItem('user_decimal_separator', '.')
    const { formatGrade, formatDecimal } = useGradeFormatter()
    expect(formatGrade(7.5)).toBe('7.5')
    expect(formatDecimal(8.4166, 2)).toBe('8.42')
  })

  it('handles fractional format setting correctly', () => {
    localStorage.setItem('teacher_register_settings', JSON.stringify({ gradeFormat: 'fractional' }))
    const { formatGrade } = useGradeFormatter()
    expect(formatGrade(7.5)).toBe('7½')
    expect(formatGrade(6.25)).toBe('6+')
    expect(formatGrade(7.75)).toBe('8-')
  })

  it('handles centesimal format setting correctly', () => {
    localStorage.setItem('teacher_register_settings', JSON.stringify({ gradeFormat: 'centesimal' }))
    const { formatGrade } = useGradeFormatter()
    expect(formatGrade(7.5)).toBe('75/100')
  })
})
