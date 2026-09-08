import XCTest
@testable import SecretaryApp

final class SecretaryCommunicationsTests: XCTestCase {

    func testCircularReadPercentage() {
        let read = 420.0
        let total = 472.0
        let percentage = (read / total) * 100.0
        XCTAssertGreaterThan(percentage, 80.0)
    }
}
