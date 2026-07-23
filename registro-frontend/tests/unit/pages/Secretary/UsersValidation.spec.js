import { describe, it, expect, vi } from 'vitest'

// Pure unit tests for Italian Codice Fiscale validation logic
// and user form logic extracted from secretary/Users.vue

const CF_REGEX = /^[A-Z]{6}[0-9]{2}[A-Z][0-9]{2}[A-Z][0-9]{3}[A-Z]$/i

describe('Codice Fiscale Regex Validation', () => {
    it('accepts a valid uppercase Codice Fiscale', () => {
        const cf = 'RSSMRA85M01H501Z'
        expect(CF_REGEX.test(cf)).toBe(true)
    })

    it('accepts a valid lowercase Codice Fiscale (case-insensitive flag)', () => {
        const cf = 'rssmra85m01h501z'
        expect(CF_REGEX.test(cf)).toBe(true)
    })

    it('rejects a CF that is too short', () => {
        expect(CF_REGEX.test('RSSMRA85')).toBe(false)
    })

    it('rejects a CF that is too long', () => {
        expect(CF_REGEX.test('RSSMRA85M01H501ZX')).toBe(false)
    })

    it('rejects a CF with invalid character in first 6 position', () => {
        // First 6 must be letters
        expect(CF_REGEX.test('RSS1RA85M01H501Z')).toBe(false)
    })

    it('rejects a CF with invalid digits in position 7-8', () => {
        // Positions 7-8 must be digits
        expect(CF_REGEX.test('RSSMRAAAM01H501Z')).toBe(false)
    })

    it('rejects an empty string', () => {
        expect(CF_REGEX.test('')).toBe(false)
    })

    it('rejects a CF with special characters', () => {
        expect(CF_REGEX.test('RSSMR@85M01H501Z')).toBe(false)
    })
})

// Validation rule function (as used in the q-input :rules prop)
const cfValidationRule = (val) =>
    !val || CF_REGEX.test(val) || 'Formato Codice Fiscale non valido'

describe('Codice Fiscale Validation Rule', () => {
    it('passes validation for an empty string (optional field)', () => {
        expect(cfValidationRule('')).toBe(true)
    })

    it('passes validation for null (optional)', () => {
        expect(cfValidationRule(null)).toBe(true)
    })

    it('passes validation for a valid CF', () => {
        expect(cfValidationRule('RSSMRA85M01H501Z')).toBe(true)
    })

    it('returns error message for an invalid CF', () => {
        const result = cfValidationRule('INVALID')
        expect(result).toBe('Formato Codice Fiscale non valido')
    })
})

// Role options validation
const roleOptions = [
    { label: 'Studente', value: 'student' },
    { label: 'Docente', value: 'teacher' },
    { label: 'Genitore', value: 'parent' },
    { label: 'Segreteria', value: 'secretary' },
    { label: 'Dirigente', value: 'principal' },
    { label: 'Vice Dirigente', value: 'vice_principal' },
    { label: 'Admin', value: 'admin' },
]

describe('User Role Options', () => {
    it('includes all required roles', () => {
        const values = roleOptions.map(r => r.value)
        expect(values).toContain('student')
        expect(values).toContain('teacher')
        expect(values).toContain('parent')
        expect(values).toContain('secretary')
        expect(values).toContain('principal')
    })

    it('has unique role values', () => {
        const values = roleOptions.map(r => r.value)
        const unique = new Set(values)
        expect(unique.size).toBe(values.length)
    })

    it('all options have both label and value', () => {
        roleOptions.forEach(opt => {
            expect(opt.label).toBeTruthy()
            expect(opt.value).toBeTruthy()
        })
    })
})
