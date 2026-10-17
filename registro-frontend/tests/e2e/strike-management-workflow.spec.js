import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import StrikeManagement from '@/pages/ata/StrikeManagement.vue'
import strikeService from '@/services/strikeService'

vi.mock('@/services/strikeService', () => ({
  default: {
    getStrikeNotices: vi.fn(),
    getNoticeSummary: vi.fn(),
    createStrikeNotice: vi.fn(),
    deleteStrikeNotice: vi.fn()
  }
}))

describe('Strike Management Workflow E2E', () => {
  let pinia

  const commonStubs = {
    'q-page': { template: '<div class="q-page"><slot /></div>' },
    'q-card': { template: '<div class="q-card"><slot /></div>' },
    'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
    'q-table': {
      props: ['rows', 'columns'],
      template: '<div class="q-table" :data-count="rows?.length"><slot name="body" v-for="row in rows" :row="row" /><slot /></div>'
    },
    'q-dialog': {
      props: ['modelValue'],
      template: '<div v-if="modelValue" class="q-dialog"><slot /></div>'
    },
    'q-tabs': { template: '<div class="q-tabs"><slot /></div>' },
    'q-tab': { template: '<div class="q-tab"><slot /></div>' },
    'q-tab-panels': { template: '<div class="q-tab-panels"><slot /></div>' },
    'q-tab-panel': { template: '<div class="q-tab-panel"><slot /></div>' },
    'q-btn': {
      props: ['label'],
      template: '<button class="q-btn">{{ label }}<slot /></button>'
    },
    'q-icon': true,
    'q-badge': true,
    'q-chip': true,
    'q-input': {
      props: ['modelValue'],
      template: '<input class="q-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
    },
    'q-select': {
      template: '<div class="q-select"><slot /></div>'
    },
    'q-toggle': true,
    'q-circular-progress': true,
    'q-linear-progress': true,
    'q-tooltip': true,
    'q-avatar': true
  }

  const mockNoticesList = [
    {
      id: 'notice-active-1',
      title: 'Sciopero Generale Scuola',
      proclaimed_by: 'Organizzazioni Sindacali',
      strike_date: '2026-10-20',
      declaration_deadline: '2026-10-18T12:00:00Z',
      content: 'Sciopero intera giornata',
      notes: 'Garantire servizi minimi',
      created_by_name: 'DSGA Rossi',
      is_published: true,
      is_expired: false
    }
  ]

  const mockSummaryData = {
    notice: mockNoticesList[0],
    total_staff: 40,
    total_answered: 30,
    participates_count: 10,
    not_participates_count: 15,
    undecided_count: 5,
    unanswered_count: 10,
    participates_percent: 25.0,
    not_participates_percent: 37.5,
    undecided_percent: 12.5,
    unanswered_percent: 25.0,
    by_role: [
      { role: 'teacher', role_display: 'Docenti', total: 30, participates: 8, percent: 26.6 },
      { role: 'ata', role_display: 'Personale ATA', total: 10, participates: 2, percent: 20.0 }
    ],
    staff: [
      { user_id: 'u-1', first_name: 'Giulia', last_name: 'Verdi', role: 'teacher', role_display: 'Docente', intention: 'participates' }
    ]
  }

  beforeEach(() => {
    vi.clearAllMocks()
    strikeService.getStrikeNotices.mockResolvedValue(mockNoticesList)
    strikeService.getNoticeSummary.mockResolvedValue(mockSummaryData)

    pinia = createTestingPinia({
      initialState: {
        auth: {
          user: { id: 'dsga-1', role: 'dsga', name: 'Mario Rossi' },
          token: null
        }
      },
      stubActions: false
    })
  })

  it('1. loads and renders active strike notices for DSGA', async () => {
    const wrapper = mount(StrikeManagement, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    expect(strikeService.getStrikeNotices).toHaveBeenCalled()
    expect(wrapper.vm.notices.length).toBe(1)
    expect(wrapper.vm.notices[0].title).toBe('Sciopero Generale Scuola')
  })

  it('2. opens create notice dialog and validates required fields', async () => {
    const wrapper = mount(StrikeManagement, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    expect(wrapper.vm.showCreateDialog).toBe(false)

    wrapper.vm.openCreateDialog()
    expect(wrapper.vm.showCreateDialog).toBe(true)

    // Fill form
    wrapper.vm.createForm.title = 'Sciopero Straordinario'
    wrapper.vm.createForm.proclaimed_by = 'Sindacato Autonomo'
    wrapper.vm.createForm.strike_date = '2026-11-10'
    wrapper.vm.createForm.declaration_deadline = '2026-11-08T12:00'
    wrapper.vm.createForm.notes = 'Plesso Centrale e Succursale'
    wrapper.vm.createForm.publish_to_bacheca = true

    strikeService.createStrikeNotice.mockResolvedValue({ id: 'notice-new' })

    await wrapper.vm.submitCreateNotice()
    expect(strikeService.createStrikeNotice).toHaveBeenCalled()
    expect(wrapper.vm.showCreateDialog).toBe(false)
  })

  it('3. selects notice and fetches statistical summary', async () => {
    const wrapper = mount(StrikeManagement, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    await wrapper.vm.selectNotice('notice-active-1')

    expect(strikeService.getNoticeSummary).toHaveBeenCalledWith('notice-active-1')
    expect(wrapper.vm.selectedNoticeId).toBe('notice-active-1')
    expect(wrapper.vm.summary).toEqual(mockSummaryData)
    expect(wrapper.vm.summary.participates_percent).toBe(25.0)
  })

  it('4. performs notice deletion and reloads list', async () => {
    strikeService.deleteStrikeNotice.mockResolvedValue({ message: 'ok' })

    const wrapper = mount(StrikeManagement, {
      global: {
        plugins: [Quasar, pinia],
        stubs: commonStubs,
        mocks: { t: (k) => k }
      }
    })

    await flushPromises()
    wrapper.vm.confirmDeleteNotice(wrapper.vm.notices[0])
    expect(wrapper.vm.showDeleteDialog).toBe(true)
    expect(wrapper.vm.noticeToDelete).toEqual(wrapper.vm.notices[0])

    await wrapper.vm.executeDeleteNotice()
    expect(strikeService.deleteStrikeNotice).toHaveBeenCalledWith('notice-active-1')
    expect(wrapper.vm.showDeleteDialog).toBe(false)
  })
})
