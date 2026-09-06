import XCTest
@testable import StudentApp

final class StudentViewModelTests: XCTestCase {
    var viewModel: StudentViewModel!

    override func setUp() {
        super.setUp()
        viewModel = StudentViewModel()
        viewModel.loadSampleData()
    }

    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }

    func testInitialStateIsEmpty() {
        let freshVM = StudentViewModel()
        XCTAssertTrue(freshVM.grades.isEmpty)
        XCTAssertTrue(freshVM.homework.isEmpty)
        XCTAssertTrue(freshVM.attendance.isEmpty)
    }

    func testCalculateGPA() {
        let gpa = viewModel.calculateGPA()
        // (8.5 + 7.0 + 8.0 + 9.0 + 7.0) / 5 = 39.5 / 5 = 7.9
        XCTAssertEqual(gpa, 7.9, accuracy: 0.01)
    }

    func testToggleHomework() {
        let initialState = viewModel.homework.first(where: { $0.id == "1" })?.isCompleted ?? true
        XCTAssertFalse(initialState)

        let toggled = viewModel.toggleHomework(id: "1")
        XCTAssertTrue(toggled)
        XCTAssertTrue(viewModel.homework.first(where: { $0.id == "1" })?.isCompleted ?? false)
    }
}
