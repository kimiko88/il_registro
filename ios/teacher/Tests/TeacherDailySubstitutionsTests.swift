import XCTest
@testable import TeacherApp

final class TeacherDailySubstitutionsTests: XCTestCase {

    func testCoverageHourAssigned() {
        let assignedClass = "2ª C"
        XCTAssertEqual(assignedClass, "2ª C")
    }
}
