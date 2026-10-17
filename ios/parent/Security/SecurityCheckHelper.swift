import Foundation
#if canImport(UIKit)
import UIKit
#endif

public class SecurityCheckHelper {
    public static let shared = SecurityCheckHelper()

    private let suspiciousFilePaths = [
        "/Applications/Cydia.app",
        "/Library/MobileSubstrate/MobileSubstrate.dylib",
        "/bin/bash",
        "/usr/sbin/sshd",
        "/etc/apt",
        "/private/var/lib/apt/",
        "/usr/bin/ssh",
        "/private/var/stash",
        "/System/Library/LaunchDaemons/com.saurik.Cydia.Startup.plist"
    ]

    public init() {}

    public func isDeviceJailbroken() -> Bool {
        return checkSuspiciousFiles() || checkSandboxEscape() || checkCydiaURLScheme()
    }

    private func checkSuspiciousFiles() -> Bool {
        let fileManager = FileManager.default
        return suspiciousFilePaths.any { fileManager.fileExists(atPath: $0) }
    }

    private func checkSandboxEscape() -> Bool {
        let testPath = "/private/jailbreak_test_\(UUID().uuidString).txt"
        do {
            try "test".write(toFile: testPath, atomically: true, encoding: .utf8)
            try? FileManager.default.removeItem(atPath: testPath)
            return true
        } catch {
            return false
        }
    }

    private func checkCydiaURLScheme() -> Bool {
        #if canImport(UIKit)
        if let url = URL(string: "cydia://package/com.example.package") {
            #if !targetEnvironment(simulator)
            if UIApplication.shared.canOpenURL(url) {
                return true
            }
            #endif
        }
        #endif
        return false
    }

    public func isSimulator() -> Bool {
        #if targetEnvironment(simulator)
        return true
        #else
        return false
        #endif
    }

    public func getSecurityReport() -> [String: Bool] {
        return [
            "isJailbroken": isDeviceJailbroken(),
            "isSimulator": isSimulator(),
            "hasSuspiciousFiles": checkSuspiciousFiles(),
            "hasSandboxEscape": checkSandboxEscape()
        ]
    }
}

private extension Array {
    func any(_ predicate: (Element) -> Bool) -> Bool {
        for element in self {
            if predicate(element) { return true }
        }
        return false
    }
}
