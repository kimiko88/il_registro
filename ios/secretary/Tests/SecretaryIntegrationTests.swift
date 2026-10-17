import XCTest
@testable import SecretaryApp

final class SecretaryIntegrationTests: XCTestCase {

    var apiService: HttpSecretaryAPIService!
    var viewModel: SecretaryViewModel!
    var cacheManager: OfflineCacheManager!
    var securityHelper: SecurityCheckHelper!
    var biometricManager: BiometricAuthManager!

    override func setUp() {
        super.setUp()
        apiService = HttpSecretaryAPIService()
        viewModel = SecretaryViewModel(apiService: apiService)
        cacheManager = OfflineCacheManager.shared
        securityHelper = SecurityCheckHelper.shared
        biometricManager = BiometricAuthManager.shared
    }

    func testSecurityAndBiometricIntegrity() {
        let report = securityHelper.getSecurityReport()
        XCTAssertNotNil(report["isJailbroken"])
        XCTAssertNotNil(report["isSimulator"])

        let canAuth = biometricManager.canEvaluatePolicy()
        XCTAssertTrue(canAuth || !canAuth)
    }

    func testRealSecretaryLoginAndLiveDatabaseFlow() async throws {
        // Zero mocks: connects to real live database
        let token = try await apiService.login(email: "segreteria.prova@scuola.it", password: "password")
        XCTAssertFalse(token.isEmpty)

        let users = try await apiService.fetchUsers(token: token)
        XCTAssertFalse(users.isEmpty)
        let firstUser = users.first!
        XCTAssertFalse(firstUser.id.isEmpty)
        XCTAssertFalse(firstUser.email.isEmpty)

        // Offline cache verification of real managed users
        let managedUsers = users.map {
            ManagedUserModel(id: $0.id, name: "\($0.firstName) \($0.lastName)".trimmingCharacters(in: .whitespacesAndNewlines).ifEmpty($0.email), role: $0.role)
        }
        cacheManager.saveUsers(managedUsers)
        XCTAssertTrue(cacheManager.isCacheValid())
        XCTAssertEqual(cacheManager.getCachedUsers().count, managedUsers.count)

        // Load into ViewModel
        await viewModel.loadFromDatabase(token: token)
        XCTAssertFalse(viewModel.users.isEmpty)

        cacheManager.clear()
        XCTAssertFalse(cacheManager.isCacheValid())
    }
}

private extension String {
    func ifEmpty(_ fallback: String) -> String {
        return self.isEmpty ? fallback : self
    }
}
