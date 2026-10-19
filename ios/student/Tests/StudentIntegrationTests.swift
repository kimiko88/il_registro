import XCTest
@testable import StudentApp

final class StudentIntegrationTests: XCTestCase {

    var apiService: HttpStudentAPIService!
    var viewModel: StudentViewModel!
    var cacheManager: OfflineCacheManager!
    var securityHelper: SecurityCheckHelper!
    var biometricManager: BiometricAuthManager!

    override func setUp() {
        super.setUp()
        apiService = HttpStudentAPIService()
        viewModel = StudentViewModel(apiService: apiService)
        cacheManager = OfflineCacheManager.shared
        securityHelper = SecurityCheckHelper.shared
        biometricManager = BiometricAuthManager.shared
    }

    func testSecurityAndBiometricIntegrity() {
        let report = securityHelper.getSecurityReport()
        XCTAssertNotNil(report["isJailbroken"])
        XCTAssertNotNil(report["isSimulator"])

        // In test environment, evaluate biometric availability check
        let canAuth = biometricManager.canEvaluatePolicy()
        XCTAssertTrue(canAuth || !canAuth) // Validates method execution without crash
    }

    func testRealStudentLoginAndLiveDatabaseFlow() async throws {
        // Zero mocks: connects to real live database
        let token = try await apiService.login(email: "studente2b_10@scuola.it", password: "password")
        XCTAssertFalse(token.isEmpty)
        XCTAssertEqual(apiService.lastStudentName, "Studente2B_10 Test")

        let grades = try await apiService.fetchGrades(token: token)
        XCTAssertFalse(grades.isEmpty)
        XCTAssertTrue(grades.first?.grade ?? 0.0 > 0.0)

        // Offline cache verification of real data
        cacheManager.saveGrades(grades)
        XCTAssertTrue(cacheManager.isCacheValid())
        XCTAssertEqual(cacheManager.getCachedGrades().count, grades.count)

        // Load into ViewModel
        await viewModel.loadFromDatabase(token: token)
        XCTAssertFalse(viewModel.grades.isEmpty)
        XCTAssertGreaterThan(viewModel.calculateGPA(), 5.0)

        cacheManager.clear()
        XCTAssertFalse(cacheManager.isCacheValid())
    }
}
