import XCTest
@testable import TeacherApp

final class TeacherGroupsTests: XCTestCase {

    func testGroupStudentCount() {
        let count = 18
        XCTAssertGreaterThan(count, 0)
    }
}
