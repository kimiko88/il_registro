import XCTest
@testable import SecretaryApp

final class SecretaryIntegrationTests: XCTestCase {
    var apiService: MockSecretaryAPIService!
    var viewModel: SecretaryViewModel!

    override func setUp() {
        super.setUp()
        apiService = MockSecretaryAPIService()
        viewModel = SecretaryViewModel()
    }

    override func tearDown() {
        apiService = nil
        viewModel = nil
        super.tearDown()
    }

    func testSecretaryIntegrationFlow() async throws {
        // 1. Authenticate Secretary
        let token = try await apiService.login(email: "segreteria@scuola.it", password: "password123")
        XCTAssertFalse(token.isEmpty)

        // 2. Fetch Users
        let users = try await apiService.fetchUsers(token: token)
        XCTAssertEqual(users.count, 2)

        viewModel.loadData()
        let addResult = viewModel.addUser(name: "Prof. Alberto Neri", role: "Docente")
        XCTAssertTrue(addResult)
        XCTAssertEqual(viewModel.users.count, 5)

        // 3. Issue Certificate
        let cert = try await apiService.issueCertificate(token: token, studentId: "s1", type: "Frequenza")
        XCTAssertEqual(cert.status, "pronto")
        XCTAssertTrue(cert.pdfUrl.contains(".pdf"))
    }
}
