import Foundation

/// Public accessor for student string localization.
/// Pass languageCode to resolve a specific locale independently of the device language.
public func studentLocalizedString(
    _ key: String,
    tableName: String? = nil,
    value: String = "",
    comment: String = "",
    languageCode: String? = nil
) -> String {
    let bundle = Bundle.module
    let requestedLanguage = (languageCode ?? Locale.preferredLanguages.first ?? "it")
        .replacingOccurrences(of: "_", with: "-")

    // Try the most specific locale first (for example zh-Hans-CN), then
    // progressively fall back to its parent locale (zh-Hans, then zh).
    let components = requestedLanguage.split(separator: "-")
    let languageCandidates = components.isEmpty
        ? ["it"]
        : (1...components.count).reversed().map {
            components.prefix($0).joined(separator: "-")
        }

    for candidate in languageCandidates {
        guard let languagePath = bundle.path(forResource: candidate, ofType: "lproj"),
              let languageBundle = Bundle(path: languagePath) else {
            continue
        }

        let localized = languageBundle.localizedString(
            forKey: key,
            value: value.isEmpty ? nil : value,
            table: tableName
        )
        if localized != key {
            return localized
        }
    }

    // Standard bundle resolution (including the package's development language).
    let localized = bundle.localizedString(
        forKey: key,
        value: value.isEmpty ? nil : value,
        table: tableName
    )
    if localized != key {
        return localized
    }

    // Italian is the school's default localization and the final package fallback.
    if let italianPath = bundle.path(forResource: "it", ofType: "lproj"),
       let italianBundle = Bundle(path: italianPath) {
        let italianValue = italianBundle.localizedString(
            forKey: key,
            value: value.isEmpty ? nil : value,
            table: tableName
        )
        if italianValue != key {
            return italianValue
        }
    }

    // Final fallback to the host application.
    return Bundle.main.localizedString(
        forKey: key,
        value: value.isEmpty ? key : value,
        table: tableName
    )
}

/// Module-internal overload of NSLocalizedString that routes to Bundle.module.
func NSLocalizedString(_ key: String, comment: String) -> String {
    studentLocalizedString(key, comment: comment)
}

/// Module-internal full overload of NSLocalizedString.
func NSLocalizedString(
    _ key: String,
    tableName: String? = nil,
    bundle: Bundle? = nil,
    value: String = "",
    comment: String
) -> String {
    if let bundle {
        return bundle.localizedString(
            forKey: key,
            value: value.isEmpty ? nil : value,
            table: tableName
        )
    }
    return studentLocalizedString(key, tableName: tableName, value: value, comment: comment)
}
