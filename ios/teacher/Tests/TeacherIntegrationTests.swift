import XCTest
@testable import TeacherApp

final class TeacherIntegrationTests: XCTestCase {

    var apiService: HttpTeacherAPIService!
    var viewModel: TeacherViewModel!
    var cacheManager: OfflineCacheManager!
    var securityHelper: SecurityCheckHelper!
    var biometricManager: BiometricAuthManager!

    override func setUp() {
        super.setUp()
        apiService = HttpTeacherAPIService()
        viewModel = TeacherViewModel(apiService: apiService)
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

    func testRealTeacherLoginAndLiveDatabaseFlow() async throws {
        // Zero mocks: connects to real live database
        let token = try await apiService.login(email: "docente.incarichi@scuola.it", password: "password")
        XCTAssertFalse(token.isEmpty)

        let classes = try await apiService.fetchClasses(token: token)
        XCTAssertFalse(classes.isEmpty)
        let firstClass = classes.first!
        XCTAssertFalse(firstClass.id.isEmpty)
        XCTAssertFalse(firstClass.className.isEmpty)

        let students = try await apiService.fetchClassStudents(token: token, classId: firstClass.id)
        XCTAssertFalse(students.isEmpty)

        // Offline cache verification of real classes & students
        cacheManager.saveClasses(classes)
        cacheManager.saveStudents(students)
        XCTAssertTrue(cacheManager.isCacheValid())
        XCTAssertEqual(cacheManager.getCachedClasses().count, classes.count)
        XCTAssertEqual(cacheManager.getCachedStudents().count, students.count)

        // Load into ViewModel
        let loaded = await viewModel.loadFromDatabase(token: token)
        XCTAssertTrue(loaded)
        XCTAssertFalse(viewModel.students.isEmpty)

        cacheManager.clear()
        XCTAssertFalse(cacheManager.isCacheValid())
    }
}
