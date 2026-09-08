import XCTest
@testable import SecretaryApp

final class SecretaryViewModelTests: XCTestCase {
    var viewModel: SecretaryViewModel!

    override func setUp() {
        super.setUp()
        viewModel = SecretaryViewModel()
        viewModel.loadSampleData()
    }

    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }

    func testInitialStateIsEmpty() {
        let freshVM = SecretaryViewModel()
        XCTAssertTrue(freshVM.users.isEmpty)
        XCTAssertTrue(freshVM.classes.isEmpty)
        XCTAssertTrue(freshVM.certificates.isEmpty)
    }

    func testAddUser() {
        let emptyFail = viewModel.addUser(name: "", role: "Docente")
        XCTAssertFalse(emptyFail)

        let success = viewModel.addUser(name: "Prof. Alberto Neri", role: "Docente")
        XCTAssertTrue(success)
        XCTAssertEqual(viewModel.users.count, 3)
    }

    func testIssueCertificate() {
        let cert = viewModel.issueCertificate(title: "Certificato di Diploma")
        XCTAssertEqual(cert.status, "pronto")
        XCTAssertTrue(cert.pdfUrl.contains(".pdf"))
        XCTAssertEqual(viewModel.certificates.count, 3)
    }
}
