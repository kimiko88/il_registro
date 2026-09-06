import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import GlobalSearch from '@/components/Common/GlobalSearch.vue'
import { createPinia, setActivePinia } from 'pinia'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

const mockRouter = {
  push: vi.fn()
}

vi.mock('vue-router', () => ({
  useRouter: () => mockRouter
}))

describe('GlobalSearch Component (Spotlight Ctrl+K)', () => {
  let pinia

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    vi.clearAllMocks()
  })

  const globalConfig = {
    plugins: [pinia],
    stubs: {
      teleport: true,
      'q-icon': true,
      'q-btn': true,
      'q-separator': true,
      'q-list': true,
      'q-item': true,
      'q-item-section': true,
      'q-item-label': true,
      'q-avatar': true,
      'q-chip': true,
      'q-spinner': true,
      'q-scroll-area': true
    }
  }

  it('renders closed by default and opens via open() expose', async () => {
    const wrapper = mount(GlobalSearch, {
      global: globalConfig
    })

    expect(wrapper.find('.global-search-overlay').exists()).toBe(false)
    expect(wrapper.vm.isOpen).toBe(false)

    wrapper.vm.open()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.isOpen).toBe(true)
    expect(wrapper.find('.global-search-overlay').exists()).toBe(true)

    wrapper.vm.close()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.isOpen).toBe(false)
  })

  it('provides quick actions and default items when query is empty', async () => {
    const wrapper = mount(GlobalSearch, {
      global: globalConfig
    })

    wrapper.vm.open()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.defaultItems.length).toBeGreaterThan(0)
    const darkAction = wrapper.vm.defaultItems.find(i => i.id === 'qa-dark')
    expect(darkAction).toBeDefined()
    expect(typeof darkAction.action).toBe('function')
  })

  it('calls tenant /search endpoint on query entry', async () => {
    api.get.mockResolvedValueOnce({
      data: {
        results: [
          { id: '1', type: 'student', title: 'Mario Rossi', description: 'Classe 3A' },
          { id: '2', type: 'communication', title: 'Circolare n. 42', description: 'Chiusura scuola' }
        ]
      }
    })

    const wrapper = mount(GlobalSearch, {
      global: globalConfig
    })

    wrapper.vm.open()
    wrapper.vm.query = 'Rossi'
    await wrapper.vm.$nextTick()

    // Wait for debounce (250ms)
    await new Promise(resolve => setTimeout(resolve, 300))

    expect(api.get).toHaveBeenCalledWith('/search', expect.objectContaining({
      params: expect.objectContaining({ q: 'Rossi' })
    }))
    expect(wrapper.vm.results.length).toBeGreaterThan(0)
  })

  it('handles keyboard navigation with moveDown and moveUp', async () => {
    const wrapper = mount(GlobalSearch, {
      global: globalConfig
    })

    wrapper.vm.open()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.highlightedIndex).toBe(0)
    wrapper.vm.moveDown()
    expect(wrapper.vm.highlightedIndex).toBe(1)
    wrapper.vm.moveUp()
    expect(wrapper.vm.highlightedIndex).toBe(0)
    wrapper.vm.moveUp()
    expect(wrapper.vm.highlightedIndex).toBe(0)
  })
})
