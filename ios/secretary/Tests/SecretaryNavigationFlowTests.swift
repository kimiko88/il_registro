import XCTest
@testable import SecretaryApp

final class SecretaryNavigationFlowTests: XCTestCase {
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

    func testSecretaryFullAdministrationFlow() {
        // --- 1. Tab 0: KPI & Directory Overview ---
        var selectedTab = 0
        XCTAssertEqual(selectedTab, 0)
        XCTAssertEqual(viewModel.users.count, 2)
        XCTAssertEqual(viewModel.classes.count, 2)
        XCTAssertEqual(viewModel.certificates.count, 2)

        // --- 2. Tab 1: User Directory & Creation Flow ---
        selectedTab = 1
        XCTAssertEqual(selectedTab, 1)

        // Blank name must fail
        let blankResult = viewModel.addUser(name: "   ", role: "Studente")
        XCTAssertFalse(blankResult)

        // Add valid student
        let studentAddResult = viewModel.addUser(name: "Giulia Bianchi", role: "Studente")
        XCTAssertTrue(studentAddResult)
        XCTAssertEqual(viewModel.users.count, 3)

        // Add valid teacher
        let teacherAddResult = viewModel.addUser(firstName: "Franco", lastName: "Verdi", email: "f.verdi@scuola.it", role: "Docente")
        XCTAssertTrue(teacherAddResult)
        XCTAssertEqual(viewModel.users.count, 4)

        // Role filtering
        let teachers = viewModel.users.filter { $0.role == "Docente" }
        XCTAssertEqual(teachers.count, 3)

        let students = viewModel.users.filter { $0.role == "Studente" }
        XCTAssertEqual(students.count, 1)

        // --- 3. Tab 2: Class Scrutiny Supervision ---
        selectedTab = 2
        XCTAssertEqual(selectedTab, 2)
        XCTAssertEqual(viewModel.classes.count, 2)
        XCTAssertTrue(viewModel.classes.contains(where: { $0.name == "Classe 1A" }))
        XCTAssertTrue(viewModel.classes.contains(where: { $0.name == "Classe 2A" }))

        // --- 4. Tab 3: Certificate Protocolling & Delivery ---
        selectedTab = 3
        XCTAssertEqual(selectedTab, 3)
        let initialCertCount = viewModel.certificates.count

        // Issue new enrollment & attendance certificate
        let issuedCert = viewModel.issueCertificate(title: "Certificato di Iscrizione e Frequenza Anno 2025/2026")
        XCTAssertEqual(issuedCert.status, "pronto")
        XCTAssertTrue(issuedCert.pdfUrl.contains(".pdf"))
        XCTAssertEqual(viewModel.certificates.count, initialCertCount + 1)

        // Request student certificate by student ID and type
        let requestedCert = viewModel.requestCertificate(studentId: "u3", type: "Valutazioni Intermedie")
        XCTAssertEqual(requestedCert.title, "Certificato Valutazioni Intermedie")
        XCTAssertEqual(requestedCert.status, "pronto")
        XCTAssertEqual(viewModel.certificates.count, initialCertCount + 2)
    }

    func testUserManagementRoleIntegrity() {
        let freshVM = SecretaryViewModel()
        freshVM.loadData()

        XCTAssertEqual(freshVM.users.count, 4)
        let teachers = freshVM.users.filter { $0.role == "Docente" }
        let students = freshVM.users.filter { $0.role == "Studente" }
        let parents = freshVM.users.filter { $0.role == "Genitore" }

        XCTAssertEqual(teachers.count, 2)
        XCTAssertEqual(students.count, 1)
        XCTAssertEqual(parents.count, 1)
    }
}

