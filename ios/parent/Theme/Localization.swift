import Foundation

public func parentLocalizedString(_ key: String, tableName: String? = nil, value: String = "", comment: String = "") -> String {
    let bundle = Bundle.module

    let preferredLang = Locale.preferredLanguages.first?.prefix(2).lowercased() ?? "it"
    if let langPath = bundle.path(forResource: String(preferredLang), ofType: "lproj"),
       let langBundle = Bundle(path: langPath) {
        let val = langBundle.localizedString(forKey: key, value: value.isEmpty ? nil : value, table: tableName)
        if val != key {
            return val
        }
    }

    let localized = bundle.localizedString(forKey: key, value: value.isEmpty ? nil : value, table: tableName)
    if localized != key {
        return localized
    }

    if let itPath = bundle.path(forResource: "it", ofType: "lproj"),
       let itBundle = Bundle(path: itPath) {
        let itVal = itBundle.localizedString(forKey: key, value: value.isEmpty ? nil : value, table: tableName)
        if itVal != key {
            return itVal
        }
    }

    return Bundle.main.localizedString(forKey: key, value: value.isEmpty ? key : value, table: tableName)
}

func NSLocalizedString(_ key: String, comment: String) -> String {
    parentLocalizedString(key, comment: comment)
}

func NSLocalizedString(_ key: String, tableName: String? = nil, bundle: Bundle? = nil, value: String = "", comment: String) -> String {
    if let b = bundle {
        return b.localizedString(forKey: key, value: value.isEmpty ? nil : value, table: tableName)
    }
    return parentLocalizedString(key, tableName: tableName, value: value, comment: comment)
}
