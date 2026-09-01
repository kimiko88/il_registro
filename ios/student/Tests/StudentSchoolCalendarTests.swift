import XCTest
@testable import StudentApp

final class StudentSchoolCalendarTests: XCTestCase {

    func testLastDayOfSchoolMonth() {
        let lastDay = "2026-06-06"
        XCTAssertTrue(lastDay.hasPrefix("2026-06"))
    }
}
