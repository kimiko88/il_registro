import SwiftUI

@main
struct SecretaryApp: App {
    @State private var isLoggedIn = false

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                SecretaryDashboardView()
            } else {
                LoginView(isLoggedIn: $isLoggedIn)
            }
        }
    }
}
