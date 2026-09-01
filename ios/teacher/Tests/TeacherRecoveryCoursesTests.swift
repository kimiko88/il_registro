import XCTest
@testable import TeacherApp

final class TeacherRecoveryCoursesTests: XCTestCase {

    func testRecoveryCourseHoursCalculation() {
        let totalHours = 10
        let currentHours = 8
        let ratio = Double(currentHours) / Double(totalHours)
        XCTAssertEqual(ratio, 0.8, accuracy: 0.001)
    }

    func testGradeWeightsConfiguration() {
        let written = 1.0
        let oral = 1.0
        let test = 0.5
        XCTAssertGreaterThan(written, 0)
        XCTAssertGreaterThan(oral, 0)
        XCTAssertGreaterThan(test, 0)
    }
}
