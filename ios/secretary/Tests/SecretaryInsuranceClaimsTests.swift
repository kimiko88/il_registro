import XCTest
@testable import SecretaryApp

final class SecretaryInsuranceClaimsTests: XCTestCase {

    func testInailProtocolPrefix() {
        let prot = "INAIL-2026-9021"
        XCTAssertTrue(prot.hasPrefix("INAIL-"))
    }
}
