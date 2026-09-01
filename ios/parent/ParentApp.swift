import SwiftUI

@main
struct ParentApp: App {
    @State private var isLoggedIn = false
    @State private var token = ""

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                ParentDashboardView(token: token)
            } else {
                LoginView(isLoggedIn: $isLoggedIn, token: $token)
            }
        }
    }
}
