import { defineStore } from 'pinia'

function getCurrentSchoolYear() {
  const now = new Date()
  const year = now.getFullYear()
  const month = now.getMonth() // 0 = Jan, 8 = Sep (September onwards)
  if (month >= 8) {
    return `${year}/${year + 1}`
  } else {
    return `${year - 1}/${year}`
  }
}

function getStartYearFromStr(syStr) {
  if (!syStr) return null
  const parts = String(syStr).replace('-', '/').split('/')
  const val = parseInt(parts[0], 10)
  return isNaN(val) ? null : val
}

function getSchoolYearFromDate(dateInput) {
  if (!dateInput) return getCurrentSchoolYear()
  const d = new Date(dateInput)
  if (isNaN(d.getTime())) return getCurrentSchoolYear()
  const year = d.getFullYear()
  const month = d.getMonth()
  const startYear = month >= 8 ? year : year - 1
  return `${startYear}/${startYear + 1}`
}

function generateSchoolYearsForUser(createdAt) {
  const currentSY = getCurrentSchoolYear()
  const currentStartYear = getStartYearFromStr(currentSY)

  let regStartYear = currentStartYear - 1
  if (createdAt) {
    const regSY = getSchoolYearFromDate(createdAt)
    const parsed = getStartYearFromStr(regSY)
    if (parsed && parsed <= currentStartYear) {
      regStartYear = parsed
    }
  }

  // Generate list starting from the registration school year up to current active school year
  const years = []
  for (let y = currentStartYear; y >= regStartYear; y--) {
    years.push(`${y}/${y + 1}`)
  }

  return years
}

export const useSchoolYearStore = defineStore('schoolYear', {
  state: () => {
    const saved = localStorage.getItem('selected_school_year')
    const current = getCurrentSchoolYear()
    return {
      selectedSchoolYear: saved || current,
      availableSchoolYears: generateSchoolYearsForUser(null)
    }
  },
  actions: {
    initializeForUser(user) {
      const createdAt = user?.created_at || user?.registration_date || null
      const years = generateSchoolYearsForUser(createdAt)
      this.availableSchoolYears = years

      // Default selected school year to current or saved, ensuring it exists in available options
      const saved = localStorage.getItem('selected_school_year')
      if (saved && years.includes(saved)) {
        this.selectedSchoolYear = saved
      } else if (!years.includes(this.selectedSchoolYear)) {
        this.selectedSchoolYear = years[0] || getCurrentSchoolYear()
        localStorage.setItem('selected_school_year', this.selectedSchoolYear)
      }
    },
    setSchoolYear(year) {
      if (year) {
        this.selectedSchoolYear = year
        localStorage.setItem('selected_school_year', year)
      }
    }
  }
})
