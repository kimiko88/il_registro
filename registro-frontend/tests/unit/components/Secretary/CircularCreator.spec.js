import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import CircularCreator from '@/components/Secretary/CircularCreator.vue'
import { useCommunicationsStore } from '@/stores/communications'

describe('CircularCreator', () => {
    let wrapper
    let commStore

    beforeEach(() => {
        wrapper = mount(CircularCreator, {
            global: {
                plugins: [
                    Quasar,
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            classes: {
                                classes: [{ id: 'c1', name: '1', section: 'A' }]
                            }
                        }
                    })
                ],
                stubs: {
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
                    'q-input': true,
                    'q-checkbox': true,
                    'q-select': { template: '<div class="select-stub"></div>' },
                    'q-editor': true,
                    'q-file': true,
                    'q-icon': true,
                    'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' }
                }
            }
        })
        commStore = useCommunicationsStore()
    })

    it('shows warning if no recipients selected', async () => {
        wrapper.vm.form.recipients.teachers = false
        await wrapper.vm.sendCircular()

        expect(wrapper.vm.sending).toBe(false)
        expect(commStore.sendMessage).not.toHaveBeenCalled()
    })

    it('shows class selector when students or parents selected', async () => {
        expect(wrapper.find('.select-stub').exists()).toBe(false)

        // Select students
        wrapper.vm.form.recipients.students = true
        await wrapper.vm.$nextTick()

        expect(wrapper.find('.select-stub').exists()).toBe(true)
    })

    it('sends circular with valid data', async () => {
        wrapper.vm.form.title = 'Test'
        wrapper.vm.form.recipients.teachers = true

        // Mock sendMessage to resolve
        commStore.sendMessage.mockResolvedValue({})

        await wrapper.vm.sendCircular()

        expect(commStore.sendMessage).toHaveBeenCalledWith(expect.objectContaining({
            title: 'Test',
            type: 'circular'
        }))
        expect(wrapper.vm.sending).toBe(false)
        expect(wrapper.emitted('sent')).toBeTruthy()
    })
})
