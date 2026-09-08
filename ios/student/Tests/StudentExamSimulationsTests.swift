import XCTest
@testable import StudentApp

final class StudentExamSimulationsTests: XCTestCase {

    func testMaxPoints() {
        let maxPoints = 20
        XCTAssertEqual(maxPoints, 20)
    }
}
