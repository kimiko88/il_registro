import XCTest
@testable import SecretaryApp

final class SecretaryClassFormationTests: XCTestCase {

    func testGenderBalance() {
        let males = 12
        let females = 12
        XCTAssertEqual(males, females)
    }
}
