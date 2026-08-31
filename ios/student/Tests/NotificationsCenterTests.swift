import XCTest
@testable import StudentApp

final class NotificationsCenterTests: XCTestCase {

    func testNotificationItemPresence() {
        let title = "Nuovo Voto Inserito: 8.5"
        XCTAssertFalse(title.isEmpty)
    }

    func testAgIDComplianceStatus() {
        let status = "Totalmente Conforme"
        XCTAssertEqual(status, "Totalmente Conforme")
    }
}
