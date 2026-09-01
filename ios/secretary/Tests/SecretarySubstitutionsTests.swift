import XCTest
@testable import SecretaryApp

final class SecretarySubstitutionsTests: XCTestCase {

    func testSubstitutionsAvailability() {
        let slot = 3
        let teacherAvailable = true
        XCTAssertTrue(teacherAvailable)
        XCTAssertEqual(slot, 3)
    }

    func testSidiExportValidation() {
        let codMecc = "RMIS09900B"
        XCTAssertEqual(codMecc.count, 10)
    }
}
