import XCTest
@testable import StudentApp

final class StudentIntegrationTests: XCTestCase {

    var apiService: HttpStudentAPIService!
    var viewModel: StudentViewModel!

    override func setUp() {
        super.setUp()
        apiService = HttpStudentAPIService()
        viewModel = StudentViewModel()
    }

    func testStudentFullFlow() {
        XCTAssertNotNil(apiService)
        XCTAssertNotNil(viewModel)

        viewModel.loadSampleData()
        let gpa = viewModel.calculateGPA()
        XCTAssertEqual(gpa, 7.9, accuracy: 0.01)

        let toggled = viewModel.toggleHomework(id: "1")
        XCTAssertTrue(toggled)
    }

    func testRealStudentLoginAndFetchGrades() async throws {
        let token = try await apiService.login(email: "studentea_1@scuola.it", password: "password")
        XCTAssertFalse(token.isEmpty)
        XCTAssertEqual(apiService.lastStudentName, "Studente2A_1 Test")

        let grades = try await apiService.fetchGrades(token: token)
        XCTAssertFalse(grades.isEmpty)
        XCTAssertTrue(grades.contains(where: { $0.subject == "Matematica" }))

        await viewModel.loadFromDatabase(token: token)
        XCTAssertFalse(viewModel.grades.isEmpty)
        XCTAssertGreaterThan(viewModel.calculateGPA(), 6.0)
    }
}
