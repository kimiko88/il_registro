import SwiftUI

@main
struct TeacherApp: App {
    @State private var isLoggedIn = false
    @State private var token = ""

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                TeacherRegisterView(token: token)
            } else {
                LoginView(isLoggedIn: $isLoggedIn, token: $token)
            }
        }
    }
}
