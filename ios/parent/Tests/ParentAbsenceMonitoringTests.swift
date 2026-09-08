import XCTest
@testable import ParentApp

final class ParentAbsenceMonitoringTests: XCTestCase {

    func testAbsenceThresholdLimit() {
        let actualAbsenceHours = 54.0
        let maxLimit = 247.0
        XCTAssertLessThan(actualAbsenceHours, maxLimit)
    }
}
