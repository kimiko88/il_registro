import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import MainLayout from '@/layouts/MainLayout.vue'
import { useAuthStore } from '@/stores/auth'

const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRoute: () => ({
    path: '/teacher/classes',
    meta: { title: 'Le Mie Classi' }
  }),
  useRouter: () => ({
    push: mockPush
  })
}))

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useQuasar: () => ({
      dark: { isActive: false, toggle: vi.fn() },
      fullscreen: { isActive: false, toggle: vi.fn() },
      screen: { lt: { md: false } },
      notify: vi.fn()
    })
  }
})

describe('MainLayout.vue — Dynamic Role Navigation & Layout Features', () => {
  let wrapper

  beforeEach(() => {
    setActivePinia(createPinia())
    mockPush.mockClear()
  })

  afterEach(() => {
    if (wrapper) wrapper.unmount()
  })

  const mountLayout = () => {
    return mount(MainLayout, {
      global: {
        stubs: {
          'q-layout': { template: '<div class="q-layout"><slot /></div>' },
          'q-header': { template: '<header><slot /></header>' },
          'q-toolbar': { template: '<nav><slot /></nav>' },
          'q-toolbar-title': { template: '<div><slot /></div>' },
          'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
          'q-btn-dropdown': { template: '<div class="btn-dropdown"><slot /></div>' },
          'q-avatar': { template: '<div class="avatar"><slot /></div>' },
          'q-icon': { template: '<i class="icon" />' },
          'q-space': { template: '<div class="space" />' },
          'q-badge': { template: '<span class="badge"><slot /></span>' },
          'q-chip': { template: '<span class="chip"><slot /></span>' },
          'q-tooltip': { template: '<span class="tooltip"><slot /></span>' },
          'q-drawer': { template: '<aside><slot /></aside>' },
          'q-scroll-area': { template: '<div><slot /></div>' },
          'q-list': { template: '<ul><slot /></ul>' },
          'q-item': { template: '<li @click="$emit(\'click\')"><slot /></li>' },
          'q-item-section': { template: '<div><slot /></div>' },
          'q-item-label': { template: '<div><slot /></div>' },
          'q-expansion-item': { template: '<div><slot /></div>' },
          'q-select': true,
          'q-toggle': true,
          'q-spinner': true,
          'q-breadcrumbs': { template: '<nav class="breadcrumbs"><slot /></nav>' },
          'q-breadcrumbs-el': { template: '<span><slot /></span>' },
          'q-page-container': { template: '<main id="main-content"><slot /></main>' },
          'router-view': true,
          'GlobalSearch': true,
          'OnboardingTour': true,
          'HelpDrawer': true,
          'HelpCenterPanel': true,
          'SessionReauthDialog': true
        }
      }
    })
  }

  it('renders MainLayout with teacher role and builds menus', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 1, name: 'Mario Rossi', role: 'teacher' }

    wrapper = mountLayout()
    expect(wrapper.exists()).toBe(true)
    expect(wrapper.find('#main-content').exists()).toBe(true)
  })

  it('navigates teacher profile to /teacher/settings', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 1, name: 'Mario Rossi', role: 'teacher' }

    wrapper = mountLayout()
    
    // Trigger profile navigation directly via vm
    wrapper.vm.navigateToProfile()
    expect(mockPush).toHaveBeenCalledWith('/teacher/settings')
  })

  it('navigates student profile to /student/profile', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 2, name: 'Luigi Verdi', role: 'student' }

    wrapper = mountLayout()
    wrapper.vm.navigateToProfile()
    expect(mockPush).toHaveBeenCalledWith('/student/profile')
  })

  it('navigates admin notifications to /admin/dashboard', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 3, name: 'Admin User', role: 'admin' }

    wrapper = mountLayout()
    wrapper.vm.navigateToNotifications()
    expect(mockPush).toHaveBeenCalledWith('/admin/dashboard')
  })

  it('identifies teacher role properly including coordinator', async () => {
    const authStore = useAuthStore()
    authStore.user = { id: 4, name: 'Prof Coord', role: 'coordinator' }

    wrapper = mountLayout()
    expect(wrapper.vm.isTeacherRole).toBe(true)
    expect(wrapper.vm.isTeacherCoordinator).toBe(true)
  })
})
