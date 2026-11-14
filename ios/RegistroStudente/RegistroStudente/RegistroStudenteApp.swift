//
//  RegistroStudenteApp.swift
//  RegistroStudente
//

import SwiftUI
import StudentApp

@main
struct RegistroStudenteApp: App {
    init() {
        MdmConfigManager.shared.applyManagedConfiguration()
    }

    var body: some Scene {
        WindowGroup {
            StudentRootView()
        }
    }
}
