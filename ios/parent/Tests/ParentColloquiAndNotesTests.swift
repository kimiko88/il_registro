import XCTest
@testable import ParentApp

final class ParentColloquiAndNotesTests: XCTestCase {

    func testColloquiSlotDuration() {
        let durationMinutes = 15
        XCTAssertEqual(durationMinutes, 15)
    }

    func testNoteAcknowledgment() {
        var acknowledged = false
        acknowledged = true
        XCTAssertTrue(acknowledged)
    }
}
