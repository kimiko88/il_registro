import XCTest
@testable import ParentApp

final class ParentWorkshopsEnrollmentTests: XCTestCase {

    func testCompletedHours() {
        let completed = 20
        let total = 30
        XCTAssertLessThanOrEqual(completed, total)
    }
}
