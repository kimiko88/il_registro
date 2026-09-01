import XCTest
@testable import TeacherApp

final class TeacherGlobalReportsTests: XCTestCase {

    func testJudgmentStatus() {
        let status = "Deliberato dal CdC"
        XCTAssertEqual(status, "Deliberato dal CdC")
    }
}
