import XCTest
@testable import StudentApp

final class MdmConfigManagerTests: XCTestCase {

    func testMdmManagerSharedInstance() {
        let manager = MdmConfigManager.shared
        XCTAssertNotNil(manager)
        // Should not crash when no managed config exists
        manager.applyManagedConfiguration()
    }
}
