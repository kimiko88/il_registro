import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import GradeBadge from '@/components/Common/GradeBadge.vue'

const customStubs = {
  'q-badge': {
    template: '<span class="q-badge" :class="[color ? `bg-${color}` : \'\']" :role="$attrs.role" :aria-label="$attrs[\'aria-label\']"><slot /></span>',
    props: ['color', 'dense']
  },
  'q-icon': {
    template: '<i :class="name"></i>',
    props: ['name', 'size']
  },
  'q-tooltip': {
    template: '<div class="q-tooltip"><slot /></div>'
  }
}

describe('GradeBadge.vue', () => {
  it('renders a sufficient grade correctly with green/positive style and check icon', () => {
    const wrapper = mount(GradeBadge, {
      props: {
        value: 8
      },
      global: {
        stubs: customStubs
      }
    })

    expect(wrapper.text()).toContain('8')
    expect(wrapper.attributes('role')).toBe('status')
    expect(wrapper.attributes('aria-label')).toContain('8')
    expect(wrapper.attributes('aria-label')).toContain('Sufficiente')
    expect(wrapper.find('i').classes()).toContain('check_circle')
  })

  it('renders an insufficient grade correctly with negative style and priority icon', () => {
    const wrapper = mount(GradeBadge, {
      props: {
        value: 4.5
      },
      global: {
        stubs: customStubs
      }
    })

    expect(wrapper.text()).toContain('4.5')
    expect(wrapper.attributes('aria-label')).toContain('4.5')
    expect(wrapper.attributes('aria-label')).toContain('Insufficiente')
    expect(wrapper.find('i').classes()).toContain('priority_high')
  })

  it('renders absence (-1 or "A") with special event_busy icon and absent text', () => {
    const wrapper = mount(GradeBadge, {
      props: {
        value: -1
      },
      global: {
        stubs: customStubs
      }
    })

    expect(wrapper.text()).toContain('A')
    expect(wrapper.attributes('aria-label')).toContain('Assente')
    expect(wrapper.find('i').classes()).toContain('event_busy')
  })

  it('extracts values properly from a grade object and renders evaluation type', () => {
    const gradeObj = {
      grade_value: 7.5,
      evaluation_type: 'Written',
      subject_name: 'Matematica',
      weight: 1.5,
      date: '2026-03-10'
    }

    const wrapper = mount(GradeBadge, {
      props: {
        value: gradeObj,
        showType: true
      },
      global: {
        stubs: customStubs
      }
    })

    expect(wrapper.text()).toContain('7.5')
    expect(wrapper.text()).toContain('(S)')
    expect(wrapper.attributes('aria-label')).toContain('7.5')
    expect(wrapper.attributes('aria-label')).toContain('Sufficiente')
  })
})
