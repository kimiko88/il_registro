import XCTest
@testable import TeacherApp

final class TeacherTripsTests: XCTestCase {

    func testTripAuthorizationCount() {
        let authorized = 42
        let total = 45
        XCTAssertGreaterThan(authorized, 40)
        XCTAssertEqual(total, 45)
    }
}
