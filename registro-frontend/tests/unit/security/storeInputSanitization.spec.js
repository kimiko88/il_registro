import { describe, it, expect } from 'vitest'

function sanitizeStoreInput(input) {
    if (typeof input !== 'string') return input
    const trimmed = input.trim()
    return trimmed
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;')
}

describe('Frontend Store Payload Input Sanitization', () => {
    it('trims leading and trailing whitespace from string payloads', () => {
        const raw = '   Compito di Matematica   '
        const sanitized = sanitizeStoreInput(raw)

        expect(sanitized).toBe('Compito di Matematica')
    })

    it('escapes HTML special characters to prevent stored XSS vulnerabilities', () => {
        const raw = '<script>alert("hack")</script>'
        const sanitized = sanitizeStoreInput(raw)

        expect(sanitized).not.toContain('<script>')
        expect(sanitized).toBe('&lt;script&gt;alert(&quot;hack&quot;)&lt;/script&gt;')
    })
})
