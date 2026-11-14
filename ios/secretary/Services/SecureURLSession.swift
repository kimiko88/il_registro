import Foundation

private final class RejectRedirectDelegate: NSObject, URLSessionTaskDelegate {
    func urlSession(
        _ session: URLSession,
        task: URLSessionTask,
        willPerformHTTPRedirection response: HTTPURLResponse,
        newRequest request: URLRequest,
        completionHandler: @escaping (URLRequest?) -> Void
    ) {
        // Never forward an authenticated API request to a redirect destination.
        completionHandler(nil)
    }
}

enum SecureURLSession {
    static let shared = URLSession(
        configuration: .ephemeral,
        delegate: RejectRedirectDelegate(),
        delegateQueue: nil
    )
}
