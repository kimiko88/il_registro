import XCTest
@testable import StudentApp

final class MockWebSocketDelegate: WebSocketDelegate {
    var connected = false
    var lastMessage = ""
    var disconnected = false

    func webSocketDidConnect() {
        connected = true
    }

    func webSocketDidReceiveMessage(_ message: String) {
        lastMessage = message
    }

    func webSocketDidDisconnect() {
        disconnected = true
    }
}

final class WebSocketManagerTests: XCTestCase {
    var manager: WebSocketManager!
    var delegate: MockWebSocketDelegate!

    override func setUp() {
        super.setUp()
        manager = WebSocketManager()
        delegate = MockWebSocketDelegate()
        manager.delegate = delegate
    }

    override func tearDown() {
        manager = nil
        delegate = nil
        super.tearDown()
    }

    func testPlaceholderNeverClaimsConnectionOrDeliversSyntheticMessages() {
        let url = URL(string: "wss://registro-backend-fdu2.onrender.com/api/v1/ws")!
        manager.connect(url: url, token: "test_token")

        XCTAssertFalse(manager.isConnected)
        XCTAssertFalse(delegate.connected)
        XCTAssertTrue(delegate.disconnected)

        manager.simulateMessage("LIVE_UPDATE")
        XCTAssertEqual(delegate.lastMessage, "")
    }
}
