import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import AgendaEventDialog from '@/components/Teacher/AgendaEventDialog.vue'

describe('AgendaEventDialog.vue - Anti-Overlap Test Warnings', () => {
  let wrapper

  beforeEach(() => {
    vi.clearAllMocks()

    wrapper = mount(AgendaEventDialog, {
      props: {
        modelValue: true,
        initialDate: '2026-09-15',
        classOptions: [{ label: '4A', value: 'class-1' }]
      },
      global: {
        plugins: [
          [Quasar, {}],
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              auth: { user: { id: 'teacher-1' } },
              agenda: {
                events: [
                  {
                    id: 'ev-1',
                    class_id: 'class-1',
                    type: 'verifica',
                    title: 'Verifica Matematica',
                    date: '2026-09-15',
                    teacher_name: 'Prof. Gauss'
                  },
                  {
                    id: 'ev-2',
                    class_id: 'class-1',
                    type: 'verifica',
                    title: 'Verifica Fisica',
                    date: '2026-09-16',
                    teacher_name: 'Prof. Newton'
                  }
                ]
              }
            },
            stubActions: true
          })
        ],
        stubs: {
          'q-dialog': { template: '<div><slot /></div>' },
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-form': { template: '<form><slot /></form>' },
          'q-input': true,
          'q-select': true,
          'q-btn': true,
          'q-toggle': true,
          'q-icon': true,
          'q-banner': { template: '<div class="q-banner-stub"><slot name="avatar" /><slot /></div>' },
          'transition': { template: '<div><slot /></div>' }
        }
      }
    })
  })

  it('does not show conflict warning for homework (compito)', async () => {
    expect(wrapper.text()).not.toContain('Possibile Sovrapposizione Verifiche')
  })

  it('detects and shows warning banner when type is verifica and dates overlap', async () => {
    // Switch type to 'verifica' and class to 'class-1' on 2026-09-15
    wrapper.vm.form.type = 'verifica'
    wrapper.vm.form.class_id = 'class-1'
    wrapper.vm.form.date = '2026-09-15'

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.testConflicts.hasConflict).toBe(true)
    expect(wrapper.vm.testConflicts.sameDay.length).toBe(1)
    expect(wrapper.vm.testConflicts.sameDay[0].title).toBe('Verifica Matematica')
    expect(wrapper.text()).toContain('Possibile Sovrapposizione Verifiche')
    expect(wrapper.text()).toContain('Verifica Matematica')
  })

  it('detects weekly overlap when 2 or more tests occur in the same week', async () => {
    // Change date to 2026-09-17 (Thursday, same week as 15th and 16th)
    wrapper.vm.form.type = 'verifica'
    wrapper.vm.form.class_id = 'class-1'
    wrapper.vm.form.date = '2026-09-17'

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.testConflicts.sameDay.length).toBe(0)
    expect(wrapper.vm.testConflicts.sameWeek.length).toBe(2)
    expect(wrapper.vm.testConflicts.hasConflict).toBe(true)
    expect(wrapper.text()).toContain('Possibile Sovrapposizione Verifiche')
  })
})
