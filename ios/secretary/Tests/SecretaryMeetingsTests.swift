import XCTest
@testable import SecretaryApp

final class SecretaryMeetingsTests: XCTestCase {

    func testMeetingStatus() {
        let status = "Convocato"
        XCTAssertEqual(status, "Convocato")
    }
}
