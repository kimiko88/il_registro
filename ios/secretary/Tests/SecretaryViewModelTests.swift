import XCTest
@testable import SecretaryApp

final class SecretaryViewModelTests: XCTestCase {
    var viewModel: SecretaryViewModel!

    override func setUp() {
        super.setUp()
        viewModel = SecretaryViewModel()
    }

    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }

    func testAddUser() {
        let emptyFail = viewModel.addUser(name: "", role: "Docente")
        XCTAssertFalse(emptyFail)

        let success = viewModel.addUser(name: "Prof. Alberto Neri", role: "Docente")
        XCTAssertTrue(success)
        XCTAssertEqual(viewModel.users.count, 5)
    }

    func testIssueCertificate() {
        let cert = viewModel.issueCertificate(title: "Certificato di Diploma")
        XCTAssertEqual(cert.status, "pronto")
        XCTAssertTrue(cert.pdfUrl.contains(".pdf"))
        XCTAssertEqual(viewModel.certificates.count, 3)
    }
}
