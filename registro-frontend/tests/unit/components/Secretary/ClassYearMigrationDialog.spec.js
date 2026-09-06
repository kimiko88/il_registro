import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import ClassYearMigrationDialog from '@/components/Secretary/ClassYearMigrationDialog.vue'
import ClassScheduleDialog from '@/components/Secretary/ClassScheduleDialog.vue'

const mockNotify = vi.fn()

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      notify: mockNotify,
      dark: { isActive: false }
    })
  }
})

const mockAdminService = vi.hoisted(() => ({
  getSchoolClasses: vi.fn(),
  createClass: vi.fn()
}))
vi.mock('@/services/adminService', () => ({ default: mockAdminService }))

const mockApi = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn()
}))
vi.mock('@/services/api', () => ({ default: mockApi }))

const commonStubs = {
  'q-dialog': {
    template: '<div class="q-dialog" v-if="modelValue"><slot /></div>',
    props: ['modelValue']
  },
  'q-card': { template: '<div class="q-card"><slot /></div>' },
  'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
  'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
  'q-avatar': { template: '<div class="q-avatar"><slot /></div>' },
  'q-space': true,
  'q-icon': true,
  'q-badge': true,
  'q-spinner-dots': true,
  'q-scroll-area': { template: '<div class="q-scroll-area"><slot /></div>' },
  'q-list': { template: '<div class="q-list"><slot /></div>' },
  'q-item': { template: '<div class="q-item"><slot /></div>' },
  'q-item-section': { template: '<div class="q-item-section"><slot /></div>' },
  'q-item-label': { template: '<div class="q-item-label"><slot /></div>' },
  'q-btn': {
    template: '<button :type="type" @click="$emit(\'click\')"><slot />{{ label }}</button>',
    props: ['type', 'label', 'disabled', 'loading']
  },
  'q-btn-toggle': {
    template: '<div class="q-btn-toggle"><button v-for="opt in options" :key="opt.value" @click="$emit(\'update:modelValue\', opt.value)">{{ opt.label }}</button></div>',
    props: ['modelValue', 'options']
  },
  'q-select': {
    template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><slot /></select>',
    props: ['modelValue', 'options', 'label']
  },
  ScheduleGrid: {
    template: '<div class="schedule-grid"><button id="mock-grid-save" @click="$emit(\'save\', [{ day: 1, hour: 1 }])">Save Grid</button></div>',
    emits: ['save']
  }
}

describe('Secretary Class Dialog Components', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockAdminService.getSchoolClasses.mockResolvedValue({
      data: [
        { id: 'c2', name: '2A', section: 'A', academic_year: '2025/2026' }
      ]
    })
    mockAdminService.createClass.mockImplementation(async (payload) => ({
      data: { id: 'created-' + payload.name, ...payload }
    }))
    mockApi.get.mockResolvedValue({
      data: {
        users: [
          { id: 'st1', first_name: 'Giulia', last_name: 'Bianchi', email: 'g.b@school.it' }
        ]
      }
    })
    mockApi.post.mockResolvedValue({ data: { success: true } })
  })

  describe('ClassScheduleDialog.vue', () => {
    it('renders and forwards save event from ScheduleGrid', async () => {
      const wrapper = mount(ClassScheduleDialog, {
        props: {
          modelValue: true,
          currentClass: { id: 'class-1', name: '1', section: 'A', academic_year: '2024/2025' },
          assignments: [],
          currentSchedule: []
        },
        global: {
          plugins: [[Quasar, {}]],
          mocks: {
            $t: (key) => key
          },
          stubs: commonStubs
        }
      })

      expect(wrapper.text()).toContain('1A')
      expect(wrapper.text()).toContain('2024/2025')

      const saveBtn = wrapper.find('#mock-grid-save')
      expect(saveBtn.exists()).toBe(true)
      await saveBtn.trigger('click')

      expect(wrapper.emitted('save')).toBeTruthy()
      expect(wrapper.emitted('save')[0][0]).toEqual([{ day: 1, hour: 1 }])
    })
  })

  describe('ClassYearMigrationDialog.vue', () => {
    it('renders step 1 and advances to step 2 with loaded students', async () => {
      const wrapper = mount(ClassYearMigrationDialog, {
        props: {
          modelValue: true,
          currentYearStr: '2024/2025',
          academicYearOptions: ['2023/2024', '2024/2025', '2025/2026'],
          classes: [
            { id: 'c1', name: '1A', section: 'A', academic_year: '2024/2025' }
          ],
          schoolId: 'school-123'
        },
        global: {
          plugins: [[Quasar, {}]],
          mocks: {
            $t: (key) => key
          },
          stubs: commonStubs
        }
      })

      expect(wrapper.text()).toContain('secretaryClasses.migrationTitle')

      // Find the Next button for Step 1
      const buttons = wrapper.findAll('button')
      const nextBtn = buttons.find(b => b.text().includes('secretaryClasses.btnConfigureStudents'))
      expect(nextBtn).toBeDefined()

      await nextBtn.trigger('click')

      // Wait for promises to resolve
      await new Promise(r => setTimeout(r, 20))
      await wrapper.vm.$nextTick()

      expect(mockAdminService.getSchoolClasses).toHaveBeenCalledWith('school-123', '2025/2026')
      expect(mockApi.get).toHaveBeenCalled()
      expect(wrapper.text()).toContain('secretaryClasses.studentsToProcess')
      expect(wrapper.text()).toContain('Bianchi Giulia')
    })

    it('advances through step 3 and executes migration', async () => {
      const wrapper = mount(ClassYearMigrationDialog, {
        props: {
          modelValue: true,
          currentYearStr: '2024/2025',
          academicYearOptions: ['2023/2024', '2024/2025', '2025/2026'],
          classes: [
            { id: 'c1', name: '1A', section: 'A', academic_year: '2024/2025' }
          ],
          schoolId: 'school-123'
        },
        global: {
          plugins: [[Quasar, {}]],
          mocks: {
            $t: (key) => key
          },
          stubs: commonStubs
        }
      })

      // Go to Step 2
      const step1Btn = wrapper.findAll('button').find(b => b.text().includes('secretaryClasses.btnConfigureStudents'))
      await step1Btn.trigger('click')
      await new Promise(r => setTimeout(r, 20))
      await wrapper.vm.$nextTick()

      // Go to Step 3
      const step2Btn = wrapper.findAll('button').find(b => b.text().includes('secretaryClasses.btnVerifySummary'))
      await step2Btn.trigger('click')
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('secretaryClasses.migrationSummaryTitle')

      // Execute migration
      const execBtn = wrapper.findAll('button').find(b => b.text().includes('secretaryClasses.btnExecuteMigration'))
      await execBtn.trigger('click')
      await new Promise(r => setTimeout(r, 20))
      await wrapper.vm.$nextTick()

      expect(mockApi.post).toHaveBeenCalledWith('/classes/migrate-students', expect.objectContaining({
        source_academic_year: '2024/2025',
        target_academic_year: '2025/2026'
      }))
      expect(wrapper.emitted('migrated')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0][0]).toBe(false)
    })
  })
})
