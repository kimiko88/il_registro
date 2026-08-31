import XCTest
@testable import SecretaryApp

final class SecretaryWhistleblowingTests: XCTestCase {

    func testComplianceLaw() {
        let law = "D.Lgs. 24/2023"
        XCTAssertTrue(law.contains("24/2023"))
    }
}
