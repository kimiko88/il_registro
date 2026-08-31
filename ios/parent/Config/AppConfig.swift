import Foundation

public struct AppConfig {
    public static var baseURL: String = "http://localhost:8080/api/v1"
    public static var wsURL: String = "ws://localhost:8080/api/v1/ws"
    public static let timeoutInterval: TimeInterval = 30.0
}
