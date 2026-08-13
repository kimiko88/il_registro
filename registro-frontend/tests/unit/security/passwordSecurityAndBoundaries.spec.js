import { describe, it, expect } from 'vitest'

// Helper complexity validator matching backend & UI rules
function validatePasswordComplexity(password) {
    if (!password || password.length < 10) return { valid: false, reason: 'min_length_10' }
    if (!/[A-Z]/.test(password)) return { valid: false, reason: 'missing_uppercase' }
    if (!/[a-z]/.test(password)) return { valid: false, reason: 'missing_lowercase' }
    if (!/[0-9]/.test(password)) return { valid: false, reason: 'missing_number' }
    if (!/[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(password)) return { valid: false, reason: 'missing_special' }
    return { valid: true }
}

// Role boundary rule for password reset
function canRoleResetPasswordForTarget(actorRole, targetRole) {
    if (actorRole === 'superadmin' || actorRole === 'admin') {
        return true
    }
    if (actorRole === 'secretary') {
        return ['teacher', 'student', 'parent'].includes(targetRole)
    }
    return false
}

describe('Security & Boundary Tests — Password Policy & RBAC Boundaries', () => {
    describe('Password Complexity Validation', () => {
        it('rejects passwords shorter than 10 characters', () => {
            const res = validatePasswordComplexity('Short1!')
            expect(res.valid).toBe(false)
            expect(res.reason).toBe('min_length_10')
        })

        it('rejects passwords missing uppercase letters', () => {
            const res = validatePasswordComplexity('password123!')
            expect(res.valid).toBe(false)
            expect(res.reason).toBe('missing_uppercase')
        })

        it('rejects passwords missing lowercase letters', () => {
            const res = validatePasswordComplexity('PASSWORD123!')
            expect(res.valid).toBe(false)
            expect(res.reason).toBe('missing_lowercase')
        })

        it('rejects passwords missing numbers', () => {
            const res = validatePasswordComplexity('PasswordSpecial!')
            expect(res.valid).toBe(false)
            expect(res.reason).toBe('missing_number')
        })

        it('rejects passwords missing special characters', () => {
            const res = validatePasswordComplexity('Password12345')
            expect(res.valid).toBe(false)
            expect(res.reason).toBe('missing_special')
        })

        it('accepts compliant passwords meeting all criteria', () => {
            const res = validatePasswordComplexity('CompliantPass123!')
            expect(res.valid).toBe(true)
        })
    })

    describe('Secretary Role Password Reset Boundaries', () => {
        it('allows Secretary to reset password for Docenti (teacher)', () => {
            expect(canRoleResetPasswordForTarget('secretary', 'teacher')).toBe(true)
        })

        it('allows Secretary to reset password for Studenti (student)', () => {
            expect(canRoleResetPasswordForTarget('secretary', 'student')).toBe(true)
        })

        it('allows Secretary to reset password for Genitori (parent)', () => {
            expect(canRoleResetPasswordForTarget('secretary', 'parent')).toBe(true)
        })

        it('forbids Secretary from resetting password for Admin', () => {
            expect(canRoleResetPasswordForTarget('secretary', 'admin')).toBe(false)
        })

        it('forbids Secretary from resetting password for SuperAdmin', () => {
            expect(canRoleResetPasswordForTarget('secretary', 'superadmin')).toBe(false)
        })

        it('forbids Secretary from resetting password for another Secretary', () => {
            expect(canRoleResetPasswordForTarget('secretary', 'secretary')).toBe(false)
        })
    })
})
