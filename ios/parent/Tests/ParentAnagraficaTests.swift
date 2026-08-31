import XCTest
@testable import ParentApp

final class ParentAnagraficaTests: XCTestCase {

    func testEmergencyPhoneFormat() {
        let phone = "+39 333 1234567"
        XCTAssertTrue(phone.hasPrefix("+39"))
    }
}
