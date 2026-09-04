import SwiftUI

#if !SWIFT_PACKAGE
@main
#endif
struct StudentApp: App {
    var body: some Scene {
        WindowGroup {
            StudentRootView()
        }
    }
}
