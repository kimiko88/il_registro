import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SchoolManagement from '@/pages/admin/SchoolManagement.vue'
import adminService from '@/services/adminService'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

// Mock adminService
vi.mock('@/services/adminService', () => ({
  default: {
    getSchools: vi.fn(),
    deleteSchool: vi.fn(),
    createSchool: vi.fn(),
    updateSchool: vi.fn(),
    getTiers: vi.fn()
  }
}))

// Mock permissions
vi.mock('@/composables/usePermissions', () => ({
  usePermissions: () => ({
    isSuperAdmin: true,
    canCreateSchools: true,
    canDeleteSchools: true,
    canEditSchool: () => true
  })
}))

describe('SchoolManagement - School Tiers & Normative Levels', () => {
  let wrapper

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()

    adminService.getSchools.mockResolvedValue({
      data: {
        items: [
          {
            id: 'school-infanzia-1',
            name: 'Scuola Infanzia Munari',
            code: 'RMAA001',
            school_level: 'infanzia',
            city: 'Roma',
            province: 'RM',
            student_count: 50,
            teacher_count: 5,
            is_active: true
          },
          {
            id: 'school-primaria-1',
            name: 'Scuola Primaria Rodari',
            code: 'RMEE001',
            school_level: 'primaria',
            city: 'Roma',
            province: 'RM',
            student_count: 220,
            teacher_count: 18,
            is_active: true
          },
          {
            id: 'school-sec2-1',
            name: 'Liceo Scientifico Righi',
            code: 'RMPS001',
            school_level: 'secondaria_secondo_grado',
            city: 'Roma',
            province: 'RM',
            student_count: 850,
            teacher_count: 65,
            is_active: true
          }
        ],
        total: 3
      }
    })

    wrapper = mount(SchoolManagement, {
      global: {
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-table': {
            template: `
              <div class="q-table-stub">
                <slot name="body-cell-name" :row="{ name: 'Liceo Scientifico Righi', code: 'RMPS001' }" />
                <slot name="body-cell-school_level" :row="{ school_level: 'infanzia', type: 'infanzia' }" />
                <slot name="body-cell-school_level" :row="{ school_level: 'primaria', type: 'primaria' }" />
                <slot name="body-cell-school_level" :row="{ school_level: 'secondaria_secondo_grado', type: 'secondaria_secondo_grado' }" />
              </div>
            `
          },
          'q-badge': {
            template: '<span class="q-badge" :data-color="color"><slot /></span>',
            props: ['color']
          },
          'q-btn': {
            template: '<button @click="$emit(\'click\', $event)">{{ label }}<slot /></button>',
            props: ['label']
          },
          'q-select': {
            template: '<div class="q-select-stub"><slot /></div>',
            props: ['modelValue', 'options', 'label', 'hint']
          },
          'q-input': {
            template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
            props: ['modelValue']
          },
          'q-dialog': { template: '<div><slot /></div>' },
          'q-toggle': true,
          'q-icon': true,
          'q-tooltip': true,
          'q-space': true
        }
      }
    })
  })

  it('populates tier filter options and tier select options correctly', () => {
    const tierFilter = wrapper.vm.tierFilterOptions
    expect(tierFilter).toBeDefined()
    expect(tierFilter.length).toBeGreaterThanOrEqual(6)

    const tierValues = tierFilter.map(o => o.value)
    expect(tierValues).toContain('infanzia')
    expect(tierValues).toContain('primaria')
    expect(tierValues).toContain('secondaria_primo_grado')
    expect(tierValues).toContain('secondaria_secondo_grado')
    expect(tierValues).toContain('comprensivo')
    expect(tierValues).toContain('omnicomprensivo')
  })

  it('maps correct badge colors for each school tier', () => {
    expect(wrapper.vm.getSchoolTierBadgeColor('infanzia')).toBe('pink-7')
    expect(wrapper.vm.getSchoolTierBadgeColor('primaria')).toBe('amber-9')
    expect(wrapper.vm.getSchoolTierBadgeColor('secondaria_primo_grado')).toBe('blue-8')
    expect(wrapper.vm.getSchoolTierBadgeColor('secondaria_secondo_grado')).toBe('indigo-8')
    expect(wrapper.vm.getSchoolTierBadgeColor('comprensivo')).toBe('teal-8')
    expect(wrapper.vm.getSchoolTierBadgeColor('omnicomprensivo')).toBe('purple-8')
  })

  it('returns human-readable labels for school tiers', () => {
    expect(wrapper.vm.getSchoolTierLabel('infanzia')).toBeTruthy()
    expect(wrapper.vm.getSchoolTierLabel('primaria')).toBeTruthy()
    expect(wrapper.vm.getSchoolTierLabel('secondaria_secondo_grado')).toBeTruthy()
  })

  it('initializes schoolForm with default tier secondaria_secondo_grado on openCreate', () => {
    wrapper.vm.openCreate()
    expect(wrapper.vm.showCreateDialog).toBe(true)
    expect(wrapper.vm.editingSchool).toBeNull()
    expect(wrapper.vm.schoolForm.school_level).toBe('secondaria_secondo_grado')
  })

  it('populates school_level from existing school on editSchool', () => {
    const schoolToEdit = {
      id: 'sc-primaria',
      name: 'Scuola Primaria Rodari',
      code: 'RMEE001',
      school_level: 'primaria',
      address: 'Via Roma 10',
      city: 'Roma',
      province: 'RM',
      zip_code: '00100',
      phone: '06123456',
      email: 'rmee001@istruzione.it',
      is_active: true
    }

    wrapper.vm.editSchool(schoolToEdit)
    expect(wrapper.vm.showCreateDialog).toBe(true)
    expect(wrapper.vm.editingSchool).toEqual(schoolToEdit)
    expect(wrapper.vm.schoolForm.school_level).toBe('primaria')
  })

  it('sends school_level payload when creating a new school', async () => {
    adminService.createSchool.mockResolvedValue({ data: { id: 'new-school-1' } })

    wrapper.vm.openCreate()
    wrapper.vm.schoolForm.name = 'Nuova Primaria Test'
    wrapper.vm.schoolForm.code = 'RMIC0099'
    wrapper.vm.schoolForm.school_level = 'primaria'
    wrapper.vm.schoolForm.address = 'Via Nuova 1'
    wrapper.vm.schoolForm.city = 'Roma'
    wrapper.vm.schoolForm.province = 'RM'
    wrapper.vm.schoolForm.zip_code = '00100'

    await wrapper.vm.saveSchool()

    expect(adminService.createSchool).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'Nuova Primaria Test',
        code: 'RMIC0099',
        school_level: 'primaria'
      })
    )
  })

  it('sends school_level payload when updating an existing school', async () => {
    adminService.updateSchool.mockResolvedValue({ data: { message: 'ok' } })

    wrapper.vm.editSchool({
      id: 'sc-edit-99',
      name: 'Liceo Scientifico',
      code: 'RMPS99',
      school_level: 'secondaria_secondo_grado',
      address: 'Via Cavour 1',
      city: 'Roma',
      province: 'RM',
      zip_code: '00100',
      is_active: true
    })

    wrapper.vm.schoolForm.school_level = 'omnicomprensivo'

    await wrapper.vm.saveSchool()

    expect(adminService.updateSchool).toHaveBeenCalledWith(
      'sc-edit-99',
      expect.objectContaining({
        school_level: 'omnicomprensivo'
      })
    )
  })
})
