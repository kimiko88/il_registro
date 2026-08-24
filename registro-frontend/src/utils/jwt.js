/**
 * Safely decodes a Base64URL encoded string with proper padding.
 * @param {string} str
 * @returns {string}
 */
export function base64UrlDecode(str) {
  if (!str) return ''
  let base64 = str.replace(/-/g, '+').replace(/_/g, '/')
  while (base64.length % 4 !== 0) {
    base64 += '='
  }
  try {
    return decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
  } catch {
    try {
      return atob(base64)
    } catch {
      return ''
    }
  }
}

/**
 * Parses payload object from a JWT token.
 * @param {string} token
 * @returns {Record<string, any> | null}
 */
export function parseJwt(token) {
  if (!token || typeof token !== 'string') return null
  const parts = token.split('.')
  if (parts.length !== 3) return null
  const jsonStr = base64UrlDecode(parts[1])
  if (!jsonStr) return null
  try {
    return JSON.parse(jsonStr)
  } catch {
    return null
  }
}

/**
 * Checks if a JWT token is expired (with an optional buffer in seconds).
 * @param {string} token
 * @param {number} bufferSeconds
 * @returns {boolean}
 */
export function isTokenExpired(token, bufferSeconds = 0) {
  const payload = parseJwt(token)
  if (!payload || !payload.exp) return true
  const currentTime = Math.floor(Date.now() / 1000)
  return payload.exp <= currentTime + bufferSeconds
}

/**
 * Extracts user role from JWT token payload.
 * @param {string} token
 * @returns {string | null}
 */
export function getRoleFromToken(token) {
  const payload = parseJwt(token)
  return payload?.role || payload?.user_role || null
}
