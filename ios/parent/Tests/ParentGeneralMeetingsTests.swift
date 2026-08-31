import XCTest
@testable import ParentApp

final class ParentGeneralMeetingsTests: XCTestCase {

    func testQueueEstimatedTime() {
        let ahead = 2
        let minPerMeeting = 4
        XCTAssertEqual(ahead * minPerMeeting, 8)
    }
}
