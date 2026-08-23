// Shared Italian grade options and conversion utilities

export const ITALIAN_GRADE_OPTIONS = [
  '10', '10-',
  '9½', '9+', '9', '9-',
  '8½', '8+', '8', '8-',
  '7½', '7+', '7', '7-',
  '6½', '6+', '6', '6-',
  '5½', '5+', '5', '5-',
  '4½', '4+', '4', '4-',
  '3½', '3+', '3', '3-',
  '2½', '2+', '2', '2-',
  '1½', '1+', '1',
  'A'
]

export function gradeToNumeric(val) {
  if (val === undefined || val === null || val === '') return null
  if (typeof val === 'number') {
    if (isNaN(val)) return null
    if (val === -1) return -1
    if (val < 1.0 || val > 10.0) return null
    return val
  }
  const clean = String(val).trim().toUpperCase()
  if (clean === 'A' || clean === 'ASSENTE' || clean === 'ABSENT') return -1
  if (clean === 'NC' || clean === 'NON CLASSIFICATO' || clean === 'E' || clean === 'ESENTE' || clean === 'NV' || clean === 'NON VALUTATO') return null

  // Judgments to standard numerical scale
  if (clean === 'ECCELLENTE' || clean === 'OTTIMO' || clean === 'O' || clean === 'AVANZATO') return 10
  if (clean === 'DISTINTO' || clean === 'D') return 8
  if (clean === 'BUONO' || clean === 'B' || clean === 'INTERMEDIO') return 7
  if (clean === 'DISCRETO' || clean === 'BASE') return 6.5
  if (clean === 'S' || clean === 'SUFF' || clean === 'SUFFICIENTE') return 6
  if (clean === 'QUASI SUFFICIENTE' || clean === 'QUASI SUFF' || clean === 'MEDIOCRE' || clean === 'INIZIALE') return 5
  if (clean === 'INS' || clean === 'INSUFFICIENTE' || clean === 'NON RAGGIUNTO') return 5
  if (clean === 'GRAVEMENTE INSUFFICIENTE' || clean === 'GRAVE INSUFFICIENTE' || clean === 'GRAVEMENTE INS') return 3

  let candidate = null

  // Handle half notation first (e.g. 8 1/2 or 9½) before general / fraction split
  if (clean.endsWith('1/2') || clean.endsWith('½')) {
    const base = parseFloat(clean.replace('1/2', '').replace('½', '').trim().replace(',', '.'))
    if (!isNaN(base)) candidate = base + 0.5
  } else if (clean.includes('/')) {
    // Handle ratio scores like 18/20 or grade ranges like 7/8
    const parts = clean.split('/')
    if (parts.length === 2) {
      const n1 = parseFloat(parts[0].replace(',', '.'))
      const n2 = parseFloat(parts[1].replace(',', '.'))
      if (!isNaN(n1) && !isNaN(n2) && n2 > 0) {
        if (n1 >= 1 && n1 <= 10 && n2 >= 1 && n2 <= 10 && Math.abs(n2 - n1) <= 1.5) {
          // Grade range e.g. 7/8 -> 7.5, 6/7 -> 6.5
          candidate = (n1 + n2) / 2
        } else if (n1 <= n2) {
          // Normalized test score e.g. 18/20 -> 9, 0/5 -> 0 (rejected later)
          candidate = Math.round(((n1 / n2) * 10) * 100) / 100
        }
      }
    }
  } else if (clean.includes('-') && !clean.endsWith('-')) {
    // Handle hyphen ranges e.g. 7-8 -> 7.5 (excluding trailing minus e.g. '8-')
    const parts = clean.split('-')
    if (parts.length === 2) {
      const n1 = parseFloat(parts[0].replace(',', '.'))
      const n2 = parseFloat(parts[1].replace(',', '.'))
      if (!isNaN(n1) && !isNaN(n2) && n1 > 0 && n2 > 0) {
        candidate = (n1 + n2) / 2
      }
    }
  } else if (clean.endsWith('+')) {
    const base = parseFloat(clean.slice(0, -1).trim().replace(',', '.'))
    if (!isNaN(base)) candidate = base + 0.25
  } else if (clean.startsWith('+')) {
    const base = parseFloat(clean.slice(1).trim().replace(',', '.'))
    if (!isNaN(base)) candidate = base + 0.25
  } else if (clean.endsWith('-')) {
    const base = parseFloat(clean.slice(0, -1).trim().replace(',', '.'))
    if (!isNaN(base)) candidate = base - 0.25
  } else {
    candidate = parseFloat(clean.replace(',', '.'))
  }

  if (candidate === null || isNaN(candidate)) return null
  if (candidate === -1) return -1
  // Strictly enforce Italian 1-10 grading range
  if (candidate < 1.0 || candidate > 10.0) return null

  return Math.round(candidate * 100) / 100
}

export function formatGrade(val, customSeparator = null) {
  if (val === undefined || val === null || val === '' || val === '-') return '-'
  let num = Number(val)
  if (isNaN(num)) return String(val)
  if (num === -1) return 'A'

  // Clamp within bounds [-1, 10]
  if (num < -1) num = -1
  if (num > 10) num = 10
  if (num > 0 && num < 1) num = 1

  const sep = customSeparator || (typeof localStorage !== 'undefined' ? (localStorage.getItem('user_decimal_separator') || ',') : ',')
  const rounded = Math.round(num * 100) / 100
  const integerPart = Math.min(10, Math.floor(rounded))
  const decimalPart = Math.round((rounded - integerPart) * 100) / 100

  // 10 is the maximum allowed integer grade
  if (integerPart >= 10) return '10'

  // Check if rounding pushed it to the next integer
  if (decimalPart >= 0.96) {
    return `${Math.min(10, integerPart + 1)}`
  }

  if (Math.abs(decimalPart - 0.5) < 0.04) {
    return `${integerPart}½`
  }
  if (Math.abs(decimalPart - 0.25) < 0.04) {
    return `${integerPart}+`
  }
  if (Math.abs(decimalPart - 0.75) < 0.04) {
    return `${Math.min(10, integerPart + 1)}-`
  }

  if (decimalPart === 0) return `${integerPart}`

  // Single decimal digit e.g. 5.3 -> '5,3', or two decimals e.g. 5.35 -> '5,35'
  const roundedDec = Math.round(decimalPart * 10)
  if (roundedDec === 10) {
    return `${Math.min(10, integerPart + 1)}`
  }
  if (Math.abs(decimalPart * 10 - roundedDec) < 0.05) {
    return `${integerPart}${sep}${roundedDec}`
  }
  return `${integerPart}${sep}${Math.round(decimalPart * 100)}`
}

export function getGradeColor(val) {
  if (val === null || val === undefined || val === '' || val === '-') return 'white'
  const num = typeof val === 'string' ? gradeToNumeric(val) : Number(val)
  if (num === null || isNaN(num)) return 'white'
  if (num < 0) return 'grey-3'
  if (num < 5) return 'red-2'
  if (num < 6) return 'amber-2'
  if (num <= 10) return 'green-2'
  return 'white'
}
