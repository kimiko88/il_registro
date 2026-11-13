import SwiftUI

#if !SWIFT_PACKAGE
@main
#endif
struct StudentApp: App {
    init() {
        MdmConfigManager.shared.applyManagedConfiguration()
    }

    var body: some Scene {
        WindowGroup {
            StudentRootView()
        }
    }
}
