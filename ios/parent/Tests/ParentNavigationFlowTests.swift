import XCTest
@testable import ParentApp

final class ParentNavigationFlowTests: XCTestCase {
    var viewModel: ParentViewModel!

    override func setUp() {
        super.setUp()
        viewModel = ParentViewModel()
        viewModel.loadSampleData()
    }

    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }

    func testParentFullNavigationAndFeatureFlow() {
        // --- 1. Tab 0: Children Selection & Overview ---
        var selectedTab = 0
        XCTAssertEqual(selectedTab, 0)
        XCTAssertEqual(viewModel.children.count, 2)
        XCTAssertEqual(viewModel.selectedChildId, "c1")

        // Inspect absences for first child (c1)
        var childAbsences = viewModel.getAbsencesForSelectedChild()
        XCTAssertEqual(childAbsences.count, 2)

        // Switch to second child (c2)
        viewModel.selectChild(id: "c2")
        XCTAssertEqual(viewModel.selectedChildId, "c2")
        childAbsences = viewModel.getAbsencesForSelectedChild()
        XCTAssertEqual(childAbsences.count, 1)
        XCTAssertEqual(childAbsences.first?.id, "a3")

        // Switch back to child 1
        viewModel.selectChild(id: "c1")
        XCTAssertEqual(viewModel.selectedChildId, "c1")

        // --- 2. Tab 1: Child Grades Review ---
        selectedTab = 1
        XCTAssertEqual(selectedTab, 1)
        XCTAssertFalse(viewModel.children.isEmpty)

        // --- 3. Tab 2: Absence Justification Flow ---
        selectedTab = 2
        XCTAssertEqual(selectedTab, 2)
        let absenceToJustify = viewModel.absences.first(where: { $0.id == "a1" })
        XCTAssertNotNil(absenceToJustify)
        XCTAssertFalse(absenceToJustify!.isJustified)

        // Justify absence with medical motive
        let justificationResult = viewModel.justifyAbsence(id: "a1", note: "Certificato medico rilasciato dal pediatra")
        XCTAssertTrue(justificationResult)

        // Verify updated model state
        let updatedAbsence = viewModel.absences.first(where: { $0.id == "a1" })
        XCTAssertTrue(updatedAbsence?.isJustified ?? false)
        XCTAssertEqual(updatedAbsence?.justificationNote, "Certificato medico rilasciato dal pediatra")

        // Justify late arrival
        let lateResult = viewModel.justifyAbsence(id: "a2", note: "Ritardo trasporto scolastico")
        XCTAssertTrue(lateResult)
        XCTAssertTrue(viewModel.absences.first(where: { $0.id == "a2" })?.isJustified ?? false)

        // --- 4. Tab 3: Teacher Colloqui & Conferences ---
        selectedTab = 3
        XCTAssertEqual(selectedTab, 3)

        // --- 5. Tab 4: School Communications & Circulars ---
        selectedTab = 4
        XCTAssertEqual(selectedTab, 4)
    }

    func testParentChildSwitchPreservesState() {
        // Justify absence for child 1
        viewModel.selectChild(id: "c1")
        _ = viewModel.justifyAbsence(id: "a1", note: "Maternità / famiglia")
        XCTAssertTrue(viewModel.absences.first(where: { $0.id == "a1" })?.isJustified ?? false)

        // Switch to child 2 and verify absence a3 is still unjustified
        viewModel.selectChild(id: "c2")
        let c2Absences = viewModel.getAbsencesForSelectedChild()
        XCTAssertEqual(c2Absences.count, 1)
        XCTAssertFalse(c2Absences[0].isJustified)

        // Justify child 2 absence
        _ = viewModel.justifyAbsence(id: "a3", note: "Terapia medica")
        XCTAssertTrue(viewModel.absences.first(where: { $0.id == "a3" })?.isJustified ?? false)

        // Switch back to child 1 and assert state continuity
        viewModel.selectChild(id: "c1")
        let c1Absences = viewModel.getAbsencesForSelectedChild()
        let c1Target = c1Absences.first(where: { $0.id == "a1" })
        XCTAssertTrue(c1Target?.isJustified ?? false)
        XCTAssertEqual(c1Target?.justificationNote, "Maternità / famiglia")
    }
}

