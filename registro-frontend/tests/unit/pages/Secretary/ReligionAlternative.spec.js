import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ReligionAlternative from '@/pages/secretary/ReligionAlternative.vue'
import { religionAlternativeService } from '@/services/religionAlternativeService'

vi.mock('@/services/religionAlternativeService', () => ({
    religionAlternativeService: {
        getOptions: vi.fn(),
        getEvaluations: vi.fn(),
        saveOption: vi.fn(),
        saveEvaluation: vi.fn()
    }
}))

describe('ReligionAlternative.vue', () => {
    it('renders religion options, counts, and evaluations', async () => {
        religionAlternativeService.getOptions.mockResolvedValueOnce({
            data: [
                { id: 'opt-1', student_id: 's-1', option_type: 'materia_alternativa', notes: 'Gruppo aperto' },
                { id: 'opt-2', student_id: 's-2', option_type: 'uscita_scuola', notes: 'Esonero orario' }
            ]
        })
        religionAlternativeService.getEvaluations.mockResolvedValueOnce({
            data: [
                { id: 'eval-1', student_id: 's-1', subject_kind: 'materia_alternativa', period: 'q1', judgment_level: 'ottimo' }
            ]
        })

        const wrapper = mount(ReligionAlternative, {
            global: {
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': true,
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-separator': true,
                    'q-badge': true,
                    'q-chip': true,
                    'q-tabs': true,
                    'q-tab': true,
                    'q-table': true,
                    'q-dialog': true
                }
            }
        })

        await wrapper.vm.$nextTick()
        expect(religionAlternativeService.getOptions).toHaveBeenCalled()
        expect(wrapper.vm.options).toHaveLength(2)
        expect(wrapper.vm.countByOption('materia_alternativa')).toBe(1)
        expect(wrapper.vm.countByOption('uscita_scuola')).toBe(1)
    })
})
