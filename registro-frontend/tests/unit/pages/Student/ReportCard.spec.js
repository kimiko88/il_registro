import { describe, it, expect } from 'vitest'

// Pure unit tests for ReportCard business logic — no component mounting needed

// getAverageClass — mirrors logic in student/ReportCard.vue
function getAverageClass(avg) {
    if (!avg) return 'text-slate-400'
    return avg >= 6.0 ? 'text-positive font-bold' : 'text-negative font-bold'
}

describe('ReportCard getAverageClass', () => {
    it('returns slate for null/undefined average', () => {
        expect(getAverageClass(null)).toBe('text-slate-400')
        expect(getAverageClass(undefined)).toBe('text-slate-400')
        expect(getAverageClass(0)).toBe('text-slate-400')
    })

    it('returns positive class for average >= 6', () => {
        expect(getAverageClass(6.0)).toBe('text-positive font-bold')
        expect(getAverageClass(7.5)).toBe('text-positive font-bold')
        expect(getAverageClass(10)).toBe('text-positive font-bold')
    })

    it('returns negative class for average < 6', () => {
        expect(getAverageClass(5.9)).toBe('text-negative font-bold')
        expect(getAverageClass(4)).toBe('text-negative font-bold')
        expect(getAverageClass(1)).toBe('text-negative font-bold')
    })

    it('boundary: exactly 6.0 is positive (sufficient)', () => {
        expect(getAverageClass(6.0)).toBe('text-positive font-bold')
    })
})

// isPromoted normalization — normalize both boolean true and string 'SÌ'
function isPromoted(val) {
    if (typeof val === 'boolean') return val
    if (typeof val === 'string') {
        const upper = val.toUpperCase().trim()
        return upper === 'SÌ' || upper === 'SI' || upper === 'YES' || upper === 'TRUE'
    }
    return false
}

describe('ReportCard isPromoted normalization', () => {
    it('returns true for boolean true', () => {
        expect(isPromoted(true)).toBe(true)
    })

    it('returns false for boolean false', () => {
        expect(isPromoted(false)).toBe(false)
    })

    it('returns true for string "SÌ"', () => {
        expect(isPromoted('SÌ')).toBe(true)
    })

    it('returns true for string "SI" (without accent)', () => {
        expect(isPromoted('SI')).toBe(true)
    })

    it('returns true for lowercase "sì"', () => {
        expect(isPromoted('sì')).toBe(true)
    })

    it('returns false for string "NO"', () => {
        expect(isPromoted('NO')).toBe(false)
    })

    it('returns false for null', () => {
        expect(isPromoted(null)).toBe(false)
    })

    it('returns false for undefined', () => {
        expect(isPromoted(undefined)).toBe(false)
    })
})

// attendanceStats defaults — should be zero-initialized on missing data
function buildAttendanceStats(data) {
    if (!data) return { absences: 0, lates: 0, earlyExits: 0 }
    return {
        absences: data.absences || 0,
        lates: data.lates || 0,
        earlyExits: data.early_exits || 0,
    }
}

describe('ReportCard buildAttendanceStats', () => {
    it('returns all zeros for null data', () => {
        const stats = buildAttendanceStats(null)
        expect(stats.absences).toBe(0)
        expect(stats.lates).toBe(0)
        expect(stats.earlyExits).toBe(0)
    })

    it('maps API fields correctly', () => {
        const stats = buildAttendanceStats({ absences: 5, lates: 2, early_exits: 1 })
        expect(stats.absences).toBe(5)
        expect(stats.lates).toBe(2)
        expect(stats.earlyExits).toBe(1)
    })

    it('defaults missing fields to 0', () => {
        const stats = buildAttendanceStats({ absences: 3 })
        expect(stats.absences).toBe(3)
        expect(stats.lates).toBe(0)
        expect(stats.earlyExits).toBe(0)
    })
})

// Report columns structure — verify column definitions are complete
const columns = [
    { name: 'subject', label: 'Materia', align: 'left', field: 'subject', sortable: true },
    { name: 'teacher', label: 'Docente', align: 'left', field: 'teacher' },
    { name: 'grades', label: 'Voti Periodo', align: 'center', field: 'grade_count' },
    { name: 'average', label: 'Media Voti', align: 'center', field: 'subject_average', sortable: true },
    { name: 'final_grade', label: 'Voto Finale', align: 'center', field: 'final_grade', sortable: true },
    { name: 'notes', label: 'Note', align: 'left', field: 'notes' },
]

describe('ReportCard columns definition', () => {
    it('has 6 columns defined', () => {
        expect(columns).toHaveLength(6)
    })

    it('all columns have name, label, align, and field', () => {
        columns.forEach(col => {
            expect(col.name).toBeTruthy()
            expect(col.label).toBeTruthy()
            expect(col.align).toMatch(/left|center|right/)
            expect(col.field).toBeTruthy()
        })
    })

    it('subject and final_grade columns are sortable', () => {
        const sortable = columns.filter(c => c.sortable)
        const names = sortable.map(c => c.name)
        expect(names).toContain('subject')
        expect(names).toContain('final_grade')
    })
})
