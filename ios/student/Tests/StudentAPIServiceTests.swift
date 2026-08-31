import XCTest
@testable import StudentApp

final class StudentAPIServiceTests: XCTestCase {
    var apiService: MockStudentAPIService!

    override func setUp() {
        super.setUp()
        apiService = MockStudentAPIService()
    }

    override func tearDown() {
        apiService = nil
        super.tearDown()
    }

    func testLoginSuccess() async throws {
        let token = try await apiService.login(email: "mario@scuola.it", password: "password123")
        XCTAssertEqual(token, "jwt_student_token_mock")
    }

    func testLoginFailure() async {
        do {
            _ = try await apiService.login(email: "mario@scuola.it", password: "123")
            XCTFail("Login should have failed for short password")
        } catch {
            XCTAssertNotNil(error)
        }
    }

    func testFetchGrades() async throws {
        let grades = try await apiService.fetchGrades(token: "valid_token")
        XCTAssertEqual(grades.count, 2)
    }
}
