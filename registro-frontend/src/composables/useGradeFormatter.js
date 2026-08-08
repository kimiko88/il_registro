import { ref } from 'vue'

/**
 * Composable per la formattazione dei voti e dei numeri decimali
 * in base al separatore decimali scelto dall'utente (virgola ',' o punto '.')
 * e al formato grafico (decimali, frazioni, centesimi).
 */
export function useGradeFormatter() {
  const getSeparator = () => {
    return localStorage.getItem('user_decimal_separator') || ','
  }

  const getFormat = () => {
    try {
      const reg = JSON.parse(localStorage.getItem('teacher_register_settings') || '{}')
      return reg.gradeFormat || 'decimal'
    } catch {
      return 'decimal'
    }
  }

  /**
   * Formatta un singolo voto numerico per l'interfaccia utente.
   * @param {number|string} val - Valore del voto (es. 7.5)
   * @param {string|null} customSeparator - Separatore custom facoltativo (',' o '.')
   * @returns {string} Voto formattato (es. "7,5" oppure "7.5" oppure "7½")
   */
  function formatGrade(val, customSeparator = null) {
    if (val === undefined || val === null || val === '') return '-'
    const num = Number(val)
    if (isNaN(num)) return String(val)
    if (num === -1) return 'A' // Assente

    const separator = customSeparator || getSeparator()
    const format = getFormat()

    if (format === 'fractional') {
      const integerPart = Math.floor(num)
      const decimalPart = num - integerPart

      if (Math.abs(decimalPart - 0.5) < 0.01) {
        return `${integerPart}½`
      }
      if (Math.abs(decimalPart - 0.25) < 0.01) {
        return `${integerPart}+`
      }
      if (Math.abs(decimalPart - 0.75) < 0.01) {
        return `${integerPart + 1}-`
      }
      if (Math.abs(decimalPart) < 0.01) {
        return `${integerPart}`
      }
    } else if (format === 'centesimal') {
      const centesimalVal = (num * 10).toFixed(0)
      return `${centesimalVal}/100`
    }

    // Formato decimali standard
    const str = String(num)
    return str.replace('.', separator)
  }

  /**
   * Formatta una media o valore decimale a N cifre decimali.
   * @param {number|string} num - Numero da formattare
   * @param {number} decimals - Cifre decimali (default 2)
   * @param {string|null} customSeparator - Separatore custom (',' o '.')
   */
  function formatDecimal(num, decimals = 2, customSeparator = null) {
    if (num === undefined || num === null || isNaN(Number(num))) return '-'
    const separator = customSeparator || getSeparator()
    const fixed = Number(num).toFixed(decimals)
    // Rimuove gli zeri decimali non necessari se interi (opzionale)
    const floatVal = parseFloat(fixed)
    const formattedStr = floatVal % 1 === 0 ? String(floatVal) : fixed
    return formattedStr.replace('.', separator)
  }

  return {
    formatGrade,
    formatDecimal,
    getSeparator,
    getFormat
  }
}
