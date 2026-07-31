import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { Quasar } from 'quasar'
import Monitoring from '@/pages/admin/Monitoring.vue'

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn().mockImplementation((url) => {
            if (url === '/admin/system/health') {
                return Promise.resolve({
                    data: {
                        status: 'healthy',
                        services: {
                            api: 'healthy',
                            database: 'healthy',
                            redis: 'healthy',
                            storage: 'healthy',
                            db_ping_ms: 2,
                            redis_ping_ms: 1
                        },
                        metrics: {
                            cpu_percent: 12,
                            memory_percent: 35,
                            disk_percent: 28,
                            api_latency_ms: 14
                        },
                        api_version: '1.0.0',
                        db_version: 'PostgreSQL 15',
                        environment: 'production',
                        uptime: '1d 4h 10m'
                    }
                })
            }
            return Promise.resolve({ data: {} })
        })
    }
}))

describe('Monitoring.vue', () => {
    let wrapper

    beforeEach(() => {
        wrapper = mount(Monitoring, {
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
                    'q-chip': true,
                    'q-spinner': true
                }
            }
        })
    })

    it('renders monitoring page with overall healthy status', async () => {
        await flushPromises()
        expect(wrapper.text()).toContain('Monitoraggio Sistema')
        expect(wrapper.text()).toContain('Tutti i servizi operativi')
    })
})

