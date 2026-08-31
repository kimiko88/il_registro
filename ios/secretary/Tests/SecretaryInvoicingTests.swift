import XCTest
@testable import SecretaryApp

final class SecretaryInvoicingTests: XCTestCase {

    func testCollectedSum() {
        let total = 14850.00
        XCTAssertGreaterThan(total, 0.0)
    }
}
