import XCTest
@testable import ParentApp

final class ParentCyberbullyingReportTests: XCTestCase {

    func testPracticePrefix() {
        let code = "BULL-2026-081"
        XCTAssertTrue(code.hasPrefix("BULL-"))
    }
}
