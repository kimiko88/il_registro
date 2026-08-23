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
  if (typeof val === 'number') return isNaN(val) ? null : val
  const clean = String(val).trim().toUpperCase()
  if (clean === 'A' || clean === 'ASSENTE' || clean === 'ABSENT') return -1
  if (clean === 'NC' || clean === 'NON CLASSIFICATO') return null
  if (clean === 'E' || clean === 'ESENTE') return null
  if (clean === 'S' || clean === 'SUFF' || clean === 'SUFFICIENTE') return 6
  if (clean === 'INS' || clean === 'INSUFFICIENTE') return 5
  if (clean === 'O' || clean === 'OTTIMO') return 10
  if (clean === 'D' || clean === 'DISTINTO') return 8
  if (clean === 'B' || clean === 'BUONO') return 7

  // Handle half notation first (e.g. 8 1/2 or 9½) before general / fraction split
  if (clean.endsWith('1/2') || clean.endsWith('½')) {
    const base = parseFloat(clean.replace('1/2', '').replace('½', '').trim().replace(',', '.'))
    if (!isNaN(base)) return base + 0.5
  }

  // Handle ratio scores like 15/20 or grade ranges like 7/8
  if (clean.includes('/')) {
    const parts = clean.split('/')
    if (parts.length === 2) {
      const n1 = parseFloat(parts[0].replace(',', '.'))
      const n2 = parseFloat(parts[1].replace(',', '.'))
      if (!isNaN(n1) && !isNaN(n2) && n2 > 0) {
        if (n2 > 10 && n1 <= n2) {
          // Normalized test score e.g. 18/20 -> 9
          return Math.round(((n1 / n2) * 10) * 100) / 100
        }
        // Grade range e.g. 7/8 -> 7.5
        return (n1 + n2) / 2
      }
    }
  }

  // Handle hyphen ranges e.g. 7-8 -> 7.5 (excluding trailing minus e.g. '8-')
  if (clean.includes('-') && !clean.endsWith('-')) {
    const parts = clean.split('-')
    if (parts.length === 2) {
      const n1 = parseFloat(parts[0].replace(',', '.'))
      const n2 = parseFloat(parts[1].replace(',', '.'))
      if (!isNaN(n1) && !isNaN(n2) && n1 > 0 && n2 > 0) {
        return (n1 + n2) / 2
      }
    }
  }

  if (clean.endsWith('+')) {
    const base = parseFloat(clean.slice(0, -1).trim().replace(',', '.'))
    if (!isNaN(base)) return base + 0.25
  }

  if (clean.startsWith('+')) {
    const base = parseFloat(clean.slice(1).trim().replace(',', '.'))
    if (!isNaN(base)) return base + 0.25
  }

  if (clean.endsWith('-')) {
    const base = parseFloat(clean.slice(0, -1).trim().replace(',', '.'))
    if (!isNaN(base)) return base - 0.25
  }

  const numericVal = parseFloat(clean.replace(',', '.'))
  return isNaN(numericVal) ? null : numericVal
}

export function formatGrade(val, customSeparator = null) {
  if (val === undefined || val === null || val === '' || val === '-') return '-'
  const num = Number(val)
  if (isNaN(num)) return String(val)
  if (num === -1) return 'A'

  const sep = customSeparator || (typeof localStorage !== 'undefined' ? (localStorage.getItem('user_decimal_separator') || ',') : ',')
  const rounded = Math.round(num * 100) / 100
  const integerPart = Math.floor(rounded)
  const decimalPart = Math.round((rounded - integerPart) * 100) / 100

  // Check if rounding pushed it to the next integer
  if (decimalPart >= 0.96) {
    return `${integerPart + 1}`
  }

  if (Math.abs(decimalPart - 0.5) < 0.04) {
    return `${integerPart}½`
  }
  if (Math.abs(decimalPart - 0.25) < 0.04) {
    return `${integerPart}+`
  }
  if (Math.abs(decimalPart - 0.75) < 0.04) {
    return `${integerPart + 1}-`
  }

  if (decimalPart === 0) return `${integerPart}`

  // Single decimal digit e.g. 5.3 -> '5,3', or two decimals e.g. 5.35 -> '5,35'
  const roundedDec = Math.round(decimalPart * 10)
  if (roundedDec === 10) {
    return `${integerPart + 1}`
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
  return 'green-2'
}
