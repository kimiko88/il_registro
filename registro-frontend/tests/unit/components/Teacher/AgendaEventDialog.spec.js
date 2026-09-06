import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import { createPinia, setActivePinia } from 'pinia'
import AgendaEventDialog from '@/components/Teacher/AgendaEventDialog.vue'
import { useAgendaStore } from '@/stores/agenda'
import { useAuthStore } from '@/stores/auth'

const mockNotify = vi.fn()
const mockDialog = vi.fn()

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      notify: mockNotify,
      dialog: (opts) => {
        mockDialog(opts)
        return {
          onOk: (fn) => { fn(); return { onCancel: () => {} } }
        }
      }
    })
  }
})

describe('AgendaEventDialog.vue', () => {
  let pinia
  let agendaStore
  let authStore

  beforeEach(() => {
    vi.clearAllMocks()
    pinia = createPinia()
    setActivePinia(pinia)
    agendaStore = useAgendaStore()
    authStore = useAuthStore()

    authStore.user = { id: 'teacher-1', first_name: 'Mario', last_name: 'Rossi' }

    agendaStore.createEvent = vi.fn().mockResolvedValue({ id: 'new-1' })
    agendaStore.updateEvent = vi.fn().mockResolvedValue({ id: 'updated-1' })
    agendaStore.deleteEvent = vi.fn().mockResolvedValue({})
  })

  function createWrapper(props = {}) {
    return mount(AgendaEventDialog, {
      props: {
        modelValue: true,
        classOptions: [
          { label: '3A Informatica', value: 'class-1' },
          { label: '4B Liceo', value: 'class-2' }
        ],
        ...props
      },
      global: {
        plugins: [pinia, [Quasar, {}]],
        mocks: {
          $t: (key, params) => {
            if (key === 'agendaPage.moreCount') return `+${params?.count} altri`
            return key
          }
        },
        stubs: {
          'q-dialog': { template: '<div class="q-dialog"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-form': {
            template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>'
          },
          'q-input': {
            template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
            props: ['modelValue', 'rules']
          },
          'q-select': {
            template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><slot /></select>',
            props: ['modelValue', 'options', 'rules']
          },
          'q-toggle': {
            template: '<input type="checkbox" :checked="modelValue" @change="$emit(\'update:modelValue\', $event.target.checked)" />',
            props: ['modelValue']
          },
          'q-btn': {
            template: '<button :type="type" @click="$emit(\'click\')"><slot />{{ label }}</button>',
            props: ['type', 'label']
          },
          'q-icon': true
        }
      }
    })
  }

  it('renders create mode with default or initial values', () => {
    const wrapper = createWrapper({
      initialDate: '2026-09-15',
      initialHour: '10:00'
    })

    expect(wrapper.text()).toMatch(/Nuovo Evento|agendaPage\.newEvent/)
    expect(wrapper.find('form').exists()).toBe(true)
  })

  it('populates fields correctly when editing an existing event', () => {
    const event = {
      id: 'event-99',
      teacher_id: 'teacher-1',
      title: 'Verifica di Fisica',
      description: 'Elettromagnetismo',
      type: 'verifica',
      all_day: false,
      class_id: 'class-1',
      date: '2026-09-20T08:00:00Z',
      start_time: '11:00',
      end_time: '12:00',
      visible_to_students: true
    }

    const wrapper = createWrapper({ event })

    expect(wrapper.text()).toMatch(/Salva Modifiche|agendaPage\.saveChanges/)
  })

  it('creates new event and emits saved on form submit', async () => {
    const wrapper = createWrapper({
      initialDate: '2026-09-15',
      initialHour: '10:00'
    })

    // Trigger form submit
    await wrapper.find('form').trigger('submit')

    expect(agendaStore.createEvent).toHaveBeenCalled()
    expect(wrapper.emitted('saved')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
  })

  it('updates existing event on form submit', async () => {
    const event = {
      id: 'event-99',
      teacher_id: 'teacher-1',
      title: 'Verifica di Fisica',
      description: 'Elettromagnetismo',
      type: 'verifica',
      all_day: false,
      class_id: 'class-1',
      date: '2026-09-20',
      start_time: '11:00',
      end_time: '12:00',
      visible_to_students: true
    }

    const wrapper = createWrapper({ event })
    await wrapper.find('form').trigger('submit')

    expect(agendaStore.updateEvent).toHaveBeenCalledWith('event-99', expect.objectContaining({
      title: 'Verifica di Fisica',
      type: 'verifica',
      class_id: 'class-1'
    }))
    expect(wrapper.emitted('saved')).toBeTruthy()
  })

  it('deletes event when confirmed', async () => {
    const event = {
      id: 'event-99',
      teacher_id: 'teacher-1',
      title: 'Verifica di Fisica',
      class_id: 'class-1'
    }

    const wrapper = createWrapper({ event })
    // Delete button exists in edit mode for current author
    const deleteBtn = wrapper.findAll('button').find(b => b.text().includes('Elimina Evento') || b.text().includes('agendaPage.deleteEvent'))
    expect(deleteBtn).toBeDefined()
    await deleteBtn.trigger('click')

    expect(agendaStore.deleteEvent).toHaveBeenCalledWith('event-99')
    expect(wrapper.emitted('deleted')).toBeTruthy()
    expect(wrapper.emitted('deleted')[0]).toEqual(['event-99'])
  })

  it('shows read-only notice when user is not author of event', () => {
    const event = {
      id: 'event-88',
      teacher_id: 'teacher-other',
      teacher_name: 'Prof. Bianchi',
      title: 'Consiglio di Classe'
    }

    const wrapper = createWrapper({ event })
    expect(wrapper.text()).toContain('Prof. Bianchi')
    // Delete button should not be rendered for non-authors
    const deleteBtn = wrapper.findAll('button').find(b => b.text().includes('agendaPage.deleteEvent'))
    expect(deleteBtn).toBeUndefined()
  })
})
