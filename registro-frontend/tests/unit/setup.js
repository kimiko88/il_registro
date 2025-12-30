import { config } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, vi } from 'vitest'

// Mock Quasar
config.global.mocks = {
    $q: {
        notify: vi.fn(),
        dialog: vi.fn().mockReturnValue({ onOk: vi.fn(), onCancel: vi.fn() }),
        loading: { show: vi.fn(), hide: vi.fn() },
        screen: { gt: { xs: true } },
        dark: { isActive: false }
    },
    $t: (msg) => msg // Mock I18n
}

// Global stubs
config.global.stubs = {
    'q-icon': true,
    'q-btn': true,
    'q-btn-dropdown': true,
    'q-list': true,
    'q-item': true,
    'q-item-section': true,
    'q-item-label': true,
    'q-separator': true,
    'q-avatar': true,
    'q-card': true,
    'q-card-section': true,
    'q-card-actions': true,
    'q-input': true,
    'q-select': true,
    'q-spinner': true,
    'q-tooltip': true,
    'q-menu': true,
    'q-dialog': true,
    'q-popup-proxy': true,
    'q-date': true,
    'q-chip': true,
    'q-table': true,
    'q-th': true,
    'q-tr': true,
    'q-td': true,
    'q-badge': true,
    'q-page': true,
    'q-page-container': true,
    'q-layout': true,
    'q-header': true,
    'q-toolbar': true,
    'q-toolbar-title': true,
    'q-space': true,
    'q-drawer': true,
    'q-footer': true,
    'router-link': true,
    'router-view': true
}

beforeEach(() => {
    setActivePinia(createPinia())
})
