import { describe, it, expect } from 'vitest'
import { sanitizeHTMLContent, escapeHTML } from '@/utils/sanitize'

describe('DOM XSS Sanitization & HTML Security', () => {
    it('strips inline script tags from user-provided HTML content', () => {
        const dirty = '<p>Nota di classe</p><script>alert("hacked")</script>'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean).not.toContain('<script>')
        expect(clean).toContain('<p>Nota di classe</p>')
    })

    it('strips iframes, objects, and embeds from user HTML', () => {
        const dirty = '<p>Test</p><iframe src="https://evil.com"></iframe><object data="bad"></object><embed src="bad">'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean).not.toContain('<iframe')
        expect(clean).not.toContain('<object')
        expect(clean).not.toContain('<embed')
        expect(clean).toContain('<p>Test</p>')
    })

    it('removes inline JavaScript event handlers (onerror, onclick, onload)', () => {
        const dirty = '<img src="x" onerror="alert(1)" onclick="doBadThing()" onload=evil()>'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean).not.toContain('onerror')
        expect(clean).not.toContain('onclick')
        expect(clean).not.toContain('onload')
    })

    it('neutralizes malicious javascript: URI schemes in links and src', () => {
        const dirty = '<a href="javascript:alert(document.cookie)">Clicca qui</a><img src="javascript:alert(1)">'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean).not.toContain('javascript:')
        expect(clean).toContain('href="#void"')
        expect(clean).toContain('src="#void"')
    })

    it('escapes HTML special characters using escapeHTML', () => {
        const raw = '<script>alert("XSS & Injection")</script>'
        const escaped = escapeHTML(raw)

        expect(escaped).toBe('&lt;script&gt;alert(&quot;XSS &amp; Injection&quot;)&lt;/script&gt;')
    })
})

