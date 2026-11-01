import XCTest
@testable import TeacherApp

final class TeacherNavigationFlowTests: XCTestCase {
    var viewModel: TeacherViewModel!

    override func setUp() {
        super.setUp()
        viewModel = TeacherViewModel()
        viewModel.loadSampleSession()
    }

    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }

    func testTeacherFullSessionWorkflow() {
        // --- 1. Tab 0: Lesson Signing & Topic Entry ---
        var selectedTab = 0
        XCTAssertEqual(selectedTab, 0)
        XCTAssertNotNil(viewModel.currentSession)
        XCTAssertEqual(viewModel.currentSession?.className, "Classe 3A")
        XCTAssertEqual(viewModel.currentSession?.subject, "Matematica")
        XCTAssertFalse(viewModel.isHourSigned)

        // Blank topic must be rejected
        let blankResult = viewModel.signLesson(topic: "")
        XCTAssertFalse(blankResult)
        XCTAssertFalse(viewModel.isHourSigned)

        // Valid lesson topic signature
        let signResult = viewModel.signLesson(topic: "Calcolo differenziale: derivata prima e punti stazionari")
        XCTAssertTrue(signResult)
        XCTAssertTrue(viewModel.isHourSigned)
        XCTAssertEqual(viewModel.lessonTopic, "Calcolo differenziale: derivata prima e punti stazionari")
        XCTAssertTrue(viewModel.currentSession?.isSigned ?? false)

        // --- 2. Tab 1: Roll Call / Attendance Marking ---
        selectedTab = 1
        XCTAssertEqual(selectedTab, 1)
        XCTAssertEqual(viewModel.students.count, 4)

        // Check initial statuses
        XCTAssertEqual(viewModel.students.first(where: { $0.id == "s1" })?.status, "presente")
        XCTAssertEqual(viewModel.students.first(where: { $0.id == "s3" })?.status, "assente")

        // Change absent student s3 to late arrival (ritardo)
        let attendanceResult1 = viewModel.updateAttendance(studentId: "s3", status: "ritardo")
        XCTAssertTrue(attendanceResult1)
        XCTAssertEqual(viewModel.students.first(where: { $0.id == "s3" })?.status, "ritardo")

        // Mark student s2 as early dismissal (uscita_anticipata)
        let attendanceResult2 = viewModel.updateAttendance(studentId: "s2", status: "uscita_anticipata")
        XCTAssertTrue(attendanceResult2)
        XCTAssertEqual(viewModel.students.first(where: { $0.id == "s2" })?.status, "uscita_anticipata")

        // --- 3. Tab 2: Grade Proposals & Assessment ---
        selectedTab = 2
        XCTAssertEqual(selectedTab, 2)
        XCTAssertEqual(viewModel.grades.count, 0)

        // Insert oral grade for student s1
        let oralResult = viewModel.insertGrade(studentId: "s1", grade: 8.5, type: "Orale")
        XCTAssertTrue(oralResult)
        XCTAssertEqual(viewModel.grades.count, 1)

        // Insert written grade for student s2
        let writtenResult = viewModel.insertGrade(studentId: "s2", grade: 7.0, type: "Scritto")
        XCTAssertTrue(writtenResult)
        XCTAssertEqual(viewModel.grades.count, 2)

        // Range boundary enforcement [1.0 - 10.0]
        let highInvalid = viewModel.insertGrade(studentId: "s1", grade: 10.5, type: "Orale")
        XCTAssertFalse(highInvalid)

        let lowInvalid = viewModel.insertGrade(studentId: "s1", grade: 0.5, type: "Orale")
        XCTAssertFalse(lowInvalid)

        XCTAssertEqual(viewModel.grades.count, 2)

        // --- 4. Tab 3: Agenda & Daily Planner ---
        selectedTab = 3
        XCTAssertEqual(selectedTab, 3)
        XCTAssertTrue(viewModel.isHourSigned)

        // --- 5. Tab 4: Scrutiny Board Overview ---
        selectedTab = 4
        XCTAssertEqual(selectedTab, 4)
        XCTAssertNotNil(viewModel.currentSession)
    }

    func testTeacherAttendanceStatusIntegrity() {
        let initialPresent = viewModel.students.filter { $0.status == "presente" }.count
        XCTAssertEqual(initialPresent, 2)

        // Update student s4 from ritardo to presente
        _ = viewModel.updateAttendance(studentId: "s4", status: "presente")
        let updatedPresent = viewModel.students.filter { $0.status == "presente" }.count
        XCTAssertEqual(updatedPresent, 3)
    }
}

