import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import Classes from '@/pages/secretary/Classes.vue'
import Scrutiny from '@/pages/secretary/Scrutiny.vue'
import SchoolDetail from '@/pages/admin/SchoolDetail.vue'
import { useSchoolYearStore } from '@/stores/schoolYear'
import { useClassesStore } from '@/stores/classes'
import { useScrutinyStore } from '@/stores/scrutiny'
import { useAuthStore } from '@/stores/auth'
import adminService from '@/services/adminService'

vi.mock('vue-router', () => ({
  useRoute: () => ({
    params: { id: 'school-e2e-1' }
  }),
  useRouter: () => ({
    push: vi.fn()
  })
}))

vi.mock('@/services/adminService', () => ({
  default: {
    getSchool: vi.fn(),
    getSchoolClasses: vi.fn(),
    getSchoolUsers: vi.fn(),
    createClass: vi.fn(),
    updateClass: vi.fn(),
    deleteClass: vi.fn(),
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: [] }),
    post: vi.fn().mockResolvedValue({ data: { success: true } })
  }
}))

describe('E2E Workflow: Centralized School Year Switching across Secretary & Admin Pages', () => {
  let pinia
  let schoolYearStore
  let classesStore
  let scrutinyStore

  beforeEach(() => {
    vi.clearAllMocks()
    pinia = createPinia()
    setActivePinia(pinia)

    schoolYearStore = useSchoolYearStore()
    classesStore = useClassesStore()
    scrutinyStore = useScrutinyStore()
    const authStore = useAuthStore()
    authStore.user = { id: 'u-1', school_id: 1 }

    // Setup initial store values
    schoolYearStore.initializeForUser({ created_at: '2022-09-01T00:00:00Z' })
    schoolYearStore.setSchoolYear('2024/2025')

    classesStore.fetchClasses = vi.fn().mockResolvedValue()
    classesStore.classes = [
      { id: 'c-1', name: '1', section: 'A', academic_year: '2024/2025' }
    ]

    scrutinyStore.fetchOverview = vi.fn().mockResolvedValue()
    scrutinyStore.overview = [
      { class_id: 'c-1', class_name: '1A', completed_subjects: 8, total_subjects: 8, status: 'completed' }
    ]

    adminService.getSchool.mockResolvedValue({
      data: { id: 'school-e2e-1', name: 'Istituto Comprensivo Volta' }
    })
    adminService.getSchoolClasses.mockResolvedValue({
      data: [{ id: 'c-1', name: '1', section: 'A', academic_year: '2024/2025' }]
    })
    adminService.getSchoolUsers.mockResolvedValue({
      data: { users: [] }
    })
    adminService.createClass.mockResolvedValue({ success: true })
  })

  it('Step 1: Secretary Classes page uses globally selected school year', async () => {
    expect(schoolYearStore.selectedSchoolYear).toBe('2024/2025')

    const wrapper = mount(Classes, {
      global: {
        plugins: [pinia],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-table': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-select': { template: '<div><slot /></div>' },
          'q-input': { template: '<input />' },
          'q-icon': true,
          'q-badge': true,
          'q-dialog': { template: '<div v-if="modelValue"><slot /></div>', props: ['modelValue'] }
        }
      }
    })

    await flushPromises()
    expect(wrapper.exists()).toBe(true)
    expect(classesStore.fetchClasses).toHaveBeenCalledWith(
      expect.objectContaining({ school_id: expect.anything() }),
      expect.anything()
    )
    wrapper.unmount()
  })

  it('Step 2: Switching school year reactively triggers Scrutiny overview reload', async () => {
    const wrapper = mount(Scrutiny, {
      global: {
        plugins: [pinia],
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-table': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-input': { template: '<input />' },
          'q-icon': true,
          'q-badge': true,
          'q-dialog': true
        }
      }
    })

    await flushPromises()
    expect(scrutinyStore.fetchOverview).toHaveBeenCalledTimes(1)

    // Switch to 2023/2024
    schoolYearStore.setSchoolYear('2023/2024')
    await flushPromises()
    await new Promise(resolve => setTimeout(resolve, 10))

    expect(scrutinyStore.fetchOverview).toHaveBeenCalledTimes(2)
    expect(schoolYearStore.selectedSchoolYear).toBe('2023/2024')
    expect(wrapper.exists()).toBe(true)
    wrapper.unmount()
  })

  it('Step 3: Admin SchoolDetail form presets new class academic_year from active global year', async () => {
    schoolYearStore.setSchoolYear('2025/2026')

    const wrapper = mount(SchoolDetail, {
      global: {
        plugins: [pinia],
        stubs: {
          'q-table': { template: '<div><slot /></div>' },
          'q-tabs': true,
          'q-tab': true,
          'q-tab-panels': true,
          'q-tab-panel': true
        }
      }
    })

    await flushPromises()
    await new Promise(resolve => setTimeout(resolve, 50))

    wrapper.vm.openClassDialog()
    expect(wrapper.vm.showClassDialog).toBe(true)
    expect(wrapper.vm.classForm.academic_year).toBe('2025/2026')

    wrapper.vm.classForm.name = '5'
    wrapper.vm.classForm.section = 'D'
    await wrapper.vm.saveClass()

    expect(adminService.createClass).toHaveBeenCalledWith(
      expect.objectContaining({
        name: '5',
        section: 'D',
        academic_year: '2025/2026',
        school_id: 'school-e2e-1'
      })
    )
    wrapper.unmount()
  })
})
