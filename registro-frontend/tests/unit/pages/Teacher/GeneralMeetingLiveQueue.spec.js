import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import GeneralMeetingLiveQueue from '@/pages/teacher/GeneralMeetingLiveQueue.vue'
import colloquiService from '@/services/colloquiService'

vi.mock('@/services/colloquiService', () => ({
  default: {
    listGeneralMeetings: vi.fn(),
    listQueueTickets: vi.fn(),
    updateTicketStatus: vi.fn()
  }
}))

describe('GeneralMeetingLiveQueue.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    colloquiService.listGeneralMeetings.mockResolvedValue({
      data: [
        { id: 'm-1', title: 'Ricevimento Generale 1° Quadrimestre', event_date: '2025-11-20' }
      ]
    })
    colloquiService.listQueueTickets.mockResolvedValue({
      data: [
        { id: 't-1', ticket_number: 1, scheduled_time: '15:00', parent_name: 'Giuseppe Verdi', student_name: 'Anna Verdi', status: 'in_colloquio' },
        { id: 't-2', ticket_number: 2, scheduled_time: '15:07', parent_name: 'Elena Bianchi', student_name: 'Marco Bianchi', status: 'prenotato' }
      ]
    })
  })

  it('renders live queue and displays current and next tickets', async () => {
    const wrapper = mount(GeneralMeetingLiveQueue, {
      global: {
        plugins: [createTestingPinia()],
        mocks: {
          $t: (msg) => msg
        },
        stubs: {
          'q-page': { template: '<div class="q-page"><slot /></div>' },
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
          'q-card-actions': { template: '<div class="q-card-actions"><slot /></div>' },
          'q-table': { template: '<div class="q-table"><slot /></div>' },
          'q-btn': { template: '<button class="q-btn"><slot /></button>' },
          'q-badge': { template: '<span class="q-badge"><slot /></span>' },
          'q-icon': true,
          'q-tooltip': true,
          'q-skeleton': true,
          'q-select': { template: '<div class="q-select-stub"><slot /></div>' }
        }
      }
    })

    await flushPromises()
    expect(wrapper.exists()).toBe(true)
    expect(colloquiService.listGeneralMeetings).toHaveBeenCalled()
    expect(colloquiService.listQueueTickets).toHaveBeenCalled()
  })
})
