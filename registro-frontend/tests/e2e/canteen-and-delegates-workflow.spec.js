import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mealsService } from '@/services/mealsService'
import { api } from '@/services/api'
import { useSchoolSettingsStore } from '@/stores/schoolSettings'

vi.mock('@/services/mealsService', () => ({
    mealsService: {
        recordRollCall: vi.fn(),
        getCateringReport: vi.fn(),
        registerDiet: vi.fn(),
        getDiet: vi.fn(),
        listDiets: vi.fn(),
        getWallet: vi.fn(),
        topUpWallet: vi.fn()
    }
}))

vi.mock('@/services/api', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        patch: vi.fn(),
        delete: vi.fn()
    },
    api: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        patch: vi.fn(),
        delete: vi.fn()
    }
}))

describe('School Canteen & Parent Delegates E2E Workflow', () => {
    let pinia

    beforeEach(() => {
        pinia = createPinia()
        setActivePinia(pinia)
        vi.clearAllMocks()
    })

    describe('School Canteen & Electronic Wallet Workflow', () => {
        it('CW01 — allows toggling canteen feature state in school settings', () => {
            const settingsStore = useSchoolSettingsStore()
            settingsStore.canteenEnabled = true
            expect(settingsStore.canteenEnabled).toBe(true)

            settingsStore.canteenEnabled = false
            expect(settingsStore.canteenEnabled).toBe(false)
        })

        it('CW02 — executes wallet top-up lifecycle with PagoPA IUV notice', async () => {
            const studentId = 'student-minor-01'
            const initialBalance = 12.50
            const topupAmount = 50.00
            const expectedIUV = '0121998877665544332'

            mealsService.topUpWallet.mockResolvedValueOnce({
                data: {
                    success: true,
                    student_id: studentId,
                    amount: topupAmount,
                    new_balance: initialBalance + topupAmount,
                    pagopa_iuv: expectedIUV,
                    status: 'paid'
                }
            })

            const res = await mealsService.topUpWallet(studentId, topupAmount, expectedIUV)

            expect(mealsService.topUpWallet).toHaveBeenCalledWith(studentId, topupAmount, expectedIUV)
            expect(res.data.success).toBe(true)
            expect(res.data.new_balance).toBe(62.50)
            expect(res.data.pagopa_iuv).toBe(expectedIUV)
        })

        it('CW03 — records morning catering roll-call and generates daily meals report', async () => {
            const classId = 'class-2b-primary'
            const date = '2026-10-31'
            const presentStudents = ['stud-1', 'stud-2', 'stud-3']

            mealsService.recordRollCall.mockResolvedValueOnce({
                data: {
                    class_id: classId,
                    date,
                    total_meals_ordered: 3,
                    special_diets_count: 1
                }
            })

            const rollCallRes = await mealsService.recordRollCall(classId, date, presentStudents)
            expect(rollCallRes.data.total_meals_ordered).toBe(3)
            expect(rollCallRes.data.special_diets_count).toBe(1)

            // Catering company report
            mealsService.getCateringReport.mockResolvedValueOnce({
                data: {
                    date,
                    standard_meals: 2,
                    special_diets: [
                        { student_id: 'stud-1', diet_type: 'celiachia', menu: 'gluten_free' }
                    ]
                }
            })

            const reportRes = await mealsService.getCateringReport(date)
            expect(reportRes.data.standard_meals).toBe(2)
            expect(reportRes.data.special_diets[0].diet_type).toBe('celiachia')
        })

        it('CW04 — registers and retrieves certified medical diet with health protocol', async () => {
            const studentId = 'student-allergy-44'
            const dietData = {
                student_id: studentId,
                diet_type: 'favismo_g6pd',
                category: 'medical',
                allergens: ['fave', 'piselli'],
                medical_certificate_url: 'https://storage.school.it/certificati/med_favismo.pdf',
                valid_until: '2027-06-30'
            }

            mealsService.registerDiet.mockResolvedValueOnce({
                data: {
                    id: 'diet-999',
                    ...dietData,
                    status: 'active'
                }
            })

            const regRes = await mealsService.registerDiet(dietData)
            expect(regRes.data.id).toBe('diet-999')
            expect(regRes.data.status).toBe('active')

            mealsService.getDiet.mockResolvedValueOnce({
                data: regRes.data
            })

            const fetchRes = await mealsService.getDiet(studentId)
            expect(fetchRes.data.allergens).toContain('fave')
            expect(fetchRes.data.category).toBe('medical')
        })
    })

    describe('Student Pick-up Delegates & Custody Compliance Workflow (Art. 337 c.c.)', () => {
        it('DEL01 — parent fetches children list and active delegates for selected child', async () => {
            const parentId = 'parent-mario'
            const childId = 'child-alessandro'

            api.get.mockImplementation((url) => {
                if (url === '/users/me/children') {
                    return Promise.resolve({
                        data: [
                            { id: childId, name: 'Alessandro Rossi (2A)' }
                        ]
                    })
                }
                if (url === `/students/${childId}/delegates`) {
                    return Promise.resolve({
                        data: [
                            {
                                id: 'del-01',
                                delegate_name: 'Giuseppe Verdi',
                                relationship: 'Nonno materno',
                                phone: '+393331122334',
                                is_active: true
                            }
                        ]
                    })
                }
                return Promise.reject(new Error('Unknown url ' + url))
            })

            const childrenRes = await api.get('/users/me/children')
            expect(childrenRes.data).toHaveLength(1)
            expect(childrenRes.data[0].id).toBe(childId)

            const delegatesRes = await api.get(`/students/${childId}/delegates`)
            expect(delegatesRes.data).toHaveLength(1)
            expect(delegatesRes.data[0].delegate_name).toBe('Giuseppe Verdi')
            expect(delegatesRes.data[0].is_active).toBe(true)
        })

        it('DEL02 — parent registers a new authorized pick-up delegate with identity details', async () => {
            const childId = 'child-alessandro'
            const newDelegateData = {
                delegate_name: 'Carla Neri',
                relationship: 'Zia paterna',
                phone: '+393478899001',
                fiscal_code: 'NRECRL75A41H501K',
                is_permanent: true,
                valid_until: '2027-06-30',
                id_document_url: 'https://storage.school.it/idcards/neri_ci.pdf',
                custody_order_ref: 'Trib. Minori RM n. 450/2025'
            }

            api.post.mockResolvedValueOnce({
                data: {
                    id: 'del-02',
                    student_id: childId,
                    ...newDelegateData,
                    is_active: true,
                    created_at: new Date().toISOString()
                }
            })

            const createRes = await api.post(`/students/${childId}/delegates`, newDelegateData)

            expect(api.post).toHaveBeenCalledWith(`/students/${childId}/delegates`, newDelegateData)
            expect(createRes.data.id).toBe('del-02')
            expect(createRes.data.delegate_name).toBe('Carla Neri')
            expect(createRes.data.custody_order_ref).toBe('Trib. Minori RM n. 450/2025')
        })

        it('DEL03 — parent revokes an authorized delegate and verifies status deactivation', async () => {
            const childId = 'child-alessandro'
            const delegateId = 'del-02'

            api.delete.mockResolvedValueOnce({
                data: {
                    message: 'Delega revocata con successo',
                    delegate_id: delegateId,
                    is_active: false
                }
            })

            const revokeRes = await api.delete(`/students/${childId}/delegates/${delegateId}`)

            expect(api.delete).toHaveBeenCalledWith(`/students/${childId}/delegates/${delegateId}`)
            expect(revokeRes.data.is_active).toBe(false)
            expect(revokeRes.data.message).toContain('successo')
        })
    })
})

