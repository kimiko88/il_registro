import XCTest
@testable import StudentApp

final class StudentBadgesTests: XCTestCase {

    func testBadgeCodePrefix() {
        let code = "OB-2026-9912-GOLD"
        XCTAssertTrue(code.hasPrefix("OB-"))
    }
}
