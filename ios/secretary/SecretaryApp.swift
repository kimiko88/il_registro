import SwiftUI

@main
struct SecretaryApp: App {
    @State private var isLoggedIn = false
    @State private var token = ""

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                SecretaryDashboardView(token: token)
            } else {
                LoginView(isLoggedIn: $isLoggedIn, token: $token)
            }
        }
    }
}
