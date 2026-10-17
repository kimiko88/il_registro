import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import Classes from '@/pages/teacher/Classes.vue'
import { useSchoolYearStore } from '@/stores/schoolYear'
import api from '@/services/api'

import { Quasar } from 'quasar'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

vi.mock('@/services/notesService', () => ({
  default: {
    getNotes: vi.fn().mockResolvedValue({ data: [] }),
    deleteNote: vi.fn().mockResolvedValue({})
  }
}))

describe('Teacher Classes.vue — Academic Year Filtering & Class Display', () => {
  let pinia
  let schoolYearStore

  const mockClasses = [
    { id: 'c-2025', name: '3A', section: 'A', academic_year: '2025/2026', students: 20 },
    { id: 'c-2024', name: '2A', section: 'A', academic_year: '2024/2025', students: 18 }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    pinia = createPinia()
    setActivePinia(pinia)
    schoolYearStore = useSchoolYearStore()

    schoolYearStore.selectedSchoolYear = '2025/2026'
    schoolYearStore.availableSchoolYears = ['2025/2026', '2024/2025']

    api.get.mockImplementation((url, config) => {
      if (url === '/teacher/classes') {
        const sy = config?.params?.school_year || schoolYearStore.selectedSchoolYear
        const filtered = mockClasses.filter(c => c.academic_year === sy)
        return Promise.resolve({ data: filtered })
      }
      if (url === '/users') {
        return Promise.resolve({ data: { users: [{ id: 's1', first_name: 'Mario', last_name: 'Rossi', email: 'mario@school.it' }] } })
      }
      return Promise.resolve({ data: [] })
    })
  })

  function createWrapper() {
    return mount(Classes, {
      global: {
        plugins: [Quasar, pinia],
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-toolbar': { template: '<div><slot /></div>' },
          'q-toolbar-title': { template: '<div><slot /></div>' },
          'q-tabs': { template: '<div><slot /></div>' },
          'q-tab': { template: '<button><slot /></button>' },
          'q-tab-panels': { template: '<div><slot /></div>' },
          'q-tab-panel': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-icon': true,
          'q-avatar': { template: '<div><slot /></div>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-list': { template: '<div class="q-list"><slot /></div>' },
          'q-item': { template: '<div class="q-item" @click="$emit(\'click\')"><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-select': {
            props: ['modelValue', 'options'],
            template: '<div class="q-select-stub"><slot name="option" v-for="opt in options" :opt="opt" :itemProps="{}" /></div>'
          },
          'q-table': { template: '<div class="q-table-stub"><slot /></div>' },
          'q-input': true,
          'q-tooltip': true,
          'NoteDialog': true
        }
      }
    })
  }

  it('loads and displays only classes belonging to selected school year', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    expect(wrapper.text()).toContain('3A')
    expect(wrapper.text()).not.toContain('2A')
  })

  it('updates classes when selectedSchoolYear changes in store', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    expect(wrapper.text()).toContain('3A')

    // Switch to 2024/2025
    schoolYearStore.selectedSchoolYear = '2024/2025'
    await flushPromises()

    expect(wrapper.text()).toContain('2A')
    expect(wrapper.text()).not.toContain('3A')
  })

  it('shows empty state when no classes match the chosen school year', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    schoolYearStore.selectedSchoolYear = '2022/2023'
    await flushPromises()

    expect(wrapper.text()).toMatch(/Nessuna classe|noClassesFound/)
  })
})
