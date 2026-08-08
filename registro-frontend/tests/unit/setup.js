import { config } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, vi } from 'vitest'

const mockIconSet = {
    name: 'material-icons',
    type: { positive: 'check', negative: 'warning', info: 'info', warning: 'priority_high' },
    arrow: { dropdown: 'arrow_drop_down', expand: 'arrow_drop_down', down: 'arrow_drop_down', up: 'arrow_drop_up' },
    chevron: { left: 'chevron_left', right: 'chevron_right' },
    table: { arrow: 'arrow_upward', select: 'arrow_drop_down' },
    editor: {},
    tree: {}
}

// Mock Quasar module globally
vi.mock('quasar', async () => {
    const actual = await vi.importActual('quasar')
    return {
        ...actual,
        useQuasar: () => ({
            dark: { isActive: false },
            notify: vi.fn(),
            loading: { show: vi.fn(), hide: vi.fn() },
            dialog: vi.fn(() => ({
                onOk: vi.fn(callback => { if (callback) callback(); return { onCancel: vi.fn(), onDismiss: vi.fn() } }),
                onCancel: vi.fn(),
                onDismiss: vi.fn()
            })),
            screen: { lt: { md: false }, gt: { xs: true } },
            lang: { current: 'it' },
            iconSet: mockIconSet
        })
    }
})

// Provide $q globally
const mockQ = {
    dark: { isActive: false },
    notify: vi.fn(),
    loading: { show: vi.fn(), hide: vi.fn() },
    dialog: vi.fn(() => ({ 
        onOk: vi.fn(callback => { if(callback) callback(); return { onCancel: vi.fn(), onDismiss: vi.fn() } }),
        onCancel: vi.fn(),
        onDismiss: vi.fn()
    })),
    screen: { lt: { md: false }, gt: { xs: true } },
    lang: { current: 'it' },
    iconSet: mockIconSet
}

config.global.provide = {
    $q: mockQ
}

config.global.mocks = {
    $q: mockQ,
    $t: (msg) => msg
}

// Global stubs for Quasar components
config.global.stubs = {
    'q-btn': { template: '<button @click="$emit(\'click\', $event)"><slot /></button>' },
    'q-input': { template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />', props: ['modelValue'] },
    'q-select': { template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><slot /></select>', props: ['modelValue'] },
    'q-checkbox': { template: '<input type="checkbox" :checked="modelValue" @change="$emit(\'update:modelValue\', $event.target.checked)" />', props: ['modelValue'] },
    'q-icon': true,
    'q-badge': true,
    'q-chip': true,
    'q-avatar': true,
    'q-spinner': true,
    'q-tooltip': true,
    'q-card': true,
    'q-card-section': true,
    'q-card-actions': true,
    'q-dialog': true,
    'q-page': true,
    'q-toolbar': true,
    'q-header': true,
    'q-footer': true,
    'q-drawer': true,
    'q-layout': true,
    'q-list': true,
    'q-item': true,
    'q-item-section': true,
    'q-item-label': true,
    'q-separator': true,
    'q-popup-proxy': true,
    'q-date': true,
    'q-time': true,
    'q-form': { template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>' },
    'q-tr': { template: '<tr><slot /></tr>' },
    'q-td': { template: '<td><slot /></td>' },
    'router-link': true,
    'router-view': true
}

beforeEach(() => {
    setActivePinia(createPinia())
})
