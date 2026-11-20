import Foundation

public protocol WebSocketDelegate: AnyObject {
    func webSocketDidConnect()
    func webSocketDidReceiveMessage(_ message: String)
    func webSocketDidDisconnect()
}

/**
 * Placeholder until a real authenticated WSS transport is implemented.
 * Never report a successful connection or deliver synthetic messages.
 */
public class WebSocketManager {
    public static let shared = WebSocketManager()
    public weak var delegate: WebSocketDelegate?
    public private(set) var isConnected: Bool = false

    public init() {}

    public func connect(url: URL, token: String) {
        isConnected = false
        delegate?.webSocketDidDisconnect()
    }

    public func simulateMessage(_ message: String) {
        // No network transport exists; synthetic messages are deliberately disabled.
    }

    public func disconnect() {
        isConnected = false
        delegate?.webSocketDidDisconnect()
    }
}
