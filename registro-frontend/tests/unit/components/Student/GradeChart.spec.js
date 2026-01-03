import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import GradeChart from '@/components/Student/GradeChart.vue'

describe('Student/GradeChart.vue', () => {
    const averages = {
        'Math': 8.5,
        'History': 5.5,
        'Science': 6.0
    }

    it('renders list of subjects and averages', () => {
        const wrapper = mount(GradeChart, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    // Use standard kebab name to match template usage
                    'q-linear-progress': {
                        template: '<div class="progress" :data-value="value" :data-color="color"><slot /></div>',
                        props: ['value', 'color', 'size', 'rounded']
                    },
                    'q-badge': { template: '<span class="badge">{{ label }}</span>', props: ['label'] }
                }
            },
            props: { averages }
        })

        // Check subject names
        expect(wrapper.text()).toContain('Math')
        expect(wrapper.text()).toContain('History')
        expect(wrapper.text()).toContain('Science')

        // Check progress bars matches
        const progressBars = wrapper.findAll('.progress')
        expect(progressBars).toHaveLength(3)

        // Math: 8.5/10 = 0.85
        const mathBar = progressBars[0]
        expect(Number(mathBar.attributes('data-value'))).toBe(0.85)
    })

    it('calculates color correctly', () => {
        const wrapper = mount(GradeChart, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-linear-progress': {
                        template: '<div class="progress" :data-color="color"></div>',
                        props: ['color', 'value', 'size', 'rounded']
                    },
                    'q-badge': true
                }
            },
            props: { averages }
        })

        // Math 8.5 -> Green
        // History 5.5 -> Red
        // Science 6.0 -> Green

        const bars = wrapper.findAll('.progress')
        const colors = bars.map(b => b.attributes('data-color'))

        expect(colors).toContain('green')
        expect(colors).toContain('red')
        expect(colors.filter(c => c === 'green').length).toBe(2)
    })
})
