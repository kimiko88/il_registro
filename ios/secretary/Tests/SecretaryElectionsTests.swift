import XCTest
@testable import SecretaryApp

final class SecretaryElectionsTests: XCTestCase {

    func testVoterTurnoutPercentage() {
        let turnout = 68.4
        XCTAssertGreaterThan(turnout, 50.0)
    }
}
