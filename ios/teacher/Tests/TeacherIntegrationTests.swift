import XCTest
@testable import TeacherApp

final class TeacherIntegrationTests: XCTestCase {
    var apiService: MockTeacherAPIService!
    var viewModel: TeacherViewModel!

    override func setUp() {
        super.setUp()
        apiService = MockTeacherAPIService()
        viewModel = TeacherViewModel()
    }

    override func tearDown() {
        apiService = nil
        viewModel = nil
        super.tearDown()
    }

    func testTeacherIntegrationFlow() async throws {
        // 1. Authenticate Teacher
        let token = try await apiService.login(email: "docente@scuola.it", password: "password123")
        XCTAssertFalse(token.isEmpty)

        // 2. Sign Lesson Hour
        let signResult = try await apiService.signLessonHour(token: token, classId: "3A", topic: "Matematica")
        XCTAssertTrue(signResult)

        viewModel.loadData()
        let vmSign = viewModel.signLesson(topic: "Matematica")
        XCTAssertTrue(vmSign)
        XCTAssertTrue(viewModel.isHourSigned)

        // 3. Roll Call update
        let updated = viewModel.updateAttendance(studentId: "s3", status: "presente")
        XCTAssertTrue(updated)

        // 4. Grade insertion
        let gradeResult = try await apiService.submitGrade(token: token, studentId: "s1", grade: 8.5, type: "Scritto")
        XCTAssertTrue(gradeResult)

        let vmGrade = viewModel.insertGrade(studentId: "s1", grade: 8.5, type: "Scritto")
        XCTAssertTrue(vmGrade)
    }
}
