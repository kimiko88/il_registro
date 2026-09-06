import { describe, it, expect, vi, beforeEach } from 'vitest'
import { sanitizeCsvFormula, useTableExport } from '@/composables/useTableExport'
import * as QuasarModule from 'quasar'

vi.mock('quasar', () => ({
  exportFile: vi.fn(),
  useQuasar: vi.fn()
}))

describe('useTableExport & sanitizeCsvFormula', () => {
  let notifyMock

  beforeEach(() => {
    vi.clearAllMocks()
    notifyMock = vi.fn()
    QuasarModule.useQuasar.mockReturnValue({
      notify: notifyMock
    })
  })

  describe('sanitizeCsvFormula', () => {
    it('returns empty string for null and undefined', () => {
      expect(sanitizeCsvFormula(null)).toBe('')
      expect(sanitizeCsvFormula(undefined)).toBe('')
    })

    it('sanitizes formula injection characters (=, +, -, @, \\t, \\r)', () => {
      expect(sanitizeCsvFormula('=SUM(A1:A10)')).toBe("'=SUM(A1:A10)")
      expect(sanitizeCsvFormula('+12345')).toBe("'+12345")
      expect(sanitizeCsvFormula('-100')).toBe("'-100")
      expect(sanitizeCsvFormula('@attack')).toBe("'@attack")
      expect(sanitizeCsvFormula('\tcommand')).toBe("'\tcommand")
      expect(sanitizeCsvFormula('\rcommand')).toBe("'\rcommand")
    })

    it('escapes double quotes properly', () => {
      expect(sanitizeCsvFormula('Text "with" quotes')).toBe('Text ""with"" quotes')
      expect(sanitizeCsvFormula('=Formula "quotes"')).toBe('\'=Formula ""quotes""')
    })

    it('preserves normal text and numbers', () => {
      expect(sanitizeCsvFormula('Mario Rossi')).toBe('Mario Rossi')
      expect(sanitizeCsvFormula(12345)).toBe('12345')
      expect(sanitizeCsvFormula(0)).toBe('0')
    })
  })

  describe('exportTableCsv', () => {
    it('notifies warning and returns false if rows are empty', () => {
      const { exportTableCsv } = useTableExport()
      const result = exportTableCsv({ rows: [] })

      expect(result).toBe(false)
      expect(notifyMock).toHaveBeenCalledWith(expect.objectContaining({ type: 'warning' }))
      expect(QuasarModule.exportFile).not.toHaveBeenCalled()
    })

    it('exports rows with field, formatter and formula sanitization', () => {
      QuasarModule.exportFile.mockReturnValue(true)

      const { exportTableCsv } = useTableExport()
      const columns = [
        { name: 'name', label: 'Nome', field: 'name' },
        { name: 'score', label: 'Punteggio', field: row => row.score, format: val => `${val} pt` },
        { name: 'formula', label: 'Formula', field: 'formula' }
      ]
      const rows = [
        { name: 'Mario Rossi', score: 95, formula: '=CMD|' },
        { name: 'Luigi "Verdi"', score: 80, formula: 'Normale' }
      ]

      const success = exportTableCsv({
        filename: 'studenti.csv',
        columns,
        rows
      })

      expect(success).toBe(true)
      expect(QuasarModule.exportFile).toHaveBeenCalledTimes(1)
      const [filename, content, mimeType] = QuasarModule.exportFile.mock.calls[0]
      expect(filename).toBe('studenti.csv')
      expect(mimeType).toBe('text/csv;charset=utf-8;')
      expect(content).toContain('"Nome","Punteggio","Formula"')
      expect(content).toContain('"Mario Rossi","95 pt","\'=CMD|"')
      expect(content).toContain('"Luigi ""Verdi""","80 pt","Normale"')
      expect(notifyMock).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }))
    })

    it('notifies error when browser denies file download', () => {
      QuasarModule.exportFile.mockReturnValue(false)

      const { exportTableCsv } = useTableExport()
      const success = exportTableCsv({
        columns: [{ name: 'id', label: 'ID', field: 'id' }],
        rows: [{ id: 1 }]
      })

      expect(success).toBe(false)
      expect(notifyMock).toHaveBeenCalledWith(expect.objectContaining({ type: 'negative' }))
    })
  })
})
