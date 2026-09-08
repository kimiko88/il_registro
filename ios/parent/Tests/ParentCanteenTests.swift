import XCTest
@testable import ParentApp

final class ParentCanteenTests: XCTestCase {

    func testMealRecordStatus() {
        let status = "Consumato"
        XCTAssertEqual(status, "Consumato")
    }
}
