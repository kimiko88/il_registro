import XCTest
@testable import StudentApp

final class StudentIntegrationTests: XCTestCase {
    var apiService: MockStudentAPIService!
    var cacheManager: OfflineCacheManager!
    var viewModel: StudentViewModel!

    override func setUp() {
        super.setUp()
        apiService = MockStudentAPIService()
        cacheManager = OfflineCacheManager()
        viewModel = StudentViewModel()
    }

    override func tearDown() {
        apiService = nil
        cacheManager = nil
        viewModel = nil
        super.tearDown()
    }

    func testFullStudentIntegrationPipeline() async throws {
        // 1. Authenticate
        let token = try await apiService.login(email: "mario@scuola.it", password: "password123")
        XCTAssertFalse(token.isEmpty)

        // 2. Fetch Grades & Cache
        let remoteGrades = try await apiService.fetchGrades(token: token)
        cacheManager.saveGrades(remoteGrades)
        XCTAssertTrue(cacheManager.isCacheValid())

        // 3. Load into ViewModel and verify GPA
        viewModel.loadData()
        let gpa = viewModel.calculateGPA()
        XCTAssertEqual(gpa, 7.9, accuracy: 0.01)

        // 4. Toggle homework status
        let toggled = viewModel.toggleHomework(id: "1")
        XCTAssertTrue(toggled)
    }
}
