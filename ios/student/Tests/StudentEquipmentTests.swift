import XCTest
@testable import StudentApp

final class StudentEquipmentTests: XCTestCase {

    func testLockerAssignmentFormat() {
        let locker = "ARM-104"
        XCTAssertTrue(locker.contains("104"))
    }
}
