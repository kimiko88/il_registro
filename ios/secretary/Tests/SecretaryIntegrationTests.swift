import XCTest
@testable import SecretaryApp

final class SecretaryIntegrationTests: XCTestCase {

    var apiService: HttpSecretaryAPIService!
    var viewModel: SecretaryViewModel!

    override func setUp() {
        super.setUp()
        apiService = HttpSecretaryAPIService()
        viewModel = SecretaryViewModel()
    }

    func testSecretaryFullFlow() {
        XCTAssertNotNil(apiService)
        XCTAssertNotNil(viewModel)

        viewModel.loadSampleData()
        XCTAssertEqual(viewModel.users.count, 2)
        XCTAssertEqual(viewModel.classes.count, 2)

        let added = viewModel.addUser(firstName: "Laura", lastName: "Verdi", email: "laura.verdi@scuola.it", role: "teacher")
        XCTAssertTrue(added)
        XCTAssertEqual(viewModel.users.count, 3)

        let cert = viewModel.requestCertificate(studentId: "st1", type: "frequenza")
        XCTAssertNotNil(cert)
    }
}
