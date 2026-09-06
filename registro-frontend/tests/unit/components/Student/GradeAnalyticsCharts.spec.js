import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import GradeAnalyticsCharts from '@/components/Student/GradeAnalyticsCharts.vue'

describe('GradeAnalyticsCharts.vue Suite', () => {
  const sampleGrades = [
    { id: '1', grade_value: 5.5, subject: 'Matematica', date: '2026-01-10' },
    { id: '2', grade_value: 6.5, subject: 'Italiano', date: '2026-01-15' },
    { id: '3', grade_value: 8.0, subject: 'Inglese', date: '2026-01-20' },
    { id: '4', grade_value: 9.0, subject: 'Storia', date: '2026-01-25' }
  ]

  const stubs = {
    'q-card': { template: '<div class="q-card"><slot /></div>' },
    'q-card-section': { template: '<div class="q-card-section"><slot /></div>' },
    'q-icon': { template: '<i class="q-icon" />' },
    'q-badge': { template: '<span class="q-badge"><slot /></span>' }
  }

  const mockTranslate = (k) => {
    const map = {
      'roleDashboards.totalGradesEvaluated': 'valutazioni',
      'roleDashboards.insufficientData': 'Inserisci almeno due voti per visualizzare il trend',
      'roleDashboards.analyticsTitle': 'Analisi Rendimento & Trend',
      'roleDashboards.analyticsSubtitle': 'Evoluzione cronologica dei voti e regolarità presenze'
    }
    return map[k] || k
  }

  const globalConfig = {
    mocks: {
      $t: mockTranslate,
      t: mockTranslate
    },
    stubs
  }

  it('renders SVG sparkline curve and interactive points when >= 2 grades provided', async () => {
    const wrapper = mount(GradeAnalyticsCharts, {
      props: {
        grades: sampleGrades,
        attendanceRate: 92
      },
      global: globalConfig
    })

    // Sparkline SVG should be present
    const svg = wrapper.find('.sparkline-svg')
    expect(svg.exists()).toBe(true)

    // Check data dots rendered (4 dots)
    const dots = wrapper.findAll('.chart-dot')
    expect(dots.length).toBe(4)

    // Distribution breakdown
    // 5.5 (<6) -> 1 fail (25%)
    // 6.5 (6-7.5) -> 1 pass (25%)
    // 8.0, 9.0 (>=8) -> 2 good (50%)
    expect(wrapper.text()).toContain('4 valutazioni')
    expect(wrapper.text()).toContain('25%')
    expect(wrapper.text()).toContain('50%')

    // Attendance display
    expect(wrapper.text()).toContain('92%')
    expect(wrapper.text()).toContain('Regolare')

    // Test hover on point
    expect(wrapper.find('rect[fill-opacity="0.9"]').exists()).toBe(false)
    await dots[2].trigger('mouseenter')
    expect(wrapper.find('rect[fill-opacity="0.9"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('8 • Inglese')

    await dots[2].trigger('mouseleave')
    expect(wrapper.find('rect[fill-opacity="0.9"]').exists()).toBe(false)
  })

  it('displays fallback state when fewer than 2 grades are available', () => {
    const wrapper = mount(GradeAnalyticsCharts, {
      props: {
        grades: [{ id: '1', grade_value: 8.0, subject: 'Arte', date: '2026-02-01' }],
        attendanceRate: 70
      },
      global: globalConfig
    })

    // Fallback message should be visible
    expect(wrapper.find('.sparkline-svg').exists()).toBe(false)
    expect(wrapper.text()).toContain('Inserisci almeno due voti per visualizzare il trend')

    // Attendance at 70% is marked as at risk
    expect(wrapper.text()).toContain('70%')
    expect(wrapper.text()).toContain('A Rischio')
  })

  it('filters out non-numeric and out-of-bounds grades gracefully', () => {
    const wrapper = mount(GradeAnalyticsCharts, {
      props: {
        grades: [
          { id: '1', grade_value: 'A', subject: 'Matematica', date: '2026-01-01' },
          { id: '2', grade_value: -1, subject: 'Fisica', date: '2026-01-02' },
          { id: '3', grade_value: 7.0, subject: 'Chimica', date: '2026-01-03' },
          { id: '4', grade_value: 8.5, subject: 'Biologia', date: '2026-01-04' }
        ],
        attendanceRate: 85
      },
      global: globalConfig
    })

    // Only the 2 valid numeric grades (7.0 and 8.5) should be plotted
    const dots = wrapper.findAll('.chart-dot')
    expect(dots.length).toBe(2)
    expect(wrapper.text()).toContain('2 valutazioni')
  })
})
