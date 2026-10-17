package it.scuola.registro.secretary

import it.scuola.registro.secretary.cache.OfflineCacheManager
import it.scuola.registro.secretary.network.HttpSecretaryApiService
import it.scuola.registro.secretary.security.BiometricAuthCallback
import it.scuola.registro.secretary.security.BiometricAuthHelper
import it.scuola.registro.secretary.security.SecurityCheckHelper
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

/**
 * End-to-End & Integration Test for Secretary App against the REAL Live Backend and Database.
 * Zero mocks / Zero stubs: queries live users and system administration data from PostgreSQL.
 */
class SecretaryLiveDatabaseE2ETest {

    private lateinit var apiService: HttpSecretaryApiService
    private lateinit var cacheManager: OfflineCacheManager
    private lateinit var securityHelper: SecurityCheckHelper
    private lateinit var biometricHelper: BiometricAuthHelper

    @Before
    fun setUp() {
        apiService = HttpSecretaryApiService("https://registro-backend-fdu2.onrender.com/api/v1")
        cacheManager = OfflineCacheManager()
        securityHelper = SecurityCheckHelper()
        biometricHelper = BiometricAuthHelper()
    }

    @Test
    fun testLiveSecretaryFullFlow_RealDatabaseIntegration() = runBlocking {
        // 1. Device Security Integrity Checks
        val securityReport = securityHelper.getSecurityReport()
        assertNotNull(securityReport)
        assertTrue(securityReport.containsKey("isRooted"))
        assertTrue(securityReport.containsKey("isEmulator"))

        // 2. Biometric Authentication Helper
        var biometricSuccess = false
        biometricHelper.authenticate(object : BiometricAuthCallback {
            override fun onSuccess() {
                biometricSuccess = true
            }
            override fun onError(errorCode: Int, errString: String) {
                biometricSuccess = false
            }
            override fun onFailed() {
                biometricSuccess = false
            }
        })
        assertTrue("Biometric authentication helper must function", biometricSuccess)

        // 3. Real Server Login with verified live secretary account
        val loginResult = apiService.login("segreteria.prova@scuola.it", "password")
        assertTrue("Real database login must succeed: ${loginResult.exceptionOrNull()?.message}", loginResult.isSuccess)

        val (token, profile) = loginResult.getOrThrow()
        assertTrue("JWT token must be populated", token.isNotEmpty())
        assertNotNull("Profile must be parsed", profile)
        assertEquals("secretary", profile?.role)

        // 4. Retrieve REAL live users managed in the school database
        val usersResult = apiService.getUsers(token)
        assertTrue("Users fetch must succeed: ${usersResult.exceptionOrNull()?.message}", usersResult.isSuccess)
        val users = usersResult.getOrThrow()
        assertTrue("School must have real users in database", users.isNotEmpty())

        val firstUser = users.first()
        assertTrue("User ID must be valid", firstUser.id.isNotBlank())
        assertTrue("User email must be valid", firstUser.email.isNotBlank())

        // 5. Retrieve certificates from database
        val certsResult = apiService.getCertificates(token)
        assertTrue("Certificates fetch must complete successfully", certsResult.isSuccess)

        // 6. Store real live users into OfflineCacheManager
        cacheManager.saveUsers(users)
        assertTrue("Cache must be valid after saving live users", cacheManager.isCacheValid())
        assertEquals("Cached count must match live count", users.size, cacheManager.getCachedUsers().size)

        // 7. Clear cache
        cacheManager.clearCache()
        assertFalse("Cache must be invalidated after clear", cacheManager.isCacheValid())
    }
}
