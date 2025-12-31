import api from './api'

/**
 * Authentication service for handling all auth-related API calls
 */
export default {
    /**
     * Login user with email and password
     * @param {string} email - User email
     * @param {string} password - User password
     * @returns {Promise} Response with user data and tokens
     */
    async login(email, password) {
        const response = await api.post('/auth/login', { email, password })
        return response.data
    },

    /**
     * Logout user by revoking refresh token
     * @param {string} refreshToken - Refresh token to revoke
     * @returns {Promise} Response message
     */
    async logout(refreshToken) {
        const response = await api.post('/auth/logout', { refresh_token: refreshToken })
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
        const response = await api.post('/auth/register', userData)
        return response.data
    }
}
