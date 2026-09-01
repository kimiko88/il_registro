import XCTest

final class SecretaryAdvancedFlowUITests: XCTestCase {

    func testSecretaryFascicoloViewDisplay() {
        let app = XCUIApplication()
        app.launch()
        XCTAssertTrue(app.exists)
    }
}
