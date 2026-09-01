import XCTest
@testable import StudentApp

final class OfflineCacheManagerTests: XCTestCase {
    var cache: OfflineCacheManager!

    override func setUp() {
        super.setUp()
        cache = OfflineCacheManager()
    }

    override func tearDown() {
        cache = nil
        super.tearDown()
    }

    func testSaveAndRetrieve() {
        let grades = [
            GradeItemModel(id: "1", subject: "Matematica", grade: 8.5, weight: 1.0, type: "Scritto", date: "30 Ago")
        ]
        cache.saveGrades(grades)
        XCTAssertEqual(cache.getCachedGrades().count, 1)
        XCTAssertTrue(cache.isCacheValid())
    }

    func testClearCache() {
        cache.saveGrades([GradeItemModel(id: "1", subject: "Fisica", grade: 7.0, weight: 1.0, type: "Scritto", date: "22 Ago")])
        cache.clear()
        XCTAssertTrue(cache.getCachedGrades().isEmpty)
        XCTAssertFalse(cache.isCacheValid())
    }
}
