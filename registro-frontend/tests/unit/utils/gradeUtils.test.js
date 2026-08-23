import { describe, it, expect } from 'vitest'
import {
  ITALIAN_GRADE_OPTIONS,
  gradeToNumeric,
  formatGrade,
  getGradeColor
} from '@/utils/gradeUtils'

describe('gradeUtils — Italian Grade Conversion & Formatting', () => {
  describe('ITALIAN_GRADE_OPTIONS', () => {
    it('contains standard Italian grade scale including halves and minuses', () => {
      expect(ITALIAN_GRADE_OPTIONS).toContain('10')
      expect(ITALIAN_GRADE_OPTIONS).toContain('10-')
      expect(ITALIAN_GRADE_OPTIONS).toContain('9½')
      expect(ITALIAN_GRADE_OPTIONS).toContain('6+')
      expect(ITALIAN_GRADE_OPTIONS).toContain('6-')
      expect(ITALIAN_GRADE_OPTIONS).toContain('A')
    })
  })

  describe('gradeToNumeric', () => {
    it('returns null for empty or undefined values', () => {
      expect(gradeToNumeric(null)).toBeNull()
      expect(gradeToNumeric(undefined)).toBeNull()
      expect(gradeToNumeric('')).toBeNull()
      expect(gradeToNumeric(NaN)).toBeNull()
    })

    it('returns numbers unchanged', () => {
      expect(gradeToNumeric(8)).toBe(8)
      expect(gradeToNumeric(6.5)).toBe(6.5)
      expect(gradeToNumeric(-1)).toBe(-1)
    })

    it('converts integer and decimal string grades', () => {
      expect(gradeToNumeric('10')).toBe(10)
      expect(gradeToNumeric('7.5')).toBe(7.5)
      expect(gradeToNumeric('7,5')).toBe(7.5)
    })

    it('converts plus and minus grades correctly', () => {
      expect(gradeToNumeric('10-')).toBe(9.75)
      expect(gradeToNumeric('6+')).toBe(6.25)
      expect(gradeToNumeric('+6')).toBe(6.25)
      expect(gradeToNumeric('7-')).toBe(6.75)
      expect(gradeToNumeric('5+')).toBe(5.25)
    })

    it('converts half symbols and fractions', () => {
      expect(gradeToNumeric('9½')).toBe(9.5)
      expect(gradeToNumeric('8 1/2')).toBe(8.5)
      expect(gradeToNumeric('6½')).toBe(6.5)
    })

    it('converts grade ranges and test ratios', () => {
      // Grade range e.g. 7/8 or 6-7 -> average
      expect(gradeToNumeric('7/8')).toBe(7.5)
      expect(gradeToNumeric('6-7')).toBe(6.5)
      // Test score ratio e.g. 18/20 -> 9
      expect(gradeToNumeric('18/20')).toBe(9)
      expect(gradeToNumeric('15/30')).toBe(5)
    })

    it('handles special grade labels (A, NC, E, S, INS, O, B, D, Eccellente, Avanzato, Mediocre)', () => {
      expect(gradeToNumeric('A')).toBe(-1)
      expect(gradeToNumeric('ASSENTE')).toBe(-1)
      expect(gradeToNumeric('ABSENT')).toBe(-1)
      expect(gradeToNumeric('NC')).toBeNull()
      expect(gradeToNumeric('E')).toBeNull()
      expect(gradeToNumeric('S')).toBe(6)
      expect(gradeToNumeric('INS')).toBe(5)
      expect(gradeToNumeric('O')).toBe(10)
      expect(gradeToNumeric('D')).toBe(8)
      expect(gradeToNumeric('B')).toBe(7)
      expect(gradeToNumeric('ECCELLENTE')).toBe(10)
      expect(gradeToNumeric('AVANZATO')).toBe(10)
      expect(gradeToNumeric('MEDIOCRE')).toBe(5)
      expect(gradeToNumeric('GRAVEMENTE INSUFFICIENTE')).toBe(3)
    })

    it('rejects out of bounds grades outside [1.0, 10.0] scale', () => {
      expect(gradeToNumeric(0)).toBeNull()
      expect(gradeToNumeric(11)).toBeNull()
      expect(gradeToNumeric('0')).toBeNull()
      expect(gradeToNumeric('11')).toBeNull()
      expect(gradeToNumeric('-5')).toBeNull()
      expect(gradeToNumeric('0/5')).toBeNull()
      expect(gradeToNumeric(15)).toBeNull()
    })
  })

  describe('formatGrade', () => {
    it('formats null, undefined, empty, or hyphen as -', () => {
      expect(formatGrade(null)).toBe('-')
      expect(formatGrade(undefined)).toBe('-')
      expect(formatGrade('')).toBe('-')
      expect(formatGrade('-')).toBe('-')
    })

    it('formats -1 as absent A', () => {
      expect(formatGrade(-1)).toBe('A')
      expect(formatGrade('-1')).toBe('A')
    })

    it('formats whole integers cleanly', () => {
      expect(formatGrade(10)).toBe('10')
      expect(formatGrade(6)).toBe('6')
      expect(formatGrade('8')).toBe('8')
    })

    it('clamps values beyond 10 to 10 and does not format as 11', () => {
      expect(formatGrade(10.02)).toBe('10')
      expect(formatGrade(10.5)).toBe('10')
      expect(formatGrade(11)).toBe('10')
      expect(formatGrade(0.5)).toBe('1')
    })

    it('formats half and quarter increments correctly', () => {
      expect(formatGrade(5.5)).toBe('5½')
      expect(formatGrade(7.25)).toBe('7+')
      expect(formatGrade(7.75)).toBe('8-')
    })

    it('corrects decimal rounding without producing anomalous 5,10 string', () => {
      expect(formatGrade(5.95)).toBe('6')
      expect(formatGrade(5.96)).toBe('6')
      expect(formatGrade(5.99)).toBe('6')
      expect(formatGrade(5.3, ',')).toBe('5,3')
      expect(formatGrade(5.3, '.')).toBe('5.3')
    })
  })

  describe('getGradeColor', () => {
    it('returns grey-3 for absent (-1 or A)', () => {
      expect(getGradeColor(-1)).toBe('grey-3')
      expect(getGradeColor('A')).toBe('grey-3')
    })

    it('returns red-2 for grades < 5', () => {
      expect(getGradeColor(4)).toBe('red-2')
      expect(getGradeColor('4½')).toBe('red-2')
      expect(getGradeColor('4+')).toBe('red-2')
    })

    it('returns amber-2 for grades between 5 and 5.99', () => {
      expect(getGradeColor(5)).toBe('amber-2')
      expect(getGradeColor('5½')).toBe('amber-2')
      expect(getGradeColor('5+')).toBe('amber-2')
      expect(getGradeColor('6-')).toBe('amber-2')
    })

    it('returns green-2 for sufficient grades >= 6', () => {
      expect(getGradeColor(6)).toBe('green-2')
      expect(getGradeColor('7½')).toBe('green-2')
      expect(getGradeColor('10')).toBe('green-2')
      expect(getGradeColor('10-')).toBe('green-2')
    })

    it('returns white for empty or invalid input', () => {
      expect(getGradeColor(null)).toBe('white')
      expect(getGradeColor('')).toBe('white')
      expect(getGradeColor('-')).toBe('white')
      expect(getGradeColor(11)).toBe('white')
      expect(getGradeColor('11')).toBe('white')
    })
  })
})
