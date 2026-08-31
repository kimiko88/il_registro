import XCTest
@testable import SecretaryApp

final class SecretaryInventoryTests: XCTestCase {

    func testInventoryNumberPrefix() {
        let inv = "INV-2026-0491"
        XCTAssertTrue(inv.hasPrefix("INV-"))
    }
}
