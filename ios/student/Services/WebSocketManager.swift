import Foundation

public protocol WebSocketDelegate: AnyObject {
    func webSocketDidConnect()
    func webSocketDidReceiveMessage(_ message: String)
    func webSocketDidDisconnect()
}

public class WebSocketManager {
    public static let shared = WebSocketManager()
    public weak var delegate: WebSocketDelegate?
    public private(set) var isConnected: Bool = false

    public init() {}

    public func connect(url: URL, token: String) {
        guard !token.isEmpty else { return }
        isConnected = true
        delegate?.webSocketDidConnect()
    }

    public func simulateMessage(_ message: String) {
        guard isConnected else { return }
        delegate?.webSocketDidReceiveMessage(message)
    }

    public func disconnect() {
        isConnected = false
        delegate?.webSocketDidDisconnect()
    }
}
