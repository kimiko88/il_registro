import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import SchoolDetail from '@/pages/admin/SchoolDetail.vue'
import adminService from '@/services/adminService'
import { useSchoolYearStore } from '@/stores/schoolYear'

vi.mock('vue-router', () => ({
  useRoute: () => ({
    params: { id: 'school-test-123' }
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

describe('Admin SchoolDetail Page — School Management & Year Sync', () => {
  let pinia
  let schoolYearStore

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    schoolYearStore = useSchoolYearStore()
    vi.clearAllMocks()

    adminService.getSchool.mockResolvedValue({
      data: {
        id: 'school-test-123',
        name: 'Liceo Scientifico Fermi',
        code: 'RMPS01000P',
        address: 'Via Roma 10',
        city: 'Roma'
      }
    })

    adminService.getSchoolClasses.mockResolvedValue({
      data: [
        { id: 'class-1', name: '1', section: 'A', academic_year: '2024/2025', students_count: 24 }
      ]
    })

    adminService.getSchoolUsers.mockResolvedValue({
      data: { users: [] }
    })
    adminService.createClass.mockResolvedValue({ success: true })
    adminService.updateClass.mockResolvedValue({ success: true })
    adminService.deleteClass.mockResolvedValue({ success: true })
  })

  const mountComponent = () => {
    return mount(SchoolDetail, {
      global: {
        plugins: [pinia],
        stubs: {
          'q-table': {
            template: '<div class="q-table"><slot /></div>',
            props: ['rows', 'columns', 'loading', 'filter']
          },
          'q-tabs': true,
          'q-tab': true,
          'q-tab-panels': true,
          'q-tab-panel': true
        }
      }
    })
  }

  it('loads school info and classes on mount', async () => {
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 20))

    expect(adminService.getSchool).toHaveBeenCalledWith('school-test-123')
    expect(adminService.getSchoolClasses).toHaveBeenCalledWith('school-test-123', expect.anything())
    expect(wrapper.vm.school.name).toBe('Liceo Scientifico Fermi')
    expect(wrapper.vm.classes.length).toBe(1)
  })

  it('initializes new class form with schoolYearStore.selectedSchoolYear', async () => {
    schoolYearStore.setSchoolYear('2024/2025')
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 50))

    wrapper.vm.openClassDialog()
    expect(wrapper.vm.showClassDialog).toBe(true)
    expect(wrapper.vm.classForm.academic_year).toBe('2024/2025')
    expect(wrapper.vm.classForm.school_id).toBe('school-test-123')
  })

  it('saves new class via adminService.createClass', async () => {
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 50))

    wrapper.vm.openClassDialog()
    wrapper.vm.classForm.name = '3'
    wrapper.vm.classForm.section = 'B'
    wrapper.vm.classForm.academic_year = '2024/2025'

    await wrapper.vm.saveClass()
    expect(adminService.createClass).toHaveBeenCalledWith(
      expect.objectContaining({
        name: '3',
        section: 'B',
        academic_year: '2024/2025',
        school_id: 'school-test-123'
      })
    )
    expect(wrapper.vm.showClassDialog).toBe(false)
  })

  it('updates existing class via adminService.updateClass', async () => {
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 10))

    const existingClass = { id: 'class-1', name: '1', section: 'A', academic_year: '2024/2025' }
    wrapper.vm.editClass(existingClass)
    expect(wrapper.vm.editingClass).toEqual(existingClass)

    wrapper.vm.classForm.section = 'C'
    await wrapper.vm.saveClass()
    expect(adminService.updateClass).toHaveBeenCalledWith(
      'class-1',
      expect.objectContaining({ section: 'C' })
    )
  })

  it('deletes class via confirmDeleteClass', async () => {
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 10))

    const existingClass = { id: 'class-1', name: '1', section: 'A' }
    wrapper.vm.confirmDeleteClass(existingClass)
    await wrapper.vm.$nextTick()

    expect(adminService.deleteClass).toHaveBeenCalledWith('class-1')
  })
})
