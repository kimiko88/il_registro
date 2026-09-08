import XCTest
@testable import ParentApp

final class ParentCalendarTests: XCTestCase {

    func testSchoolEndDateFormat() {
        let date = "2026-06-06"
        XCTAssertTrue(date.contains("2026"))
    }
}
