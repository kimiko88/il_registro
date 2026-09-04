import XCTest
@testable import StudentApp

final class StudentIntegrationTests: XCTestCase {

    var apiService: HttpStudentAPIService!
    var viewModel: StudentViewModel!

    override func setUp() {
        super.setUp()
        apiService = HttpStudentAPIService()
        viewModel = StudentViewModel()
    }

    func testStudentFullFlow() {
        XCTAssertNotNil(apiService)
        XCTAssertNotNil(viewModel)

        viewModel.loadSampleData()
        let gpa = viewModel.calculateGPA()
        XCTAssertEqual(gpa, 7.9, accuracy: 0.01)

        let toggled = viewModel.toggleHomework(id: "1")
        XCTAssertTrue(toggled)
    }
}
