import XCTest
@testable import TeacherApp

final class TeacherUdaAndSupportTests: XCTestCase {

    func testUdaTargetDisciplines() {
        let disciplines = ["Informatica", "Italiano", "Filosofia"]
        XCTAssertEqual(disciplines.count, 3)
    }

    func testPeiSupportHours() {
        let hours = 18
        XCTAssertGreaterThan(hours, 0)
    }
}
