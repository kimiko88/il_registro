import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import Scrutiny from '@/pages/secretary/Scrutiny.vue'
import { useScrutinyStore } from '@/stores/scrutiny'
import { useSchoolYearStore } from '@/stores/schoolYear'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key) => key,
    locale: { value: 'it' }
  }),
  createI18n: () => ({
    global: {
      t: (key) => key,
      locale: { value: 'it' }
    }
  })
}))

describe('Secretary Scrutiny Page — Supervision & Global School Year Reactivity', () => {
  let pinia
  let scrutinyStore
  let schoolYearStore

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    scrutinyStore = useScrutinyStore()
    schoolYearStore = useSchoolYearStore()

    // Mock store methods
    scrutinyStore.fetchOverview = vi.fn().mockResolvedValue()
    scrutinyStore.fetchClassReport = vi.fn().mockResolvedValue({
      class_name: '3B',
      students: [
        { id: 'stud-1', name: 'Mario Rossi', final_grade: 7.5, outcome: 'admitted' }
      ]
    })
    scrutinyStore.finalizeScrutiny = vi.fn().mockResolvedValue({ success: true })
    scrutinyStore.exportAll = vi.fn()

    scrutinyStore.overview = [
      { class_id: 'c-1', class_name: '1A', completed_subjects: 8, total_subjects: 8, status: 'completed' },
      { class_id: 'c-2', class_name: '2A', completed_subjects: 5, total_subjects: 8, status: 'in_progress' },
      { class_id: 'c-3', class_name: '3B', completed_subjects: 0, total_subjects: 8, status: 'pending' }
    ]
  })

  const mountComponent = () => {
    return mount(Scrutiny, {
      global: {
        plugins: [pinia],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-table': {
            template: `
              <div class="q-table">
                <div v-for="row in (rows || [])" :key="row.class_id">
                  <slot name="body-cell-status" :row="row" :props="{ row }" />
                  <slot name="body-cell-completed_subjects" :row="row" :props="{ row }" />
                  <slot name="body-cell-actions" :row="row" :props="{ row }" />
                </div>
                <slot />
              </div>
            `,
            props: ['rows', 'columns']
          },
          'q-td': { template: '<td><slot /></td>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-input': { template: '<input />' },
          'q-icon': { template: '<i></i>' },
          'q-badge': { template: '<span><slot /></span>' },
          'q-dialog': { template: '<div v-if="modelValue"><slot /></div>', props: ['modelValue'] },
          'q-linear-progress': { template: '<div></div>' }
        },
        mocks: {
          $q: {
            dark: { isActive: false },
            notify: vi.fn()
          }
        }
      }
    })
  }

  it('renders scrutiny supervision page and loads overview on mount', () => {
    const wrapper = mountComponent()
    expect(wrapper.exists()).toBe(true)
    expect(scrutinyStore.fetchOverview).toHaveBeenCalledTimes(1)
  })

  it('filters class overview based on search filterText', async () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.filteredOverview.length).toBe(3)

    wrapper.vm.filterText = '3b'
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.filteredOverview.length).toBe(1)
    expect(wrapper.vm.filteredOverview[0].class_name).toBe('3B')
  })

  it('reactively reloads scrutiny overview when selectedSchoolYear changes in schoolYearStore', async () => {
    mountComponent()
    expect(scrutinyStore.fetchOverview).toHaveBeenCalledTimes(1)

    // Simulate global academic year change
    schoolYearStore.setSchoolYear('2023/2024')
    await new Promise(resolve => setTimeout(resolve, 10))

    expect(scrutinyStore.fetchOverview).toHaveBeenCalledTimes(2)
  })

  it('opens scrutiny detail modal for selected class', async () => {
    const wrapper = mountComponent()
    const targetClass = { class_id: 'c-2', class_name: '2A' }

    await wrapper.vm.openDetail(targetClass)
    expect(scrutinyStore.fetchClassReport).toHaveBeenCalledWith('c-2')
    expect(wrapper.vm.selectedClass).toEqual(targetClass)
    expect(wrapper.vm.detailDialog).toBe(true)
  })

  it('opens confirmation modal and finalizes scrutiny for class in progress', async () => {
    const wrapper = mountComponent()
    const targetClass = { class_id: 'c-2', class_name: '2A', status: 'in_progress' }

    wrapper.vm.confirmFinalize(targetClass)
    expect(wrapper.vm.finalizeModal).toBe(true)
    expect(wrapper.vm.classToFinalize).toEqual(targetClass)

    await wrapper.vm.executeFinalize()
    expect(scrutinyStore.finalizeScrutiny).toHaveBeenCalledWith('c-2')
    expect(wrapper.vm.finalizeModal).toBe(false)
  })

  it('calls exportAll on export click', () => {
    const wrapper = mountComponent()
    wrapper.vm.exportAll()
    expect(scrutinyStore.exportAll).toHaveBeenCalled()
  })
})
