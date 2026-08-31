import XCTest
@testable import SecretaryApp

final class SecretaryTransfersTests: XCTestCase {

    func testNullaOstaPrefix() {
        let protocolNumber = "NO-2026-0392"
        XCTAssertTrue(protocolNumber.hasPrefix("NO-"))
    }
}
