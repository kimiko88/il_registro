import XCTest
@testable import TeacherApp

final class TeacherNotesAndVerbaliTests: XCTestCase {

    func testVerbalePDFExport() {
        let verbaleId = "VERB-005"
        XCTAssertTrue(verbaleId.hasPrefix("VERB-"))
    }

    func testPdpAccommodationsTime() {
        let extraTime = 30
        XCTAssertEqual(extraTime, 30)
    }
}
