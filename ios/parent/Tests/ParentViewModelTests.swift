import XCTest
@testable import ParentApp

final class ParentViewModelTests: XCTestCase {
    var viewModel: ParentViewModel!

    override func setUp() {
        super.setUp()
        viewModel = ParentViewModel()
    }

    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }

    func testSelectChildFiltering() {
        viewModel.selectChild(id: "c1")
        let c1Absences = viewModel.getAbsencesForSelectedChild()
        XCTAssertEqual(c1Absences.count, 2)

        viewModel.selectChild(id: "c2")
        let c2Absences = viewModel.getAbsencesForSelectedChild()
        XCTAssertEqual(c2Absences.count, 1)
    }

    func testJustifyAbsence() {
        let success = viewModel.justifyAbsence(id: "a1", note: "Visita dentistica")
        XCTAssertTrue(success)

        let absence = viewModel.absences.first(where: { $0.id == "a1" })
        XCTAssertTrue(absence?.isJustified ?? false)
        XCTAssertEqual(absence?.justificationNote, "Visita dentistica")
    }
}
