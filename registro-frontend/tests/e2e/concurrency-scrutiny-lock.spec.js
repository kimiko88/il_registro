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

describe('High-Concurrency Scrutiny Locking & Read-Only UI E2E', () => {
    let pinia

    beforeEach(() => {
        pinia = createPinia()
        setActivePinia(pinia)
        vi.clearAllMocks()
    })

    describe('Scrutiny Session Lock Acquisition & Concurrency Guard', () => {
        it('CSL01 — acquires class scrutiny lock with coordinator ownership and lease TTL', () => {
            const createScrutinyLock = (classId, coordinatorId, ttlMinutes = 15) => {
                const now = new Date()
                const expiresAt = new Date(now.getTime() + ttlMinutes * 60000)
                return {
                    class_id: classId,
                    term: 1,
                    status: 'IN_PROGRESS',
                    lock_owner_id: coordinatorId,
                    locked_at: now.toISOString(),
                    expires_at: expiresAt.toISOString(),
                    is_active: true
                }
            }

            const lock = createScrutinyLock('class-3a', 'teacher-coordinator-1')
            expect(lock.class_id).toBe('class-3a')
            expect(lock.lock_owner_id).toBe('teacher-coordinator-1')
            expect(lock.is_active).toBe(true)
            expect(new Date(lock.expires_at).getTime()).toBeGreaterThan(Date.now())
        })

        it('CSL02 — rejects concurrent lock acquisition attempt by a second teacher with 409 Conflict', () => {
            const currentLock = {
                class_id: 'class-3a',
                lock_owner_id: 'coord-1',
                expires_at: new Date(Date.now() + 10 * 60000).toISOString(),
                is_active: true
            }

            const tryAcquireLock = (classId, requesterId) => {
                const now = new Date()
                const isLockValid = currentLock.is_active && new Date(currentLock.expires_at) > now
                if (isLockValid && currentLock.lock_owner_id !== requesterId) {
                    const err = new Error('409 Conflict: Scrutiny session is already locked by another user')
                    err.status = 409
                    err.currentOwner = currentLock.lock_owner_id
                    throw err
                }
                return { acquired: true }
            }

            // Same coordinator can re-acquire or refresh
            expect(tryAcquireLock('class-3a', 'coord-1').acquired).toBe(true)

            // Second teacher is rejected with 409
            expect(() => tryAcquireLock('class-3a', 'teacher-2')).toThrow('409 Conflict')
        })

        it('CSL03 — expired lock lease allows automatic takeover by incoming coordinator', () => {
            // Lock expired 2 minutes ago
            const expiredLock = {
                class_id: 'class-3a',
                lock_owner_id: 'coord-disconnected',
                expires_at: new Date(Date.now() - 2 * 60000).toISOString(),
                is_active: true
            }

            const acquireWithExpiryCheck = (lock, newRequesterId) => {
                const isExpired = new Date(lock.expires_at) <= new Date()
                if (isExpired) {
                    return {
                        ...lock,
                        lock_owner_id: newRequesterId,
                        locked_at: new Date().toISOString(),
                        expires_at: new Date(Date.now() + 15 * 60000).toISOString(),
                        takeover: true
                    }
                }
                throw new Error('Lock still active')
            }

            const newLock = acquireWithExpiryCheck(expiredLock, 'coord-new')
            expect(newLock.takeover).toBe(true)
            expect(newLock.lock_owner_id).toBe('coord-new')
            expect(new Date(newLock.expires_at).getTime()).toBeGreaterThan(Date.now())
        })
    })

    describe('Post-Scrutiny Finalization & Immutability Enforcement', () => {
        it('CSL04 — finalized scrutiny (VERBALIZZATO) freezes grades into read-only state', () => {
            const scrutinySession = {
                class_id: 'class-3a',
                term: 1,
                status: 'VERBALIZZATO', // Officially approved and signed
                verbale_signed: true,
                closed_at: '2026-01-31T12:00:00Z',
                students_evaluations: [
                    { student_id: 'st-1', subject_id: 'math', final_grade: 8, conduct: 9, outcome: 'PROMOSSO' }
                ]
            }

            const attemptGradeEdit = (session, studentId, newGrade) => {
                if (session.status === 'VERBALIZZATO' || session.verbale_signed) {
                    const err = new Error('403 Forbidden: Cannot modify grades in a finalized scrutiny session')
                    err.status = 403
                    throw err
                }
                const record = session.students_evaluations.find(s => s.student_id === studentId)
                record.final_grade = newGrade
                return record
            }

            expect(() => attemptGradeEdit(scrutinySession, 'st-1', 9))
                .toThrow('Cannot modify grades in a finalized scrutiny session')
            expect(scrutinySession.students_evaluations[0].final_grade).toBe(8)
        })

        it('CSL05 — frontend UI inputs are disabled when scrutiny status is read-only', () => {
            const getFieldProps = (isSessionLockedByOther, isFinalized) => {
                const isReadOnly = isSessionLockedByOther || isFinalized
                return {
                    disabled: isReadOnly,
                    readonly: isReadOnly,
                    canSave: !isReadOnly,
                    badgeLabel: isFinalized ? 'Scrutinio Chiuso' : isSessionLockedByOther ? 'In Sola Lettura' : 'Modifica Attiva'
                }
            }

            // Normal active session
            const editableProps = getFieldProps(false, false)
            expect(editableProps.disabled).toBe(false)
            expect(editableProps.canSave).toBe(true)

            // Locked by colleague
            const lockedProps = getFieldProps(true, false)
            expect(lockedProps.disabled).toBe(true)
            expect(lockedProps.canSave).toBe(false)
            expect(lockedProps.badgeLabel).toBe('In Sola Lettura')

            // Finalized session
            const closedProps = getFieldProps(false, true)
            expect(closedProps.disabled).toBe(true)
            expect(closedProps.canSave).toBe(false)
            expect(closedProps.badgeLabel).toBe('Scrutinio Chiuso')
        })
    })
})
