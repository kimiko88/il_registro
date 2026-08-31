import SwiftUI

@main
struct TeacherApp: App {
    @State private var isLoggedIn = false

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                TeacherRegisterView()
            } else {
                LoginView(isLoggedIn: $isLoggedIn)
            }
        }
    }
}
