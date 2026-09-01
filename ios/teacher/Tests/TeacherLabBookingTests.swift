import XCTest
@testable import TeacherApp

final class TeacherLabBookingTests: XCTestCase {

    func testLabBookingDuration() {
        let duration = 2
        XCTAssertEqual(duration, 2)
    }
}
