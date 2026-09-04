import SwiftUI

#if !SWIFT_PACKAGE
@main
#endif
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
