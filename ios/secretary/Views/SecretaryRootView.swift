import SwiftUI

public struct SecretaryRootView: View {
    @State private var isLoggedIn: Bool
    @State private var token: String

    public init(startLoggedIn: Bool = false, initialToken: String = "") {
        self._isLoggedIn = State(initialValue: startLoggedIn)
        self._token = State(initialValue: initialToken)
    }

    private func logout() {
        token = ""
        isLoggedIn = false
    }

    public var body: some View {
        if isLoggedIn {
            SecretaryDashboardView(token: token, onLogout: logout)
        } else {
            LoginView(isLoggedIn: $isLoggedIn, token: $token)
        }
    }
}
