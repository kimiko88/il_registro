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
                    'q-chip': { template: '<span class="q-chip-stub" :data-color="color"><slot /></span>', props: ['color'] },
                    'q-icon': true,
                    'q-btn': true,
                    'q-linear-progress': true,
                    'q-spinner': true
                }
            }
        })
    })

    it('renders monitoring page with overall healthy status and online service chips', async () => {
        await flushPromises()
        expect(wrapper.text()).toContain('Monitoraggio Sistema')
        expect(wrapper.text()).toContain('Tutti i servizi operativi')

        // Checks that chips render ONLINE when backend sends healthy
        const chips = wrapper.findAll('.q-chip-stub')
        expect(chips.length).toBe(4)
        chips.forEach(chip => {
            expect(chip.text()).toBe('ONLINE')
            expect(chip.attributes('data-color')).toBe('positive')
        })
    })

    it('renders degraded and in-memory states correctly', async () => {
        const api = (await import('@/services/api')).default
        api.get.mockImplementationOnce(() => Promise.resolve({
            data: {
                status: 'degraded',
                services: {
                    api: 'ok',
                    database: 'degraded',
                    redis: 'in-memory',
                    storage: 'unhealthy',
                    db_ping_ms: 120
                }
            }
        }))

        const degradedWrapper = mount(Monitoring, {
            global: {
                plugins: [
                    [Quasar, {}],
                    createTestingPinia({ createSpy: vi.fn })
                ],
                stubs: {
                    'q-page': { template: '<div><slot /></div>' },
                    'q-card': { template: '<div><slot /></div>' },
                    'q-card-section': { template: '<div><slot /></div>' },
                    'q-chip': { template: '<span class="q-chip-stub" :data-color="color"><slot /></span>', props: ['color'] },
                    'q-icon': true,
                    'q-btn': true,
                    'q-linear-progress': true,
                    'q-spinner': true
                }
            }
        })

        await flushPromises()
        expect(degradedWrapper.text()).toContain('Servizi degradati')
        const chips = degradedWrapper.findAll('.q-chip-stub')
        expect(chips[0].text()).toBe('ONLINE')
        expect(chips[0].attributes('data-color')).toBe('positive')
        expect(chips[1].text()).toBe('DEGRADATO')
        expect(chips[1].attributes('data-color')).toBe('warning')
        expect(chips[2].text()).toBe('IN-MEMORY')
        expect(chips[2].attributes('data-color')).toBe('info')
        expect(chips[3].text()).toBe('OFFLINE')
        expect(chips[3].attributes('data-color')).toBe('negative')
    })
})

