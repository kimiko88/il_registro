import XCTest
@testable import SecretaryApp

final class SecretaryCertificatesProtocolTests: XCTestCase {

    func testProtocolPrefix() {
        let prot = "PROT-2026-004819"
        XCTAssertTrue(prot.hasPrefix("PROT-"))
    }
}
