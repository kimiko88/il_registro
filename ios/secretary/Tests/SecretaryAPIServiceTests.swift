import XCTest
@testable import SecretaryApp

final class SecretaryAPIServiceTests: XCTestCase {
    var apiService: MockSecretaryAPIService!

    override func setUp() {
        super.setUp()
        apiService = MockSecretaryAPIService()
    }

    override func tearDown() {
        apiService = nil
        super.tearDown()
    }

    func testLoginSuccess() async throws {
        let token = try await apiService.login(email: "segreteria@scuola.it", password: "password123")
        XCTAssertEqual(token, "jwt_secretary_token_mock")
    }

    func testFetchUsers() async throws {
        let users = try await apiService.fetchUsers(token: "valid_token")
        XCTAssertEqual(users.count, 2)
    }

    func testIssueCertificate() async throws {
        let cert = try await apiService.issueCertificate(token: "token", studentId: "s1", type: "Iscrizione")
        XCTAssertEqual(cert.status, "pronto")
        XCTAssertTrue(cert.pdfUrl.contains(".pdf"))
    }
}
