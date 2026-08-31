import Foundation

public class OfflineCacheManager {
    public static let shared = OfflineCacheManager()

    private var cachedGrades: [GradeItemModel] = []
    private var lastSync: Date?

    public init() {}

    public func saveGrades(_ grades: [GradeItemModel]) {
        cachedGrades = grades
        lastSync = Date()
    }

    public func getCachedGrades() -> [GradeItemModel] {
        return cachedGrades
    }

    public func isCacheValid(timeout: TimeInterval = 3600) -> Bool {
        guard let lastSync = lastSync else { return false }
        return Date().timeIntervalSince(lastSync) < timeout && !cachedGrades.isEmpty
    }

    public func clear() {
        cachedGrades.removeAll()
        lastSync = nil
    }
}
