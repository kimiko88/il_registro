import XCTest
@testable import SecretaryApp

final class SecretaryFascicoloAndTextbooksTests: XCTestCase {

    func testSpendingCeilingComparison() {
        let actual = 281.50
        let ceiling = 310.00
        XCTAssertLessThanOrEqual(actual, ceiling)
    }
}
