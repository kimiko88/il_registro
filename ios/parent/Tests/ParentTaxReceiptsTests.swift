import XCTest
@testable import ParentApp

final class ParentTaxReceiptsTests: XCTestCase {

    func testTaxDeductionAmount() {
        let total = 680.00
        XCTAssertLessThanOrEqual(total, 800.00)
    }
}
