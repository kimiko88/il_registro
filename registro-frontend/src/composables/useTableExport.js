import { exportFile, useQuasar } from 'quasar'

/**
 * Sanitizes values to prevent CSV formula injection (DDE attacks)
 * Prefixes cells starting with '=', '+', '-', '@', '\t', '\r' with an apostrophe.
 * Also escapes double quotes.
 */
export function sanitizeCsvFormula(val) {
  if (val === null || val === undefined) return ''
  let str = String(val)
  if (/^[=+\-@\t\r]/.test(str)) {
    str = "'" + str
  }
  return str.replace(/"/g, '""')
}

/**
 * Composable for standardized and safe CSV table export across Quasar tables.
 */
export function useTableExport() {
  const $q = useQuasar()

  /**
   * Exports tabular data to CSV with automatic header, formatter, and formula injection sanitization.
   * @param {Object} options
   * @param {string} [options.filename='export.csv'] - File name for the download
   * @param {Array} [options.columns=[]] - Quasar table columns definition
   * @param {Array} [options.rows=[]] - Data rows
   * @param {string} [options.delimiter=','] - CSV delimiter (e.g. ',' or ';')
   * @param {string} [options.mimeType='text/csv;charset=utf-8;'] - MIME type
   * @returns {boolean} Whether the export succeeded
   */
  function exportTableCsv({
    filename = 'export.csv',
    columns = [],
    rows = [],
    delimiter = ',',
    mimeType = 'text/csv;charset=utf-8;'
  } = {}) {
    if (!rows || rows.length === 0) {
      $q.notify({ type: 'warning', message: 'Nessun dato da esportare.' })
      return false
    }

    const header = columns
      .map(col => `"${sanitizeCsvFormula(col.label !== undefined ? col.label : col.name)}"`)
      .join(delimiter)

    const contentRows = rows.map(row => {
      return columns
        .map(col => {
          let val
          if (typeof col.field === 'function') {
            val = col.field(row)
          } else if (col.field !== undefined) {
            val = row[col.field]
          } else {
            val = row[col.name]
          }

          if (col.format && typeof col.format === 'function') {
            val = col.format(val, row)
          }

          return `"${sanitizeCsvFormula(val)}"`
        })
        .join(delimiter)
    })

    const content = [header, ...contentRows].join('\r\n')
    const status = exportFile(filename, content, mimeType)

    if (status !== true) {
      $q.notify({
        type: 'negative',
        message: 'Impossibile scaricare il file. Verifica i permessi del browser.'
      })
      return false
    }

    $q.notify({
      type: 'positive',
      message: 'Esportazione completata con successo',
      icon: 'download'
    })
    return true
  }

  return {
    exportTableCsv,
    sanitizeCsvFormula
  }
}

export default useTableExport
