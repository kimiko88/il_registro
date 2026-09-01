import XCTest
@testable import ParentApp

final class ParentIntegrationTests: XCTestCase {

    var apiService: HttpParentAPIService!
    var viewModel: ParentViewModel!

    override func setUp() {
        super.setUp()
        apiService = HttpParentAPIService()
        viewModel = ParentViewModel()
    }

    func testParentFullFlow() {
        XCTAssertNotNil(apiService)
        XCTAssertNotNil(viewModel)

        viewModel.loadSampleData()
        XCTAssertEqual(viewModel.children.count, 2)

        let child = viewModel.children.first!
        viewModel.selectChild(id: child.id)
        XCTAssertEqual(viewModel.selectedChildId, child.id)

        let justified = viewModel.justifyAbsence(id: "a1", reason: "Visita medica")
        XCTAssertTrue(justified)
    }
}
