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
  if (typeof val === 'number') return val
  const clean = String(val).trim().toUpperCase()
  if (clean === 'A' || clean === 'ASSENTE') return -1

  if (clean.includes('/')) {
    const parts = clean.split('/')
    if (parts.length === 2) {
      const n1 = parseFloat(parts[0])
      const n2 = parseFloat(parts[1])
      if (!isNaN(n1) && !isNaN(n2)) {
        return (n1 + n2) / 2
      }
    }
  }

  if (clean.endsWith('1/2') || clean.endsWith('½')) {
    const base = parseFloat(clean.replace('1/2', '').replace('½', '').trim())
    if (!isNaN(base)) return base + 0.5
  }

  if (clean.endsWith('+')) {
    const base = parseFloat(clean.slice(0, -1).trim())
    if (!isNaN(base)) return base + 0.25
  }

  if (clean.endsWith('-')) {
    const base = parseFloat(clean.slice(0, -1).trim())
    if (!isNaN(base)) return base - 0.25
  }

  const numericVal = parseFloat(clean.replace(',', '.'))
  return isNaN(numericVal) ? null : numericVal
}

export function formatGrade(val) {
  if (val === undefined || val === null || val === '' || val === '-') return '-'
  const num = Number(val)
  if (isNaN(num)) return val
  if (num === -1) return 'A'

  const sep = typeof localStorage !== 'undefined' ? (localStorage.getItem('user_decimal_separator') || ',') : ','
  const integerPart = Math.floor(num)
  const decimalPart = Math.round((num - integerPart) * 100) / 100

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
  return `${integerPart}${sep}${Math.round(decimalPart * 10)}`
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
