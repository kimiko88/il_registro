import { describe, it, expect } from 'vitest'
import itIT from '@/i18n/it-IT/index.js'
import enUS from '@/i18n/en-US/index.js'

// Formatters replicating Timecard.vue logic for unit verification
function formatHoursAndMinutes(hoursFloat) {
  if (hoursFloat === null || hoursFloat === undefined || isNaN(Number(hoursFloat))) {
    return '0h 00m'
  }
  const totalMinutes = Math.round(Number(hoursFloat) * 60)
  const isNegative = totalMinutes < 0
  const absMinutes = Math.abs(totalMinutes)
  const h = Math.floor(absMinutes / 60)
  const m = absMinutes % 60
  const formattedMinutes = String(m).padStart(2, '0')
  return `${isNegative ? '-' : ''}${h}h ${formattedMinutes}m`
}

function formatBalanceHours(hoursFloat) {
  if (hoursFloat === null || hoursFloat === undefined || isNaN(Number(hoursFloat))) {
    return '0h 00m'
  }
  const totalMinutes = Math.round(Number(hoursFloat) * 60)
  const isNegative = totalMinutes < 0
  const absMinutes = Math.abs(totalMinutes)
  const h = Math.floor(absMinutes / 60)
  const m = absMinutes % 60
  const formattedMinutes = String(m).padStart(2, '0')
  const sign = isNegative ? '-' : '+'
  return `${sign}${h}h ${formattedMinutes}m`
}

function formatDailyWorked(row) {
  if (!row) return '—'
  let minutes = null
  if (row.worked_minutes !== undefined && row.worked_minutes !== null) {
    minutes = Number(row.worked_minutes)
  } else if (row.hours !== undefined && row.hours !== null) {
    minutes = Math.round(Number(row.hours) * 60)
  }
  if (minutes === null || (minutes === 0 && !row.entry_time)) {
    return '—'
  }
  const isNegative = minutes < 0
  const absMinutes = Math.abs(minutes)
  const h = Math.floor(absMinutes / 60)
  const m = absMinutes % 60
  const formattedMinutes = String(m).padStart(2, '0')
  return `${isNegative ? '-' : ''}${h}h ${formattedMinutes}m`
}

function formatTimeHHMMSS(val) {
  if (!val) return '—'
  if (typeof val === 'string') {
    const trimmed = val.trim()
    if (/^\d{2}:\d{2}:\d{2}$/.test(trimmed)) return trimmed
    if (/^\d{2}:\d{2}$/.test(trimmed)) return `${trimmed}:00`
  }
  try {
    const d = new Date(val)
    if (isNaN(d.getTime())) return String(val)
    const hh = String(d.getHours()).padStart(2, '0')
    const mm = String(d.getMinutes()).padStart(2, '0')
    const ss = String(d.getSeconds()).padStart(2, '0')
    return `${hh}:${mm}:${ss}`
  } catch {
    return String(val)
  }
}

function getEntryStatusLabel(st, dict) {
  const norm = String(st || '').toLowerCase().trim()
  const t = (k) => {
    const parts = k.split('.')
    let curr = dict
    for (const p of parts) {
      if (curr && curr[p]) curr = curr[p]
      else return null
    }
    return curr
  }

  switch (norm) {
    case 'present':
    case 'presente':
      return t('timecard.entryStatusPresent')
    case 'absent':
    case 'assente':
      return t('timecard.entryStatusAbsent')
    case 'late':
    case 'in ritardo':
      return t('staffAttendance.statusLate') || 'In Ritardo'
    case 'on_strike':
    case 'strike':
    case 'in sciopero':
      return t('staffAttendance.statusOnStrike') || 'In Sciopero'
    case 'sick':
    case 'sick_leave':
    case 'malattia':
      return t('timecard.entryStatusSick')
    case 'leave':
    case 'ferie':
      return t('timecard.entryStatusLeave')
    case 'permit':
    case 'permesso':
      return t('staffAttendance.statusPermit') || t('timecard.entryStatusLeave')
    case 'mission':
    case 'missione':
      return t('staffAttendance.statusMission') || 'Missione'
    case 'holiday':
    case 'festivo':
      return t('timecard.entryStatusHoliday')
    default:
      return st || t('timecard.entryStatusPresent')
  }
}

describe('Timecard formatting and localization', () => {
  describe('Worked Hours formatting without decimals', () => {
    it('formats raw float 0.9166666666666666 h to 0h 55m', () => {
      expect(formatHoursAndMinutes(0.9166666666666666)).toBe('0h 55m')
    })

    it('formats integer and fractional hours correctly', () => {
      expect(formatHoursAndMinutes(7.5)).toBe('7h 30m')
      expect(formatHoursAndMinutes(156)).toBe('156h 00m')
      expect(formatHoursAndMinutes(0)).toBe('0h 00m')
      expect(formatHoursAndMinutes(null)).toBe('0h 00m')
    })

    it('formats overtime balance with sign', () => {
      expect(formatBalanceHours(-155.08333333333334)).toBe('-155h 05m')
      expect(formatBalanceHours(2.5)).toBe('+2h 30m')
      expect(formatBalanceHours(0)).toBe('+0h 00m')
    })

    it('formats daily row worked time from worked_minutes or hours', () => {
      expect(formatDailyWorked({ worked_minutes: 55 })).toBe('0h 55m')
      expect(formatDailyWorked({ hours: 1.25 })).toBe('1h 15m')
      expect(formatDailyWorked({ worked_minutes: 0, entry_time: '14:00:00' })).toBe('0h 00m')
      expect(formatDailyWorked({ worked_minutes: 0 })).toBe('—')
      expect(formatDailyWorked(null)).toBe('—')
    })
  })

  describe('Entry / Exit time formatting to HH:MM:SS', () => {
    it('formats ISO timestamps to HH:MM:SS', () => {
      const formatted = formatTimeHHMMSS('2026-09-14T14:09:21.349927Z')
      expect(formatted).toMatch(/^\d{2}:\d{2}:\d{2}$/)
    })

    it('formats string times directly without modifying valid HH:MM:SS', () => {
      expect(formatTimeHHMMSS('14:09:21')).toBe('14:09:21')
      expect(formatTimeHHMMSS('08:30')).toBe('08:30:00')
    })

    it('handles empty / null values', () => {
      expect(formatTimeHHMMSS(null)).toBe('—')
      expect(formatTimeHHMMSS('')).toBe('—')
    })
  })

  describe('Attendance Status Localization', () => {
    it('localizes lowercase DB status "present" in Italian', () => {
      expect(getEntryStatusLabel('present', itIT)).toBe('Presente')
    })

    it('localizes lowercase DB status "present" in English', () => {
      expect(getEntryStatusLabel('present', enUS)).toBe('Present')
    })

    it('localizes absent, late, and on_strike', () => {
      expect(getEntryStatusLabel('absent', itIT)).toBe('Assente')
      expect(getEntryStatusLabel('absent', enUS)).toBe('Absent')
      expect(getEntryStatusLabel('late', itIT)).toBe('In Ritardo')
      expect(getEntryStatusLabel('on_strike', itIT)).toBe('In Sciopero')
    })

    it('localizes sick leave, leave and holiday', () => {
      expect(getEntryStatusLabel('sick', itIT)).toBe('Malattia')
      expect(getEntryStatusLabel('leave', itIT)).toBe('Ferie / Permesso')
      expect(getEntryStatusLabel('holiday', itIT)).toBe('Festivo / Chiusura')
    })
  })
})
