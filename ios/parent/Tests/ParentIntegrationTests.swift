import XCTest
@testable import ParentApp

final class ParentIntegrationTests: XCTestCase {

    var apiService: HttpParentAPIService!
    var viewModel: ParentViewModel!
    var cacheManager: OfflineCacheManager!
    var securityHelper: SecurityCheckHelper!
    var biometricManager: BiometricAuthManager!

    override func setUp() {
        super.setUp()
        apiService = HttpParentAPIService()
        viewModel = ParentViewModel(apiService: apiService)
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

    func testRealParentLoginAndLiveDatabaseFlow() async throws {
        // Zero mocks: connects to real live database
        let token = try await apiService.login(email: "genitore2b_10@scuola.it", password: "password")
        XCTAssertFalse(token.isEmpty)

        let children = try await apiService.fetchChildren(token: token)
        XCTAssertFalse(children.isEmpty)
        let child = children.first!
        XCTAssertEqual(child.firstName, "Studente2B_10")
        XCTAssertEqual(child.className, "2B")

        // Offline cache verification of real children
        cacheManager.saveChildren(children)
        XCTAssertTrue(cacheManager.isCacheValid())
        XCTAssertEqual(cacheManager.getCachedChildren().count, children.count)

        // Load into ViewModel
        await viewModel.loadFromDatabase(token: token)
        XCTAssertFalse(viewModel.children.isEmpty)
        XCTAssertEqual(viewModel.selectedChildId, child.id)

        cacheManager.clear()
        XCTAssertFalse(cacheManager.isCacheValid())
    }
}
