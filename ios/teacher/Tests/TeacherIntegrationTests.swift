import XCTest
@testable import TeacherApp

final class TeacherIntegrationTests: XCTestCase {

    var apiService: HttpTeacherAPIService!
    var viewModel: TeacherViewModel!

    override func setUp() {
        super.setUp()
        apiService = HttpTeacherAPIService()
        viewModel = TeacherViewModel()
    }

    func testTeacherFullFlow() {
        XCTAssertNotNil(apiService)
        XCTAssertNotNil(viewModel)

        viewModel.loadSampleSession()
        XCTAssertNotNil(viewModel.currentSession)

        let signed = viewModel.signLesson(topic: "Studio di funzioni")
        XCTAssertTrue(signed)
        XCTAssertTrue(viewModel.currentSession?.isSigned == true)

        let toggled = viewModel.toggleAttendance(studentId: "st1", status: "Assente")
        XCTAssertTrue(toggled)
    }
}
