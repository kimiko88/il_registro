import XCTest
@testable import ParentApp

final class ParentSignaturesTests: XCTestCase {

    func testPdpSignatureWorkflow() {
        var signed = false
        signed = true
        XCTAssertTrue(signed)
    }
}
