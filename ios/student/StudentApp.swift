import SwiftUI

@main
struct StudentApp: App {
    @State private var isLoggedIn = false
    @State private var studentName = "Mario Rossi"

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                StudentDashboardView()
            } else {
                LoginView(isLoggedIn: $isLoggedIn, studentName: $studentName)
            }
        }
    }
}
