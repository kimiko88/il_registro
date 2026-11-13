import SwiftUI

#if !SWIFT_PACKAGE
@main
#endif
struct ParentApp: App {
    @State private var isLoggedIn = false
    @State private var token = ""

    init() {
        MdmConfigManager.shared.applyManagedConfiguration()
    }

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
