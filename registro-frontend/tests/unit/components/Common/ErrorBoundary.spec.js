import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { defineComponent, h, ref } from 'vue'
import ErrorBoundary from '@/components/Common/ErrorBoundary.vue'
import { useErrorStore } from '@/stores/error'

const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  })
}))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useI18n: () => ({
      t: (key) => key
    })
  }
})

describe('ErrorBoundary.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPush.mockClear()
  })

  it('renders default slot content when there is no error', () => {
    const wrapper = mount(ErrorBoundary, {
      slots: {
        default: '<div class="normal-content">Contenuto normale</div>'
      }
    })

    expect(wrapper.find('.normal-content').exists()).toBe(true)
    expect(wrapper.find('.error-boundary-wrapper').exists()).toBe(false)
  })

  it('captures an error in child component and displays fallback UI', async () => {
    const shouldCrash = ref(false)

    const FaultyChild = defineComponent({
      name: 'FaultyChild',
      setup() {
        return () => {
          if (shouldCrash.value) {
            throw new Error('Crash simulato nel componente figlio!')
          }
          return h('div', { class: 'safe-child' }, 'Safe Child')
        }
      }
    })

    const wrapper = mount(ErrorBoundary, {
      slots: {
        default: () => h(FaultyChild)
      }
    })

    expect(wrapper.find('.safe-child').exists()).toBe(true)

    // Trigger error on re-render
    shouldCrash.value = true
    await wrapper.vm.$nextTick()
    await flushPromises()

    expect(wrapper.find('.error-boundary-wrapper').exists()).toBe(true)
    expect(wrapper.vm.errorMessage).toContain('Crash simulato nel componente figlio!')

    // Verify error was reported to errorStore
    const errorStore = useErrorStore()
    expect(errorStore.errorHistory.length).toBeGreaterThan(0)
  })

  it('resets error state when resetError is invoked', async () => {
    const shouldCrash = ref(false)

    const FaultyChild = defineComponent({
      name: 'FaultyChild',
      setup() {
        return () => {
          if (shouldCrash.value) {
            throw new Error('Test Error')
          }
          return h('div', { class: 'recovered-content' }, 'Recovered')
        }
      }
    })

    const wrapper = mount(ErrorBoundary, {
      slots: {
        default: () => h(FaultyChild)
      }
    })

    shouldCrash.value = true
    await wrapper.vm.$nextTick()
    await flushPromises()

    expect(wrapper.find('.error-boundary-wrapper').exists()).toBe(true)

    // Stop crashing, then reset error on boundary
    shouldCrash.value = false
    wrapper.vm.resetError()
    await wrapper.vm.$nextTick()
    await flushPromises()

    expect(wrapper.vm.hasError).toBe(false)
    expect(wrapper.find('.recovered-content').exists()).toBe(true)
  })

  it('navigates to dashboard on goHome click', async () => {
    const shouldCrash = ref(false)

    const FaultyChild = defineComponent({
      setup() {
        return () => {
          if (shouldCrash.value) {
            throw new Error('Test Home Error')
          }
          return h('div', 'ok')
        }
      }
    })

    const wrapper = mount(ErrorBoundary, {
      slots: {
        default: () => h(FaultyChild)
      }
    })

    shouldCrash.value = true
    await wrapper.vm.$nextTick()
    await flushPromises()

    expect(wrapper.find('.error-boundary-wrapper').exists()).toBe(true)
    wrapper.vm.goHome()
    expect(mockPush).toHaveBeenCalledWith('/')
  })
})
