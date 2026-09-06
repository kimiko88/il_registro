import XCTest
@testable import StudentApp

class MockWebSocketDelegate: WebSocketDelegate {
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

    func testConnectAndReceive() {
        guard let url = URL(string: "wss://registro-backend-fdu2.onrender.com/api/v1/ws") else { return }
        manager.connect(url: url, token: "test_token")
        XCTAssertTrue(manager.isConnected)
        XCTAssertTrue(delegate.connected)

        manager.simulateMessage("LIVE_UPDATE")
        XCTAssertEqual(delegate.lastMessage, "LIVE_UPDATE")

        manager.disconnect()
        XCTAssertFalse(manager.isConnected)
        XCTAssertTrue(delegate.disconnected)
    }
}
