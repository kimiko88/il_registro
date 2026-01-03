import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import DocumentEditor from '@/components/Teacher/DocumentEditor.vue'

describe('Teacher/DocumentEditor.vue', () => {
    it('updates model value', async () => {
        const wrapper = mount(DocumentEditor, {
            props: { modelValue: 'Initial' },
            global: {
                stubs: {
                    'q-editor': {
                        props: ['modelValue'],
                        template: '<div class="q-editor-stub" @input="$emit(\'update:modelValue\', \'New Content\')"></div>'
                    }
                }
            }
        })

        const editor = wrapper.find('.q-editor-stub')
        await editor.trigger('input')

        expect(wrapper.emitted('update:modelValue')[0]).toEqual(['New Content'])
    })
})
