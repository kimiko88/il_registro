import Foundation
import LocalAuthentication

public class BiometricAuthManager {
    public static let shared = BiometricAuthManager()

    public init() {}

    public func canEvaluatePolicy() -> Bool {
        let context = LAContext()
        var error: NSError?
        return context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error)
    }

    public func authenticateUser(reason: String = "Accedi rapidamente con Face ID o Touch ID") async -> Bool {
        let context = LAContext()
        do {
            return try await context.evaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, localizedReason: reason)
        } catch {
            return false
        }
    }
}
