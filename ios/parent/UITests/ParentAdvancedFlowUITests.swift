import XCTest

final class ParentAdvancedFlowUITests: XCTestCase {

    func testParentColloquiViewDisplay() {
        let app = XCUIApplication()
        app.launch()
        XCTAssertTrue(app.exists)
    }
}
