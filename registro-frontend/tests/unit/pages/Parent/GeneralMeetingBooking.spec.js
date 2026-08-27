import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import GeneralMeetingBooking from '@/pages/parent/GeneralMeetingBooking.vue'
import colloquiService from '@/services/colloquiService'
import api from '@/services/api'

vi.mock('@/services/colloquiService', () => ({
  default: {
    listGeneralMeetings: vi.fn(),
    getGeneralMeeting: vi.fn(),
    listQueueTickets: vi.fn(),
    bookQueueTicket: vi.fn()
  }
}))

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

describe('GeneralMeetingBooking.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    colloquiService.listGeneralMeetings.mockResolvedValue({
      data: [
        {
          id: 'm-1',
          title: 'Ricevimento Generale Pomeridiano',
          event_date: '2025-11-20',
          start_time: '15:00',
          end_time: '19:00',
          slot_duration_minutes: 7,
          location_type: 'in_presenza',
          teacher_slots: [
            { id: 'ts-1', teacher_id: 't-1', teacher_name: 'Prof. Mario Rossi', subject_name: 'Matematica', room_or_table: 'Aula Magna', booked_count: 5, max_bookings: 30 }
          ]
        }
      ]
    })
    colloquiService.getGeneralMeeting.mockResolvedValue({
      data: {
        id: 'm-1',
        title: 'Ricevimento Generale Pomeridiano',
        event_date: '2025-11-20',
        start_time: '15:00',
        end_time: '19:00',
        slot_duration_minutes: 7,
        location_type: 'in_presenza',
        teacher_slots: [
          { id: 'ts-1', teacher_id: 't-1', teacher_name: 'Prof. Mario Rossi', subject_name: 'Matematica', room_or_table: 'Aula Magna', booked_count: 5, max_bookings: 30 }
        ]
      }
    })
    colloquiService.listQueueTickets.mockResolvedValue({
      data: [
        { id: 'tick-1', ticket_number: 6, scheduled_time: '15:35', teacher_name: 'Prof. Mario Rossi', student_name: 'Luca Bianchi', status: 'prenotato' }
      ]
    })
    api.get.mockImplementation((url) => {
      if (url.includes('/users/me/children')) return Promise.resolve({ data: [{ id: 'child-1', first_name: 'Luca', last_name: 'Bianchi' }] })
      return Promise.resolve({ data: [] })
    })
  })

  it('renders general meeting booking page with active tickets and available meetings', async () => {
    const wrapper = mount(GeneralMeetingBooking, {
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
          'q-btn': { template: '<button class="q-btn"><slot /></button>' },
          'q-badge': { template: '<span class="q-badge"><slot /></span>' },
          'q-separator': { template: '<hr />' },
          'q-icon': true,
          'q-tooltip': true,
          'q-dialog': true,
          'q-form': true,
          'q-input': true,
          'q-select': true,
          'q-space': true
        }
      }
    })

    await flushPromises()
    expect(wrapper.exists()).toBe(true)
    expect(colloquiService.listGeneralMeetings).toHaveBeenCalled()
  })
})
