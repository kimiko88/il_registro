import XCTest
@testable import TeacherApp

final class TeacherLeavesTests: XCTestCase {

    func testPermitHours() {
        let hours = 2
        XCTAssertEqual(hours, 2)
    }
}
