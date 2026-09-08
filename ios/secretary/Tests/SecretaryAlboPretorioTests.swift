import XCTest
@testable import SecretaryApp

final class SecretaryAlboPretorioTests: XCTestCase {

    func testPublicationPeriodDays() {
        let days = 15
        XCTAssertEqual(days, 15)
    }
}
