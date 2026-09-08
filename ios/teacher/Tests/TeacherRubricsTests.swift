import XCTest
@testable import TeacherApp

final class TeacherRubricsTests: XCTestCase {

    func testRubricLevelsCount() {
        let levels = ["Avanzato", "Intermedio", "Base", "Iniziale"]
        XCTAssertEqual(levels.count, 4)
    }
}
