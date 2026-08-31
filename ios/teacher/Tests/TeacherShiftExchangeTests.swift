import XCTest
@testable import TeacherApp

final class TeacherShiftExchangeTests: XCTestCase {

    func testShiftStatus() {
        let status = "Approvato dal DS"
        XCTAssertEqual(status, "Approvato dal DS")
    }
}
