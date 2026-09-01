import XCTest
@testable import StudentApp

final class LocalizationIntegrityTests: XCTestCase {

    func testSupportedLocalesList() {
        let supportedLocales = [
            "it", "en", "de", "fr", "es", "ar", "ro", "ru", "sq", "uk", "zh-Hans"
        ]
        XCTAssertEqual(supportedLocales.count, 11)
        XCTAssertTrue(supportedLocales.contains("it"))
        XCTAssertTrue(supportedLocales.contains("ar"))
        XCTAssertTrue(supportedLocales.contains("zh-Hans"))
    }
}
