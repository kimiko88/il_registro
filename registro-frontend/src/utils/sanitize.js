/**
 * Security utilities for HTML sanitization and XSS prevention.
 * Uses an inert <template> DOM fragment and strict tag/attribute allowlists
 * instead of fragile HTML filtering regular expressions.
 */

const ALLOWED_TAGS = new Set([
    'a', 'abbr', 'b', 'blockquote', 'br', 'code', 'del', 'div', 'em',
    'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'hr', 'i', 'img', 'li', 'mark',
    'ol', 'p', 'pre', 's', 'small', 'span', 'strong', 'sub', 'sup',
    'table', 'tbody', 'td', 'th', 'thead', 'tr', 'u', 'ul'
])

const DROP_SUBTREE_TAGS = new Set([
    'base', 'button', 'embed', 'form', 'frame', 'frameset', 'iframe',
    'input', 'link', 'math', 'meta', 'noscript', 'object', 'script',
    'select', 'style', 'svg', 'template', 'textarea'
])

const ALLOWED_ATTRS = new Set([
    'alt', 'class', 'colspan', 'height', 'href', 'rel', 'rowspan',
    'src', 'target', 'title', 'width'
])

function stripControlAndSpace(str) {
    return Array.from(String(str || ''))
        .filter(ch => ch.charCodeAt(0) > 0x20)
        .join('')
        .toLowerCase()
}

function sanitizeURLValue(raw) {
    const value = String(raw || '').trim()
    const normalized = stripControlAndSpace(value)
    if (
        normalized.startsWith('javascript:') ||
        normalized.startsWith('vbscript:') ||
        normalized.startsWith('data:text/html')
    ) {
        return '#void'
    }
    if (
        normalized === '' ||
        normalized.startsWith('#') ||
        normalized.startsWith('/') ||
        normalized.startsWith('http://') ||
        normalized.startsWith('https://') ||
        normalized.startsWith('mailto:') ||
        normalized.startsWith('tel:') ||
        (normalized.startsWith('data:image/') && !normalized.startsWith('data:image/svg'))
    ) {
        return value
    }
    return null
}

function sanitizeNodeTree(parent) {
    for (const child of Array.from(parent.childNodes)) {
        if (child.nodeType === 8) {
            child.remove()
            continue
        }
        if (child.nodeType !== 1) {
            continue
        }

        const tag = String(child.tagName || '').toLowerCase()
        if (tag.includes('<') || tag.includes('>') || DROP_SUBTREE_TAGS.has(tag)) {
            child.remove()
            continue
        }

        if (!ALLOWED_TAGS.has(tag)) {
            // Recurse and unwrap harmless unknown elements while keeping sanitized children
            sanitizeNodeTree(child)
            while (child.firstChild) {
                parent.insertBefore(child.firstChild, child)
            }
            child.remove()
            continue
        }

        for (const attr of Array.from(child.attributes)) {
            const name = String(attr.name || '').toLowerCase()
            if (name.startsWith('on') || !ALLOWED_ATTRS.has(name)) {
                child.removeAttribute(attr.name)
                continue
            }
            if (name === 'href' || name === 'src') {
                const safeURL = sanitizeURLValue(attr.value)
                if (safeURL === null) {
                    child.removeAttribute(attr.name)
                } else {
                    child.setAttribute(name, safeURL)
                }
            }
        }

        if (tag === 'a' && child.getAttribute('target') === '_blank') {
            child.setAttribute('rel', 'noopener noreferrer')
        }

        sanitizeNodeTree(child)
    }
}

/**
 * Strips potentially dangerous executable tags, event handlers, and javascript: links from HTML strings.
 * @param {string} dirtyHTML - Untrusted HTML string
 * @returns {string} Sanitized safe HTML
 */
export function sanitizeHTMLContent(dirtyHTML) {
    if (!dirtyHTML || typeof dirtyHTML !== 'string') return ''
    if (typeof document === 'undefined') {
        return escapeHTML(dirtyHTML)
    }
    const template = document.createElement('template')
    template.innerHTML = dirtyHTML
    sanitizeNodeTree(template.content)
    return template.innerHTML
}

/**
 * Safely strips all HTML markup and returns plain text content.
 * @param {string} html - Potentially rich HTML string
 * @returns {string} Plain text without HTML tags
 */
export function stripHTMLTags(html) {
    if (!html || typeof html !== 'string') return ''
    if (typeof document !== 'undefined') {
        const template = document.createElement('template')
        template.innerHTML = sanitizeHTMLContent(html)
        const text = template.content.textContent || ''
        return text.replace(/[<>]/g, '')
    }
    // Pure character-by-character state machine fallback for non-DOM environments
    let out = ''
    let inTag = false
    for (let i = 0; i < html.length; i++) {
        const ch = html[i]
        if (ch === '<') {
            inTag = true
        } else if (ch === '>') {
            inTag = false
        } else if (!inTag) {
            out += ch
        }
    }
    return out
}

/**
 * Escapes HTML special characters to prevent XSS when interpolating text.
 * @param {string} text - Raw unescaped text
 * @returns {string} HTML-escaped string
 */
export function escapeHTML(text) {
    if (!text || typeof text !== 'string') return ''
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    }
    return text.replace(/[&<>"']/g, m => map[m])
}
