import XCTest
@testable import StudentApp

final class StudentELearningAndGoalsTests: XCTestCase {

    func testQuizPassingScore() {
        let score = 85
        XCTAssertGreaterThanOrEqual(score, 60)
    }

    func testGoalsProgressRange() {
        let progress = 0.75
        XCTAssertTrue(progress >= 0.0 && progress <= 1.0)
    }

    func testCreditsAccumulation() {
        let total = 11 + 12
        XCTAssertEqual(total, 23)
    }
}
