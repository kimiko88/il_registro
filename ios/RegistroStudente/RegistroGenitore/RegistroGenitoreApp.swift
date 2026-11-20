//
//  RegistroGenitoreApp.swift
//  RegistroGenitore
//

import SwiftUI
import ParentApp

@main
struct RegistroGenitoreApp: App {
    init() {
        MdmConfigManager.shared.applyManagedConfiguration()
    }

    var body: some Scene {
        WindowGroup {
            ParentRootView()
        }
    }
}
