import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Analytics from '@/pages/admin/Analytics.vue'

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn().mockImplementation((url) => {
            if (url === '/admin/system/metrics') {
                return Promise.resolve({
                    data: {
                        api_success_rate: 99.8,
                        db_cpu_percent: 15.0,
                        cache_hit_rate: 94.5
                    }
                })
            }
            if (url === '/admin/analytics/user-growth') {
                return Promise.resolve({
                    data: [
                        { month: 'Gen', users: 100 },
                        { month: 'Feb', users: 150 }
                    ]
                })
            }
            return Promise.resolve({ data: {} })
        })
    }
}))

vi.mock('@/services/adminService', () => ({
    default: {
        getDashboardStats: vi.fn().mockResolvedValue({
            data: {
                total_schools: 5,
                total_users: 1200,
                total_students: 900,
                total_teachers: 100,
                recent_events: [
                    { id: 'ev-1', description: 'Accesso sistema', user_name: 'Admin', created_at: '2026-07-31T10:00:00Z' }
                ]
            }
        })
    }
}))

describe('Analytics.vue', () => {
    let wrapper

    beforeEach(() => {
        wrapper = mount(Analytics, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({ createSpy: vi.fn })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-icon': true,
                    'q-btn': true,
                    'q-linear-progress': true,
                    'q-list': { template: '<ul><slot /></ul>' },
                    'q-item': { template: '<li><slot /></li>' },
                    'q-item-section': { template: '<div><slot /></div>' },
                    'q-item-label': { template: '<span><slot /></span>' },
                    'q-spinner': true
                }
            }
        })
    })

    it('renders analytics title and system metrics cards', async () => {
        await flushPromises()
        expect(wrapper.text()).toContain('Analytics di Sistema')
        expect(wrapper.text()).toContain('Metriche di Infrastruttura')
        expect(wrapper.text()).toContain('Percentuale di Successo API')
    })
})

