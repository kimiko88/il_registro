import { describe, it, expect } from 'vitest'

// Pure unit tests for GradeWeights business logic

// Weight validation rule — same as the :rules on q-input in GradeWeights.vue
const weightRules = [
    val => val !== null && val !== '' || 'Campo obbligatorio',
    val => Number(val) >= 0 || 'Il peso deve essere >= 0',
    val => Number(val) <= 10 || 'Il peso deve essere <= 10',
]

function validateWeight(val) {
    for (const rule of weightRules) {
        const result = rule(val)
        if (result !== true) return result
    }
    return true
}

describe('GradeWeights weight validation rules', () => {
    it('accepts weight 0 (excluded from average)', () => {
        expect(validateWeight(0)).toBe(true)
    })

    it('accepts weight 1 (normal weight)', () => {
        expect(validateWeight(1)).toBe(true)
    })

    it('accepts weight 2 (double weight)', () => {
        expect(validateWeight(2)).toBe(true)
    })

    it('accepts weight 10 (max boundary)', () => {
        expect(validateWeight(10)).toBe(true)
    })

    it('rejects weight -1 (below zero)', () => {
        const result = validateWeight(-1)
        expect(result).not.toBe(true)
        expect(result).toContain('>= 0')
    })

    it('rejects weight 11 (above max)', () => {
        const result = validateWeight(11)
        expect(result).not.toBe(true)
        expect(result).toContain('<= 10')
    })

    it('rejects null weight (required)', () => {
        const result = validateWeight(null)
        expect(result).not.toBe(true)
        expect(result).toContain('obbligatorio')
    })

    it('rejects empty string weight (required)', () => {
        const result = validateWeight('')
        expect(result).not.toBe(true)
    })
})

// weightColor utility — mirrors the logic in GradeWeights.vue
const weightColor = (w) =>
    w === 0 ? 'grey' : w > 1.5 ? 'deep-orange' : w < 0.8 ? 'blue-grey' : 'primary'

describe('GradeWeights weightColor utility', () => {
    it('returns grey for weight 0 (excluded)', () => {
        expect(weightColor(0)).toBe('grey')
    })

    it('returns deep-orange for weight > 1.5 (high emphasis)', () => {
        expect(weightColor(2)).toBe('deep-orange')
    })

    it('returns blue-grey for weight < 0.8 (low emphasis)', () => {
        expect(weightColor(0.5)).toBe('blue-grey')
    })

    it('returns primary for weight = 1 (normal)', () => {
        expect(weightColor(1)).toBe('primary')
    })

    it('returns primary for weight 1.5 exactly (boundary)', () => {
        expect(weightColor(1.5)).toBe('primary')
    })
})

// categoryColor utility — mirrors logic in GradeWeights.vue
const categoryColor = (cat) =>
    ({ formative: 'blue', summative: 'deep-purple', practical: 'teal' }[cat] || 'grey')

describe('GradeWeights categoryColor utility', () => {
    it('returns blue for formative', () => {
        expect(categoryColor('formative')).toBe('blue')
    })

    it('returns deep-purple for summative', () => {
        expect(categoryColor('summative')).toBe('deep-purple')
    })

    it('returns teal for practical', () => {
        expect(categoryColor('practical')).toBe('teal')
    })

    it('returns grey for unknown category', () => {
        expect(categoryColor('unknown')).toBe('grey')
    })
})

// Save validation — same guards as in save() function
function validateSaveForm(form) {
    if (!form.grade_category || form.weight == null) {
        return 'Compila tutti i campi obbligatori'
    }
    if (form.weight < 0 || form.weight > 10) {
        return 'Il peso deve essere compreso tra 0 e 10'
    }
    return true
}

describe('GradeWeights save form validation', () => {
    it('passes with valid grade_category and weight', () => {
        expect(validateSaveForm({ grade_category: 'summative', weight: 1 })).toBe(true)
    })

    it('fails when grade_category is missing', () => {
        const result = validateSaveForm({ grade_category: '', weight: 1 })
        expect(result).toContain('obbligatori')
    })

    it('fails when weight is null', () => {
        const result = validateSaveForm({ grade_category: 'formative', weight: null })
        expect(result).toContain('obbligatori')
    })

    it('fails when weight is negative', () => {
        const result = validateSaveForm({ grade_category: 'formative', weight: -1 })
        expect(result).toContain('tra 0 e 10')
    })

    it('fails when weight exceeds 10', () => {
        const result = validateSaveForm({ grade_category: 'formative', weight: 11 })
        expect(result).toContain('tra 0 e 10')
    })
})
