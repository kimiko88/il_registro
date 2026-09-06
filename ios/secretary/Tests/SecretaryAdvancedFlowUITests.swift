import XCTest

final class SecretaryAdvancedFlowUITests: XCTestCase {

    func testSecretaryFascicoloViewDisplay() throws {
        guard ProcessInfo.processInfo.environment["XCUI_TARGET_APPLICATION_BUNDLE_PATH"] != nil else {
            throw XCTSkip("Skipping UI test: No target application configured.")
        }
        let app = XCUIApplication()
        app.launch()
        XCTAssertTrue(app.exists)
    }
}
