import { defineStore } from 'pinia'

function getCurrentSchoolYear() {
  const now = new Date()
  const year = now.getFullYear()
  const month = now.getMonth() // 0 = Jan, 8 = Sep
  if (month >= 8) {
    return `${year}/${year + 1}`
  } else {
    return `${year - 1}/${year}`
  }
}

function generateSchoolYears() {
  const current = getCurrentSchoolYear()
  const startYear = parseInt(current.split('/')[0])
  const years = []
  for (let y = startYear - 2; y <= startYear + 1; y++) {
    years.push(`${y}/${y + 1}`)
  }
  return years
}

export const useSchoolYearStore = defineStore('schoolYear', {
  state: () => ({
    selectedSchoolYear: getCurrentSchoolYear(),
    availableSchoolYears: generateSchoolYears()
  }),
  actions: {
    setSchoolYear(year) {
      if (year) {
        this.selectedSchoolYear = year
      }
    }
  }
})
