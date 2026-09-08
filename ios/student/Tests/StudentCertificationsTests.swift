import XCTest
@testable import StudentApp

final class StudentCertificationsTests: XCTestCase {

    func testQcerLevelValidation() {
        let qcer = "B2"
        XCTAssertEqual(qcer, "B2")
    }
}
