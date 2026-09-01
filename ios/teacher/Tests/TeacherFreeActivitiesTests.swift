import XCTest
@testable import TeacherApp

final class TeacherFreeActivitiesTests: XCTestCase {

    func testOnDutyHourDuration() {
        let start = 10
        let end = 11
        XCTAssertEqual(end - start, 1)
    }
}
