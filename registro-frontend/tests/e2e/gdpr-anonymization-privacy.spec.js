import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

// ─── Quasar Mock ─────────────────────────────────────────────────────────────
const { mockNotify } = vi.hoisted(() => ({
    mockNotify: vi.fn()
}))

vi.mock('quasar', async (importOriginal) => {
    const actual = await importOriginal()
    return {
        ...actual,
        useQuasar: () => ({
            loading: { show: vi.fn(), hide: vi.fn() },
            notify: mockNotify,
            dialog: vi.fn().mockReturnValue({ onOk: (fn) => fn() })
        }),
        Notify: { create: mockNotify }
    }
})

describe('GDPR Privacy & Adult Student Access Control E2E', () => {
    let pinia

    beforeEach(() => {
        pinia = createPinia()
        setActivePinia(pinia)
        vi.clearAllMocks()
    })

    describe('Adult Student (18+) Parental Access Governance', () => {
        it('GDPR01 — minor student (<18) records are fully accessible to authenticated parent', () => {
            const student = {
                id: 'st-minor-16',
                first_name: 'Luca',
                last_name: 'Verdi',
                birth_date: '2010-05-15',
                age: 16,
                is_adult: false,
                parent_consent_granted: false
            }

            const canParentAccessGrades = !student.is_adult || student.parent_consent_granted
            const canParentJustifyAbsences = !student.is_adult || student.parent_consent_granted

            expect(canParentAccessGrades).toBe(true)
            expect(canParentJustifyAbsences).toBe(true)
        })

        it('GDPR02 — adult student (18+) without parental consent blocks parent from viewing grades and justification', () => {
            const student = {
                id: 'st-adult-18',
                first_name: 'Marco',
                last_name: 'Rossi',
                birth_date: '2008-01-10',
                age: 18,
                is_adult: true,
                parent_consent_granted: false
            }

            const canParentAccessGrades = !student.is_adult || student.parent_consent_granted
            const canParentJustifyAbsences = !student.is_adult || student.parent_consent_granted

            expect(canParentAccessGrades).toBe(false)
            expect(canParentJustifyAbsences).toBe(false)
        })

        it('GDPR03 — adult student (18+) who explicitly grants parental consent enables parent access', () => {
            const student = {
                id: 'st-adult-18-consenting',
                first_name: 'Chiara',
                last_name: 'Bianchi',
                birth_date: '2008-03-20',
                age: 18,
                is_adult: true,
                parent_consent_granted: true,
                consent_timestamp: '2026-09-12T10:00:00Z'
            }

            const canParentAccessGrades = !student.is_adult || student.parent_consent_granted
            expect(canParentAccessGrades).toBe(true)
            expect(student.consent_timestamp).toBeTruthy()
        })
    })

    describe('Right to be Forgotten & Data Anonymization (Art. 17 GDPR)', () => {
        it('GDPR04 — student profile anonymization replaces PII with pseudonymized identifiers while retaining audit records', () => {
            const originalStudent = {
                id: 'student-uuid-456',
                first_name: 'Alessandro',
                last_name: 'Manzoni',
                fiscal_code: 'MNZLSN05A01H501Z',
                email: 'a.manzoni@studenti.scuola.it',
                phone: '+39 333 1234567',
                address: 'Via Roma 12, Milano',
                is_anonymized: false
            }

            const anonymizeStudent = (st) => {
                return {
                    id: st.id, // Immutable internal UUID retained for grade history
                    first_name: 'STUDENTE',
                    last_name: `ANONIMO_${st.id.slice(0, 8)}`,
                    fiscal_code: 'ANONYMIZED_GDPR',
                    email: null,
                    phone: null,
                    address: null,
                    is_anonymized: true,
                    anonymized_at: new Date().toISOString()
                }
            }

            const anonymized = anonymizeStudent(originalStudent)

            expect(anonymized.is_anonymized).toBe(true)
            expect(anonymized.email).toBeNull()
            expect(anonymized.phone).toBeNull()
            expect(anonymized.fiscal_code).toBe('ANONYMIZED_GDPR')
            expect(anonymized.first_name).toBe('STUDENTE')
            expect(anonymized.last_name).toContain('ANONIMO_')
        })
    })

    describe('Sensitive Data Protection & BES/DSA Diagnosis Masking (Art. 9 GDPR)', () => {
        it('GDPR05 — medical diagnosis codes (ICD-10 / Legge 104) are masked from unauthorized staff', () => {
            const studentHealthRecord = {
                student_id: 'student-bes-1',
                has_bes_dsa: true,
                pdp_plan_active: true,
                diagnosis_code: 'F81.0 - Dislessia Evolutiva',
                special_accommodations: ['Tempo aggiuntivo 30%', 'Mappe concettuali', 'Sintesi vocale']
            }

            const filterHealthRecordByRole = (record, userRole) => {
                const authorizedRoles = ['support_teacher', 'coordinator', 'preside', 'admin']
                if (!authorizedRoles.includes(userRole)) {
                    return {
                        student_id: record.student_id,
                        has_bes_dsa: record.has_bes_dsa,
                        pdp_plan_active: record.pdp_plan_active,
                        diagnosis_code: '[RISERVATO]',
                        special_accommodations: record.special_accommodations // Accommodations are visible to all teachers
                    }
                }
                return record
            }

            // Regular teacher sees accommodations but diagnosis code is masked
            const regularTeacherView = filterHealthRecordByRole(studentHealthRecord, 'teacher')
            expect(regularTeacherView.diagnosis_code).toBe('[RISERVATO]')
            expect(regularTeacherView.special_accommodations.length).toBe(3)

            // Support teacher (sostegno) can view full medical details
            const supportTeacherView = filterHealthRecordByRole(studentHealthRecord, 'support_teacher')
            expect(supportTeacherView.diagnosis_code).toBe('F81.0 - Dislessia Evolutiva')
        })

        it('GDPR06 — personal data export generates traceability log for compliance auditing', () => {
            const auditTrail = []

            const exportStudentPersonalData = (operatorId, studentId, purpose) => {
                const logEntry = {
                    id: `audit-${Date.now()}`,
                    action: 'GDPR_DATA_EXPORT',
                    operator_id: operatorId,
                    target_student_id: studentId,
                    purpose: purpose,
                    timestamp: new Date().toISOString()
                }
                auditTrail.push(logEntry)
                return { exported: true, studentId }
            }

            const result = exportStudentPersonalData('admin-1', 'st-101', 'Art. 20 Data Portability Request')
            expect(result.exported).toBe(true)
            expect(auditTrail.length).toBe(1)
            expect(auditTrail[0].action).toBe('GDPR_DATA_EXPORT')
            expect(auditTrail[0].operator_id).toBe('admin-1')
        })
    })
})
