import Foundation

public class OfflineCacheManager {
    public static let shared = OfflineCacheManager()

    private var cachedChildren: [ParentChildModel] = []
    private var cachedAbsences: [AbsenceModel] = []
    private var lastSync: Date?

    public init() {}

    public func saveChildren(_ children: [ParentChildModel]) {
        cachedChildren = children
        lastSync = Date()
    }

    public func getCachedChildren() -> [ParentChildModel] {
        return cachedChildren
    }

    public func saveAbsences(_ absences: [AbsenceModel]) {
        cachedAbsences = absences
        lastSync = Date()
    }

    public func getCachedAbsences() -> [AbsenceModel] {
        return cachedAbsences
    }

    public func isCacheValid(timeout: TimeInterval = 3600) -> Bool {
        guard let lastSync = lastSync else { return false }
        return Date().timeIntervalSince(lastSync) < timeout && !cachedChildren.isEmpty
    }

    public func clear() {
        cachedChildren.removeAll()
        cachedAbsences.removeAll()
        lastSync = nil
    }
}
