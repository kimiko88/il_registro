import Foundation

public final class MdmConfigManager {
    public static let shared = MdmConfigManager()

    public func applyManagedConfiguration() {
        guard let managedConfig = UserDefaults.standard.dictionary(forKey: "com.apple.configuration.managed") else {
            return
        }

        if let serverURL = managedConfig["server_url"] as? String, !serverURL.isEmpty {
            AppConfig.baseURL = serverURL
        }

        if let wsURL = managedConfig["ws_url"] as? String, !wsURL.isEmpty {
            AppConfig.wsURL = wsURL
        }
    }
}
