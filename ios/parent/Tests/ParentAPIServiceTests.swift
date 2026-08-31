import XCTest
@testable import ParentApp

final class ParentAPIServiceTests: XCTestCase {
    var apiService: MockParentAPIService!

    override func setUp() {
        super.setUp()
        apiService = MockParentAPIService()
    }

    override func tearDown() {
        apiService = nil
        super.tearDown()
    }

    func testLoginSuccess() async throws {
        let token = try await apiService.login(email: "genitore@famiglia.it", password: "password123")
        XCTAssertEqual(token, "jwt_parent_token_mock")
    }

    func testFetchChildren() async throws {
        let children = try await apiService.fetchChildren(token: "valid_token")
        XCTAssertEqual(children.count, 2)
    }

    func testJustifyAbsenceValidation() async throws {
        let success = try await apiService.justifyAbsence(token: "token", absenceId: "a1", note: "Salute")
        XCTAssertTrue(success)
    }
}
