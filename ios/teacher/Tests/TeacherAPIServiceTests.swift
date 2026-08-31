import XCTest
@testable import TeacherApp

final class TeacherAPIServiceTests: XCTestCase {
    var apiService: MockTeacherAPIService!

    override func setUp() {
        super.setUp()
        apiService = MockTeacherAPIService()
    }

    override func tearDown() {
        apiService = nil
        super.tearDown()
    }

    func testLoginSuccess() async throws {
        let token = try await apiService.login(email: "docente@scuola.it", password: "password123")
        XCTAssertEqual(token, "jwt_teacher_token_mock")
    }

    func testSignLessonHour() async throws {
        let success = try await apiService.signLessonHour(token: "token", classId: "3A", topic: "Matematica applicata")
        XCTAssertTrue(success)
    }

    func testSubmitGradeValidation() async throws {
        let valid = try await apiService.submitGrade(token: "token", studentId: "s1", grade: 8.5, type: "Scritto")
        XCTAssertTrue(valid)
    }
}
