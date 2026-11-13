import SwiftUI

#if !SWIFT_PACKAGE
@main
#endif
struct TeacherApp: App {
    @State private var isLoggedIn = false
    @State private var token = ""

    init() {
        MdmConfigManager.shared.applyManagedConfiguration()
    }

    var body: some Scene {
        WindowGroup {
            if isLoggedIn {
                TeacherRegisterView(token: token, onLogout: {
                    let tokenToRevoke = token
                    Task { await HttpTeacherAPIService().revokeSession(token: tokenToRevoke) }
                    token = ""
                    isLoggedIn = false
                })
            } else {
                LoginView(isLoggedIn: $isLoggedIn, token: $token)
            }
        }
    }
}
