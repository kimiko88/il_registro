package it.scuola.registro.parent

import it.scuola.registro.parent.cache.OfflineCacheManager
import it.scuola.registro.parent.network.HttpParentApiService
import it.scuola.registro.parent.security.BiometricAuthCallback
import it.scuola.registro.parent.security.BiometricAuthHelper
import it.scuola.registro.parent.security.SecurityCheckHelper
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

/**
 * End-to-End & Integration Test for Parent App against the REAL Live Backend and Database.
 * Zero mocks / Zero stubs: all parent credentials and children are fetched live from PostgreSQL.
 */
class ParentLiveDatabaseE2ETest {

    private lateinit var apiService: HttpParentApiService
    private lateinit var cacheManager: OfflineCacheManager
    private lateinit var securityHelper: SecurityCheckHelper
    private lateinit var biometricHelper: BiometricAuthHelper

    @Before
    fun setUp() {
        apiService = HttpParentApiService("https://registro-backend-fdu2.onrender.com/api/v1")
        cacheManager = OfflineCacheManager()
        securityHelper = SecurityCheckHelper()
        biometricHelper = BiometricAuthHelper()
    }

    @Test
    fun testLiveParentFullFlow_RealDatabaseIntegration() = runBlocking {
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

        // 3. Real Server Login with verified live database account
        val loginResult = apiService.login("genitore2b_10@scuola.it", "password")
        assertTrue("Real database login must succeed: ${loginResult.exceptionOrNull()?.message}", loginResult.isSuccess)

        val token = loginResult.getOrThrow()
        assertTrue("JWT token must be populated", token.isNotEmpty())

        // 4. Retrieve REAL live parent children from database
        val childrenResult = apiService.getChildren(token)
        assertTrue("Children fetch must succeed: ${childrenResult.exceptionOrNull()?.message}", childrenResult.isSuccess)
        val children = childrenResult.getOrThrow()
        assertTrue("Parent must have at least one child registered in database", children.isNotEmpty())

        val child = children.first()
        assertTrue("Child ID must be valid", child.id.isNotBlank())
        assertEquals("Studente2B_10", child.firstName)
        assertEquals("2B", child.className)

        // 5. Store real live children into OfflineCacheManager
        cacheManager.saveChildren(children)
        assertTrue("Cache must be valid after saving live children", cacheManager.isCacheValid())
        assertEquals("Cached count must match live count", children.size, cacheManager.getCachedChildren().size)

        // 6. Check absences for child from live server
        val absencesResult = apiService.getAbsences(token, child.id)
        assertTrue("Absences check must complete without exception", absencesResult.isSuccess)

        // 7. Clear cache
        cacheManager.clearCache()
        assertFalse("Cache must be invalidated after clear", cacheManager.isCacheValid())
    }
}
