//
//  RegistroDocenteApp.swift
//  RegistroDocente
//

import SwiftUI
import TeacherApp

@main
struct RegistroDocenteApp: App {
    init() {
        MdmConfigManager.shared.applyManagedConfiguration()
    }

    var body: some Scene {
        WindowGroup {
            TeacherRootView()
        }
    }
}
