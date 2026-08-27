/**
 * Security utilities for HTML sanitization and XSS prevention.
 */

/**
 * Strips potentially dangerous executable tags, event handlers, and javascript: links from HTML strings.
 * @param {string} dirtyHTML - Untrusted HTML string
 * @returns {string} Sanitized safe HTML
 */
export function sanitizeHTMLContent(dirtyHTML) {
    if (!dirtyHTML || typeof dirtyHTML !== 'string') return ''
    return dirtyHTML
        .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
        .replace(/<iframe\b[\s\S]*?(?:<\/iframe>|>)/gi, '')
        .replace(/<object\b[\s\S]*?(?:<\/object>|>)/gi, '')
        .replace(/<embed\b[\s\S]*?(?:<\/embed>|>)/gi, '')
        .replace(/on\w+\s*=\s*"[^"]*"/gi, '')
        .replace(/on\w+\s*=\s*'[^']*'/gi, '')
        .replace(/on\w+\s*=\s*[^\s>]+/gi, '')
        .replace(/href\s*=\s*["']?javascript:[^"'>\s]*/gi, 'href="#void"')
        .replace(/src\s*=\s*["']?javascript:[^"'>\s]*/gi, 'src="#void"')
        .replace(/data\s*:\s*text\/html/gi, 'data:text/plain')
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
