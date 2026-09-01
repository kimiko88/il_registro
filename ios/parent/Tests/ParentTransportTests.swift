import XCTest
@testable import ParentApp

final class ParentTransportTests: XCTestCase {

    func testPickupStop() {
        let stop = "Via Roma, 45"
        XCTAssertTrue(stop.contains("Roma"))
    }
}
