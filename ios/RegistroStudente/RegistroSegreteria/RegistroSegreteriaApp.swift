//
//  RegistroSegreteriaApp.swift
//  RegistroSegreteria
//

import SwiftUI
import SecretaryApp

@main
struct RegistroSegreteriaApp: App {
    init() {
        MdmConfigManager.shared.applyManagedConfiguration()
    }

    var body: some Scene {
        WindowGroup {
            SecretaryRootView()
        }
    }
}
