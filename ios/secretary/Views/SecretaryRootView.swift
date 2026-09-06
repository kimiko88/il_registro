import SwiftUI

public struct SecretaryRootView: View {
    @State private var isLoggedIn: Bool
    @State private var token: String

    public init(startLoggedIn: Bool = false, initialToken: String = "") {
        self._isLoggedIn = State(initialValue: startLoggedIn)
        self._token = State(initialValue: initialToken)
    }

    public var body: some View {
        if isLoggedIn {
            SecretaryDashboardView(token: token)
        } else {
            LoginView(isLoggedIn: $isLoggedIn, token: $token)
        }
    }
}
