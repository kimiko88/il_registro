import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SidiExports from '@/pages/secretary/SidiExports.vue'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({
      data: [
        {
          id: 'exp-1',
          export_type: 'ANS_ANAGRAFE',
          school_year: '2025/2026',
          file_name: 'FLUSSO_SIDI_ANS_ANAGRAFE.xml',
          records_count: 240
        }
      ]
    }),
    post: vi.fn().mockResolvedValue({
      data: {
        export: { id: 'new-exp', file_name: 'FLUSSO_SIDI_SCRUTINIO.xml' },
        validation: { valid: true, total_records: 240, errors: [] }
      }
    })
  }
}))

describe('SidiExports.vue Page', () => {
  it('renders SIDI export title correctly', () => {
    const wrapper = mount(SidiExports, {
      global: {
        stubs: {
          'q-page': { template: '<div><slot /></div>' },
          'q-icon': true,
          'q-btn': true,
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-select': true,
          'q-banner': true,
          'q-chip': true,
          'q-list': { template: '<div><slot /></div>' },
          'q-item': { template: '<div><slot /></div>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-avatar': true
        }
      }
    })

    expect(wrapper.text()).toContain('Flussi SIDI & Tracciati Ministeriali (MIM)')
  })
})
