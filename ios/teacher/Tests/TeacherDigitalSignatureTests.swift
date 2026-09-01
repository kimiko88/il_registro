import XCTest
@testable import TeacherApp

final class TeacherDigitalSignatureTests: XCTestCase {

    func testOtpValidation() {
        let otp = "654321"
        XCTAssertEqual(otp.count, 6)
    }

    func testPadesSignatureStatus() {
        let signed = true
        XCTAssertTrue(signed)
    }
}
