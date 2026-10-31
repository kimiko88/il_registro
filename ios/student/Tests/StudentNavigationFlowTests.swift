import XCTest
@testable import StudentApp

final class StudentNavigationFlowTests: XCTestCase {
    var viewModel: StudentViewModel!

    override func setUp() {
        super.setUp()
        viewModel = StudentViewModel()
        viewModel.loadSampleData()
    }

    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }

    func testStudentDashboardFullNavigationFlow() {
        // --- 1. Tab 0: Home Overview ---
        var selectedTab = 0
        XCTAssertEqual(selectedTab, 0)
        XCTAssertFalse(viewModel.grades.isEmpty)
        let overallGPA = viewModel.calculateGPA()
        XCTAssertEqual(overallGPA, 7.9, accuracy: 0.1)

        // --- 2. Tab 1: Grades Breakdown ---
        selectedTab = 1
        XCTAssertEqual(selectedTab, 1)
        XCTAssertEqual(viewModel.grades.count, 5)

        // Filter and inspect Math grades
        let mathGrades = viewModel.grades.filter { $0.subject == "Matematica" }
        XCTAssertEqual(mathGrades.count, 2)
        let mathSum = mathGrades.reduce(0.0) { $0 + ($1.grade * $1.weight) }
        let mathWeight = mathGrades.reduce(0.0) { $0 + $1.weight }
        let mathAvg = (mathSum / mathWeight * 10).rounded() / 10
        XCTAssertEqual(mathAvg, 7.8, accuracy: 0.1)

        // Check grade type localization mapping
        let writtenGrade = mathGrades.first(where: { $0.type == "Scritto" })
        XCTAssertNotNil(writtenGrade)
        XCTAssertFalse(writtenGrade!.localizedType.isEmpty)

        // --- 3. Tab 2: Agenda & Homework Task Tracking ---
        selectedTab = 2
        XCTAssertEqual(selectedTab, 2)
        XCTAssertEqual(viewModel.homework.count, 3)

        // Find open task
        let task1 = viewModel.homework.first(where: { $0.id == "1" })
        XCTAssertNotNil(task1)
        XCTAssertFalse(task1!.isCompleted)

        // Complete task via toggle
        let isCompletedNow = viewModel.toggleHomework(id: "1")
        XCTAssertTrue(isCompletedNow)
        XCTAssertTrue(viewModel.homework.first(where: { $0.id == "1" })?.isCompleted ?? false)

        // Toggle back to open state
        let isReopened = viewModel.toggleHomework(id: "1")
        XCTAssertFalse(isReopened)
        XCTAssertFalse(viewModel.homework.first(where: { $0.id == "1" })?.isCompleted ?? true)

        // --- 4. Tab 3: Attendance Records ---
        selectedTab = 3
        XCTAssertEqual(selectedTab, 3)
        XCTAssertEqual(viewModel.attendance.count, 2)

        let absences = viewModel.attendance.filter { $0.type == "Assenza" }
        XCTAssertEqual(absences.count, 1)
        XCTAssertTrue(absences.first?.isJustified ?? false)

        let lates = viewModel.attendance.filter { $0.type == "Ritardo" }
        XCTAssertEqual(lates.count, 1)

        // --- 5. Tab 4: Final Report Card ---
        selectedTab = 4
        XCTAssertEqual(selectedTab, 4)
        XCTAssertGreaterThan(viewModel.calculateGPA(), 6.0)
    }

    func testNavigationTabStateContinuity() {
        var activeTab = 0
        XCTAssertEqual(activeTab, 0)

        // Move to agenda and toggle homework
        activeTab = 2
        XCTAssertEqual(activeTab, 2)
        _ = viewModel.toggleHomework(id: "2")
        XCTAssertTrue(viewModel.homework.first(where: { $0.id == "2" })?.isCompleted ?? false)

        // Move to grades
        activeTab = 1
        XCTAssertEqual(activeTab, 1)
        XCTAssertEqual(viewModel.grades.count, 5)

        // Return to agenda and verify persistent state
        activeTab = 2
        XCTAssertEqual(activeTab, 2)
        XCTAssertTrue(viewModel.homework.first(where: { $0.id == "2" })?.isCompleted ?? false)

        // Return to home
        activeTab = 0
        XCTAssertEqual(activeTab, 0)
        XCTAssertEqual(viewModel.calculateGPA(), 7.9, accuracy: 0.1)
    }
}

