import SwiftUI

@main
struct StudentApp: App {
    @State private var isLoggedIn = false
    @State private var token = ""
    @State private var studentName = "Studente"

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                StudentDashboardView(token: token, studentName: studentName)
            } else {
                LoginView(isLoggedIn: $isLoggedIn, token: $token, studentName: $studentName)
            }
        }
    }
}
