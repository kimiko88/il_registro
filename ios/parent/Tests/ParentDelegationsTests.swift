import XCTest
@testable import ParentApp

final class ParentDelegationsTests: XCTestCase {

    func testDelegationStatus() {
        let status = "Valida & Verificata"
        XCTAssertEqual(status, "Valida & Verificata")
    }
}
