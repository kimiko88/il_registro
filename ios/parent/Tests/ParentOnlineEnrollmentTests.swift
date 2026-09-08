import XCTest
@testable import ParentApp

final class ParentOnlineEnrollmentTests: XCTestCase {

    func testEnrollmentAppPrefix() {
        let app = "ISC-2026-4019"
        XCTAssertTrue(app.hasPrefix("ISC-"))
    }
}
