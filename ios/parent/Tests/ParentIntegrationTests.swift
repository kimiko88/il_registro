import XCTest
@testable import ParentApp

final class ParentIntegrationTests: XCTestCase {
    var apiService: MockParentAPIService!
    var viewModel: ParentViewModel!

    override func setUp() {
        super.setUp()
        apiService = MockParentAPIService()
        viewModel = ParentViewModel()
    }

    override func tearDown() {
        apiService = nil
        viewModel = nil
        super.tearDown()
    }

    func testParentIntegrationFlow() async throws {
        // 1. Authenticate Parent
        let token = try await apiService.login(email: "genitore@famiglia.it", password: "password123")
        XCTAssertFalse(token.isEmpty)

        // 2. Fetch children
        let children = try await apiService.fetchChildren(token: token)
        XCTAssertEqual(children.count, 2)

        // 3. Select child and justify absence
        viewModel.loadData()
        viewModel.selectChild(id: "c1")
        let absences = viewModel.getAbsencesForSelectedChild()
        XCTAssertEqual(absences.count, 2)

        let justifiedApi = try await apiService.justifyAbsence(token: token, absenceId: absences[0].id, note: "Visita medica")
        XCTAssertTrue(justifiedApi)

        let justifiedVm = viewModel.justifyAbsence(id: absences[0].id, note: "Visita medica")
        XCTAssertTrue(justifiedVm)
        XCTAssertTrue(viewModel.getAbsencesForSelectedChild()[0].isJustified)
    }
}
