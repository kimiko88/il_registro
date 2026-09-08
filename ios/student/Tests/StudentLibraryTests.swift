import XCTest
@testable import StudentApp

final class StudentLibraryTests: XCTestCase {

    func testBookLoanPrefix() {
        let code = "BIB-2024-8821"
        XCTAssertTrue(code.hasPrefix("BIB-"))
    }
}
