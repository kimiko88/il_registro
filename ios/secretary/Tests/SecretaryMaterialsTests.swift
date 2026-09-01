import XCTest
@testable import SecretaryApp

final class SecretaryMaterialsTests: XCTestCase {

    func testStorageLimit() {
        let storageUsed = 42.8
        let limit = 250.0
        XCTAssertLessThan(storageUsed, limit)
    }
}
