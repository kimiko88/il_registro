import { describe, it, expect } from 'vitest'

// Core DOM XSS Sanitizer function used across frontend preview components
function sanitizeHTMLContent(dirtyHTML) {
    if (!dirtyHTML) return ''
    return dirtyHTML
        .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
        .replace(/on\w+="[^"]*"/gi, '')
        .replace(/on\w+='[^']*'/gi, '')
        .replace(/javascript:[^\s"']*/gi, '#void')
}

describe('DOM XSS Sanitization & HTML Security', () => {
    it('strips inline script tags from user-provided HTML content', () => {
        const dirty = '<p>Nota di classe</p><script>alert("hacked")</script>'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean).not.toContain('<script>')
        expect(clean).toContain('<p>Nota di classe</p>')
    })

    it('removes inline JavaScript event handlers (onerror, onclick)', () => {
        const dirty = '<img src="x" onerror="alert(1)" onclick="doBadThing()">'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean).not.toContain('onerror')
        expect(clean).not.toContain('onclick')
    })

    it('neutralizes malicious javascript: URI schemes in links', () => {
        const dirty = '<a href="javascript:alert(document.cookie)">Clicca qui</a>'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean).not.toContain('javascript:')
        expect(clean).toContain('href="#void"')
    })
})
