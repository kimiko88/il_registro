import XCTest
@testable import TeacherApp

final class TeacherViewModelTests: XCTestCase {
    var viewModel: TeacherViewModel!

    override func setUp() {
        super.setUp()
        viewModel = TeacherViewModel()
    }

    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }

    func testSignLesson() {
        let emptyFailure = viewModel.signLesson(topic: "")
        XCTAssertFalse(emptyFailure)

        let success = viewModel.signLesson(topic: "Sistemi lineari e matrici")
        XCTAssertTrue(success)
        XCTAssertTrue(viewModel.isHourSigned)
        XCTAssertEqual(viewModel.lessonTopic, "Sistemi lineari e matrici")
    }

    func testUpdateAttendance() {
        let success = viewModel.updateAttendance(studentId: "s3", status: "presente")
        XCTAssertTrue(success)
        XCTAssertEqual(viewModel.students.first(where: { $0.id == "s3" })?.status, "presente")
    }

    func testInsertGradeValidation() {
        let valid = viewModel.insertGrade(studentId: "s1", grade: 8.5, type: "Scritto")
        XCTAssertTrue(valid)

        let invalid = viewModel.insertGrade(studentId: "s1", grade: 12.0, type: "Scritto")
        XCTAssertFalse(invalid)
    }
}
