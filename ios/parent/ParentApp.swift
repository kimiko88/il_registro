import SwiftUI

@main
struct ParentApp: App {
    @State private var isLoggedIn = false

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                ParentDashboardView()
            } else {
                LoginView(isLoggedIn: $isLoggedIn)
            }
        }
    }
}
