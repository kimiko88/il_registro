import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import LessonFormDialog from '@/components/Teacher/LessonFormDialog.vue'
import HomeworkFormDialog from '@/components/Teacher/HomeworkFormDialog.vue'
import FreeActivityDialog from '@/components/Teacher/FreeActivityDialog.vue'

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

const commonStubs = {
  'q-dialog': {
    template: '<div class="q-dialog" v-if="modelValue"><slot /></div>',
    props: ['modelValue']
  },
  'q-card': { template: '<div class="q-card"><slot /></div>' },
  'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
  'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
  'q-banner': { template: '<div class="q-banner"><slot /></div>' },
  'q-form': {
    template: '<form @submit.prevent="$emit(\'submit\', $event)"><slot /></form>',
    methods: {
      validate: () => Promise.resolve(true)
    }
  },
  'q-input': {
    template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'rules', 'type', 'label']
  },
  'q-select': {
    template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><slot /><slot name="option" :opt="{}" :itemProps="{}" /></select>',
    props: ['modelValue', 'options', 'rules', 'label']
  },
  'q-checkbox': {
    template: '<input type="checkbox" :checked="modelValue" @change="$emit(\'update:modelValue\', $event.target.checked)" />',
    props: ['modelValue', 'label']
  },
  'q-btn': {
    template: '<button :type="type" @click="$emit(\'click\')"><slot />{{ label }}</button>',
    props: ['type', 'label']
  },
  'q-icon': true,
  'q-separator': true,
  'q-item': true,
  'q-item-section': true,
  'q-item-label': true
}

describe('LessonPlanner Dialog Components', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('LessonFormDialog.vue', () => {
    it('renders create mode and handles submit', async () => {
      const wrapper = mount(LessonFormDialog, {
        props: {
          modelValue: true,
          isEditing: false,
          availableSubjectOptions: [
            { subject_id: 1, subject_name: 'Matematica' },
            { subject_id: 2, subject_name: 'Fisica' }
          ],
          activityTypeOptions: [
            { label: 'Standard', value: 'standard' }
          ],
          initialData: {
            subject_id: 1,
            topic: 'Algebra',
            date: '2026-09-10',
            hour: 1,
            duration: 1,
            type: 'Frontale'
          }
        },
        global: {
          plugins: [[Quasar, {}]],
          mocks: {
            $t: (key) => key
          },
          stubs: commonStubs
        }
      })

      expect(wrapper.text()).toContain('classRegister.newLesson')
      const form = wrapper.find('form')
      expect(form.exists()).toBe(true)

      await form.trigger('submit')

      expect(wrapper.emitted('save')).toBeTruthy()
      expect(wrapper.emitted('save')[0][0]).toMatchObject({
        subject_id: 1,
        topic: 'Algebra',
        date: '2026-09-10',
        hour: 1,
        duration: 1,
        type: 'Frontale'
      })
    })

    it('renders edit mode with initial values', () => {
      const wrapper = mount(LessonFormDialog, {
        props: {
          modelValue: true,
          isEditing: true,
          initialData: {
            subject_id: 2,
            topic: 'Meccanica Quantistica',
            date: '2026-09-15',
            hour: 2,
            duration: 2,
            type: 'Laboratorio'
          }
        },
        global: {
          plugins: [[Quasar, {}]],
          mocks: {
            $t: (key) => key
          },
          stubs: commonStubs
        }
      })

      expect(wrapper.text()).toContain('classRegister.editLesson')
    })
  })

  describe('HomeworkFormDialog.vue', () => {
    it('renders and emits save with valid data', async () => {
      const wrapper = mount(HomeworkFormDialog, {
        props: {
          modelValue: true,
          availableSubjectOptions: [
            { subject_id: 1, subject_name: 'Matematica' }
          ],
          initialData: {
            subject_id: 1,
            description: 'Esercizi pag. 50',
            dueDate: '2026-09-20'
          }
        },
        global: {
          plugins: [[Quasar, {}]],
          mocks: {
            $t: (key) => key
          },
          stubs: commonStubs
        }
      })

      expect(wrapper.text()).toContain('homework.assignHomework')
      const form = wrapper.find('form')
      await form.trigger('submit')

      expect(wrapper.emitted('save')).toBeTruthy()
      expect(wrapper.emitted('save')[0][0]).toEqual({
        subject_id: 1,
        description: 'Esercizi pag. 50',
        dueDate: '2026-09-20'
      })
    })
  })

  describe('FreeActivityDialog.vue', () => {
    it('renders and emits save with free activity payload', async () => {
      const wrapper = mount(FreeActivityDialog, {
        props: {
          modelValue: true,
          isEditing: false,
          freeActivityTypeOptions: [
            { label: 'A disposizione', value: 'disposizione', icon: 'person', color: 'teal' }
          ],
          initialData: {
            date: '2026-09-12',
            start_hour: 2,
            duration: 1,
            activity_type: 'disposizione',
            description: 'Ora buco a disposizione',
            notes: 'In sala professori'
          }
        },
        global: {
          plugins: [[Quasar, {}]],
          mocks: {
            $t: (key) => key
          },
          stubs: commonStubs
        }
      })

      expect(wrapper.text()).toContain('Nuova Attività Libera')
      const form = wrapper.find('form')
      await form.trigger('submit')

      expect(wrapper.emitted('save')).toBeTruthy()
      expect(wrapper.emitted('save')[0][0]).toMatchObject({
        date: '2026-09-12',
        start_hour: 2,
        duration: 1,
        activity_type: 'disposizione',
        description: 'Ora buco a disposizione',
        notes: 'In sala professori'
      })
    })
  })
})
