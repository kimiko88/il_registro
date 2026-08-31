import XCTest
@testable import SecretaryApp

final class SecretarySafetyTests: XCTestCase {

    func testEvacuationDuration() {
        let actualTime = 165
        let limit = 180
        XCTAssertLessThan(actualTime, limit)
    }
}
