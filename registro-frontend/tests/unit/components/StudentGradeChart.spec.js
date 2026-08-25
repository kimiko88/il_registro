import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import GradeChart from '@/components/Student/GradeChart.vue'

describe('GradeChart.vue Robustness Suite', () => {
  it('renders correctly with numeric averages without NaN', () => {
    const wrapper = mount(GradeChart, {
      props: {
        averages: {
          Matematica: 8.5,
          Italiano: 5.5,
          Storia: 6.0
        }
      },
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-linear-progress': {
            props: ['value', 'color'],
            template: '<div class="progress" :data-value="value" :data-color="color"><slot /></div>'
          },
          'q-badge': {
            props: ['label'],
            template: '<span class="badge">{{ label }}</span>'
          }
        }
      }
    })

    expect(wrapper.text()).toContain('Matematica')
    expect(wrapper.text()).toContain('Italiano')
    expect(wrapper.text()).toContain('Storia')

    const progresses = wrapper.findAll('.progress')
    expect(progresses).toHaveLength(3)

    // Check Math (8.5 -> value 0.85, green)
    expect(progresses[0].attributes('data-value')).toBe('0.85')
    expect(progresses[0].attributes('data-color')).toBe('green')

    // Check Italian (5.5 -> value 0.55, red)
    expect(progresses[1].attributes('data-value')).toBe('0.55')
    expect(progresses[1].attributes('data-color')).toBe('red')
  })

  it('safely handles non-numeric averages (e.g. "-") without throwing or producing NaN', () => {
    const wrapper = mount(GradeChart, {
      props: {
        averages: {
          Fisica: '-'
        }
      },
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-linear-progress': {
            props: ['value', 'color'],
            template: '<div class="progress" :data-value="value" :data-color="color"><slot /></div>'
          },
          'q-badge': {
            props: ['label'],
            template: '<span class="badge">{{ label }}</span>'
          }
        }
      }
    })

    expect(wrapper.text()).toContain('Fisica')
    const progress = wrapper.find('.progress')
    expect(progress.attributes('data-value')).toBe('0')
    expect(progress.attributes('data-color')).toBe('grey')
    expect(wrapper.text()).toContain('-')
  })

  it('handles empty averages object gracefully', () => {
    const wrapper = mount(GradeChart, {
      props: {
        averages: {}
      },
      global: {
        mocks: { t: (k) => k },
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' }
        }
      }
    })

    expect(wrapper.findAll('.row')).toHaveLength(0)
  })
})
