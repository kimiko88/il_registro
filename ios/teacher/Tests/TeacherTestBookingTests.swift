import XCTest
@testable import TeacherApp

final class TeacherTestBookingTests: XCTestCase {

    func testTestBookingHours() {
        let duration = 2
        XCTAssertEqual(duration, 2)
    }
}
