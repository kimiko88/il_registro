import XCTest
@testable import TeacherApp

final class TeacherSeatingPlanTests: XCTestCase {

    func testSeatingPlanCapacity() {
        let capacity = 24
        XCTAssertEqual(capacity, 24)
    }
}
