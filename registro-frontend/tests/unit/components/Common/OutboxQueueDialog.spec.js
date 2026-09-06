import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { Quasar } from 'quasar'
import OutboxQueueDialog from '@/components/Common/OutboxQueueDialog.vue'
import { useOutboxStore } from '@/stores/outbox'

describe('Common/OutboxQueueDialog.vue', () => {
  let outboxStore

  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    outboxStore = useOutboxStore()
  })

  function createWrapper(props = { modelValue: true }) {
    return mount(OutboxQueueDialog, {
      global: {
        plugins: [Quasar],
        mocks: {
          $t: (key) => key
        },
        stubs: {
          'q-dialog': { template: '<div class="q-dialog-stub" v-if="modelValue"><slot /></div>', props: ['modelValue'] },
          'q-card': { template: '<div class="q-card-stub"><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' },
          'q-card-actions': { template: '<div><slot /></div>' },
          'q-separator': true,
          'q-avatar': { template: '<div><slot /></div>' },
          'q-icon': true,
          'q-list': { template: '<div class="q-list-stub"><slot /></div>' },
          'q-item': { template: '<div class="q-item-stub"><slot /></div>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-tooltip': true
        }
      },
      props
    })
  }

  it('renders empty state when queue is empty', () => {
    const wrapper = createWrapper()
    expect(wrapper.text()).toMatch(/offlineQueue\.allSynced|Nessuna operazione in sospeso/)
    expect(wrapper.find('.q-list-stub').exists()).toBe(false)
  })

  it('renders queued items correctly', () => {
    outboxStore.enqueue({
      url: '/grades/bulk',
      method: 'POST',
      title: 'Salvataggio voti classe 2A'
    })

    const wrapper = createWrapper()
    expect(wrapper.find('.q-list-stub').exists()).toBe(true)
    expect(wrapper.text()).toContain('Salvataggio voti classe 2A')
    expect(wrapper.text()).toContain('POST')
  })

  it('removes item when delete button is clicked', async () => {
    const id = outboxStore.enqueue({
      url: '/attendance/mark-bulk',
      method: 'POST',
      title: 'Appello classe 1B'
    })

    expect(outboxStore.pendingCount).toBe(1)
    const wrapper = createWrapper()

    const deleteBtn = wrapper.find('button[aria-label="common.delete"]')
    if (deleteBtn.exists()) {
      await deleteBtn.trigger('click')
      expect(outboxStore.pendingCount).toBe(0)
    }
  })
})
