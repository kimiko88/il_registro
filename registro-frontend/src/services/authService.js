import api from './api'

/**
  * Authentication service for handling all auth-related API calls
  */
export const authService = {
    /**
     * Login user with email and password
     * @param {string} email - User email
     * @param {string} password - User password
     * @returns {Promise} Response with user data and tokens
     */
    async login(email, password) {
        const sanitizedEmail = typeof email === 'string' ? email.trim() : email
        const response = await api.post('/auth/login', { email: sanitizedEmail, password })
        return response.data
    },

    /**
     * Logout user by revoking refresh token
     * @param {string} [refreshToken] - Refresh token to revoke
     * @returns {Promise} Response message
     */
    async logout(refreshToken) {
        const payload = refreshToken ? { refresh_token: refreshToken } : {}
        const response = await api.post('/auth/logout', payload)
        return response.data
    },

    /**
     * Get current user profile
     * @returns {Promise} User profile data
     */
    async getCurrentUser() {
        const response = await api.get('/auth/me')
        return response.data
    },

    /**
     * Refresh access token using refresh token
     * @param {string} refreshToken - Refresh token
     * @returns {Promise} New token pair
     */
    async refreshToken(refreshToken) {
        const response = await api.post('/auth/refresh-token', { refresh_token: refreshToken })
        return response.data
    },

    /**
     * Register new user
     * @param {Object} userData - User registration data
     * @returns {Promise} Created user data
     */
    async register(userData) {
        const payload = { ...userData }
        if (typeof payload.email === 'string') {
            payload.email = payload.email.trim()
        }
        const response = await api.post('/auth/register', payload)
        return response.data
    },

    /**
     * Change user password
     * @param {string} currentPassword - Current password
     * @param {string} newPassword - New password
     * @returns {Promise} Response data
     */
    async changePassword(currentPassword, newPassword) {
        const response = await api.post('/auth/change-password', {
            current_password: currentPassword,
            new_password: newPassword
        })
        return response.data
    },

    /**
     * Request password reset email
     * @param {string} email - User email
     * @returns {Promise} Response data
     */
    async forgotPassword(email) {
        const sanitizedEmail = typeof email === 'string' ? email.trim() : email
        const response = await api.post('/auth/forgot-password', { email: sanitizedEmail })
        return response.data
    },

    /**
     * Reset password using reset token
     * @param {string} token - Reset token
     * @param {string} newPassword - New password
     * @returns {Promise} Response data
     */
    async resetPassword(token, newPassword) {
        const response = await api.post('/auth/reset-password', {
            token,
            new_password: newPassword
        })
        return response.data
    }
}

export default authService

