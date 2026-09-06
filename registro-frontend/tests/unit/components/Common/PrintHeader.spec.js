import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PrintHeader from '@/components/Common/PrintHeader.vue'

// Stub per useI18n — restituisce una funzione t() che ritorna la chiave stessa
const mockI18n = {
  install(app) {
    app.config.globalProperties.$t = (key) => key
  }
}

const mountComponent = (props = {}) =>
  mount(PrintHeader, {
    props,
    global: {
      plugins: [mockI18n],
      mocks: {
        $t: (key) => key,
      },
      // Stub useI18n composable
      provide: {},
    },
  })

describe('PrintHeader.vue', () => {
  it('renders the school name', () => {
    const wrapper = mountComponent({ schoolName: 'Istituto Comprensivo Test' })
    expect(wrapper.find('.print-header__title').text()).toBe('Istituto Comprensivo Test')
  })

  it('renders school subtitle when provided', () => {
    const wrapper = mountComponent({
      schoolName: 'IC Test',
      schoolSubtitle: 'Via Roma 1, Milano'
    })
    expect(wrapper.find('.print-header__subtitle').text()).toBe('Via Roma 1, Milano')
  })

  it('does not render subtitle element when not provided', () => {
    const wrapper = mountComponent({ schoolName: 'IC Test', schoolSubtitle: '' })
    expect(wrapper.find('.print-header__subtitle').exists()).toBe(false)
  })

  it('renders documentType and academicYear in meta', () => {
    const wrapper = mountComponent({
      schoolName: 'IC Test',
      documentType: 'Pagella 1° Quadrimestre',
      academicYear: '2026/2027',
      printDate: false
    })
    const meta = wrapper.find('.print-header__meta').text()
    expect(meta).toContain('Pagella 1° Quadrimestre')
    expect(meta).toContain('2026/2027')
  })

  it('shows print date by default', () => {
    const wrapper = mountComponent({
      schoolName: 'IC Test',
      printDate: true
    })
    // The date is today in locale format — verify it is non-empty
    const meta = wrapper.find('.print-header__meta').text()
    expect(meta.length).toBeGreaterThan(0)
  })

  it('hides print date when printDate is false', () => {
    const wrapper = mountComponent({
      schoolName: 'IC Test',
      documentType: '',
      academicYear: '',
      printDate: false
    })
    // Meta should be empty (no date, no type, no year)
    const meta = wrapper.find('.print-header__meta').text().trim()
    expect(meta).toBe('')
  })

  it('has role="banner" for accessibility', () => {
    const wrapper = mountComponent({ schoolName: 'IC Test' })
    expect(wrapper.find('.print-header').attributes('role')).toBe('banner')
  })

  it('has aria-label for screen readers', () => {
    const wrapper = mountComponent({ schoolName: 'IC Test' })
    expect(wrapper.find('.print-header').attributes('aria-label')).toBe(
      'Intestazione istituzionale'
    )
  })
})
