import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { Quasar } from 'quasar'
import GradesList from '@/components/Student/GradesList.vue'

describe('Student/GradesList.vue', () => {
    const gradesBySubject = {
        'Math': [
            { id: 1, value: 8, description: 'Test', type: 'Oral', date: '2023-01-01' },
            { id: 2, value: 5, description: 'Quiz', type: 'Written', date: '2023-01-02' }
        ],
        'History': [
            { id: 3, value: 7, description: 'Essay', type: 'Written', date: '2023-01-03' }
        ]
    }
    const averages = { 'Math': 6.5, 'History': 7 }
    const mockGetTrend = (subject) => {
        if (subject === 'Math') return 'up'
        if (subject === 'History') return 'down'
        return 'flat'
    }

    it('renders subjects and grades', () => {
        const wrapper = mount(GradesList, {
            global: {
                plugins: [Quasar],
                stubs: {
                    'q-expansion-item': { template: '<div><slot name="header"></slot><slot></slot></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-badge': { template: '<span class="q-badge"><slot /></span>' },
                    'q-item-label': { template: '<span><slot /></span>' },
                    'q-icon': { template: '<i class="q-icon"></i>' }
                }
            },
            props: {
                gradesBySubject,
                averages,
                getTrend: mockGetTrend
            }
        })

        // Check subjects
        expect(wrapper.text()).toContain('Math')
        expect(wrapper.text()).toContain('History')

        // Check grades
        expect(wrapper.text()).toContain('Test')
        expect(wrapper.text()).toContain('Quiz')
        expect(wrapper.text()).toContain('Essay')
    })
})
