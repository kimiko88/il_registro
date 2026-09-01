import XCTest
@testable import StudentApp

final class AppConfigTests: XCTestCase {

    func testAppConfigBaseURL() {
        guard let url = URL(string: AppConfig.baseURL) else {
            XCTFail("Base URL is not a valid URL")
            return
        }
        XCTAssertTrue(url.scheme == "http" || url.scheme == "https")
        XCTAssertTrue(AppConfig.baseURL.hasSuffix("/api/v1"))
    }

    func testAppConfigWsURL() {
        guard let url = URL(string: AppConfig.wsURL) else {
            XCTFail("WS URL is not a valid URL")
            return
        }
        XCTAssertTrue(url.scheme == "ws" || url.scheme == "wss")
        XCTAssertTrue(AppConfig.wsURL.hasSuffix("/ws"))
    }

    func testTimeout() {
        XCTAssertGreaterThanOrEqual(AppConfig.timeoutInterval, 10.0)
    }
}
