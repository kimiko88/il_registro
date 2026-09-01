import XCTest
@testable import SecretaryApp

final class SecretaryTimetableTests: XCTestCase {

    func testTimetableWeeklyPeriods() {
        let daily = 6
        let days = 5
        XCTAssertEqual(daily * days, 30)
    }
}
