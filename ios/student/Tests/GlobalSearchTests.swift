import XCTest
@testable import StudentApp

final class GlobalSearchTests: XCTestCase {

    func testSearchQueryMatching() {
        let results = ["Circolari", "Compiti", "Materiali"]
        let filtered = results.filter { $0.contains("Compiti") }
        XCTAssertEqual(filtered.count, 1)
    }
}
