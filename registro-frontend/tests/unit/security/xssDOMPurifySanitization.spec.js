import { describe, it, expect } from 'vitest'
import { sanitizeHTMLContent, stripHTMLTags, escapeHTML } from '@/utils/sanitize'

describe('DOM XSS Sanitization & HTML Security', () => {
    it('strips inline script tags from user-provided HTML content', () => {
        const dirty = '<p>Nota di classe</p><script>alert("hacked")</script>'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean).not.toContain('<script>')
        expect(clean).toContain('<p>Nota di classe</p>')
    })

    it('neutralizes nested multi-character bypass payloads (<scr<script>ipt>)', () => {
        const dirty = '<scr<script>ipt>alert(1)</scr</script>ipt><img src=x ononload=load=alert(1)>'
        const clean = sanitizeHTMLContent(dirty)

        expect(clean.toLowerCase()).not.toContain('<script')
        expect(clean.toLowerCase()).not.toContain('onload')
        expect(clean.toLowerCase()).not.toContain('ononload')
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

    it('strips all HTML tags safely using stripHTMLTags', () => {
        const cleanText = stripHTMLTags('<p>Avviso</p><script>alert(1)</script><b>Importante</b>')
        expect(cleanText).not.toContain('<')
        expect(cleanText).toContain('Avviso')
        expect(cleanText).toContain('Importante')

        const bypassAttempt = stripHTMLTags('Avviso<scr<script>ipt>alert(1)</scr</script>ipt>!')
        expect(bypassAttempt).not.toContain('<')
        expect(bypassAttempt.toLowerCase()).not.toContain('<script')
    })

    it('escapes HTML special characters using escapeHTML', () => {
        const raw = '<script>alert("XSS & Injection")</script>'
        const escaped = escapeHTML(raw)

        expect(escaped).toBe('&lt;script&gt;alert(&quot;XSS &amp; Injection&quot;)&lt;/script&gt;')
    })
})

