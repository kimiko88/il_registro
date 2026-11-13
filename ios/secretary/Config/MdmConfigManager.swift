import Foundation

public final class MdmConfigManager {
    public static let shared = MdmConfigManager()

    private func validatedURL(_ value: String, requiredScheme: String) -> String? {
        guard let components = URLComponents(string: value.trimmingCharacters(in: .whitespacesAndNewlines)),
              components.scheme?.lowercased() == requiredScheme,
              let host = components.host, !host.isEmpty,
              components.user == nil,
              components.password == nil,
              components.query == nil,
              components.fragment == nil else {
            return nil
        }
        return components.string?.trimmingCharacters(in: CharacterSet(charactersIn: "/"))
    }

    public func applyManagedConfiguration() {
        guard let managedConfig = UserDefaults.standard.dictionary(forKey: "com.apple.configuration.managed") else {
            return
        }

        if let serverURL = managedConfig["server_url"] as? String,
           let validated = validatedURL(serverURL, requiredScheme: "https") {
            AppConfig.baseURL = validated
        }

        if let wsURL = managedConfig["ws_url"] as? String,
           let validated = validatedURL(wsURL, requiredScheme: "wss") {
            AppConfig.wsURL = validated
        }
    }
}
