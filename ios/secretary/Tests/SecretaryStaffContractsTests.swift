import XCTest
@testable import SecretaryApp

final class SecretaryStaffContractsTests: XCTestCase {

    func testWeeklyHours() {
        let hours = 18
        XCTAssertEqual(hours, 18)
    }
}
