import XCTest
@testable import StudentApp

final class StudentGradeSimulatorTests: XCTestCase {

    func testSimulationResultCalculation() {
        let grades = [7.0, 7.0, 7.5, 8.0]
        let avg = grades.reduce(0, +) / Double(grades.count)
        XCTAssertEqual(avg, 7.375, accuracy: 0.001)
    }
}
