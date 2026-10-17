import Foundation

public class OfflineCacheManager {
    public static let shared = OfflineCacheManager()

    private var cachedUsers: [ManagedUserModel] = []
    private var cachedClasses: [SecretaryClassModel] = []
    private var cachedCertificates: [CertificateItemModel] = []
    private var lastSync: Date?

    public init() {}

    public func saveUsers(_ users: [ManagedUserModel]) {
        cachedUsers = users
        lastSync = Date()
    }

    public func getCachedUsers() -> [ManagedUserModel] {
        return cachedUsers
    }

    public func saveClasses(_ classes: [SecretaryClassModel]) {
        cachedClasses = classes
        lastSync = Date()
    }

    public func getCachedClasses() -> [SecretaryClassModel] {
        return cachedClasses
    }

    public func saveCertificates(_ certificates: [CertificateItemModel]) {
        cachedCertificates = certificates
        lastSync = Date()
    }

    public func getCachedCertificates() -> [CertificateItemModel] {
        return cachedCertificates
    }

    public func isCacheValid(timeout: TimeInterval = 3600) -> Bool {
        guard let lastSync = lastSync else { return false }
        return Date().timeIntervalSince(lastSync) < timeout && (!cachedUsers.isEmpty || !cachedCertificates.isEmpty)
    }

    public func clear() {
        cachedUsers.removeAll()
        cachedClasses.removeAll()
        cachedCertificates.removeAll()
        lastSync = nil
    }
}
