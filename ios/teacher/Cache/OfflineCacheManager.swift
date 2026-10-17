import Foundation

public class OfflineCacheManager {
    public static let shared = OfflineCacheManager()

    private var cachedClasses: [TeacherClassSessionModel] = []
    private var cachedStudents: [RollCallStudent] = []
    private var lastSync: Date?

    public init() {}

    public func saveClasses(_ classes: [TeacherClassSessionModel]) {
        cachedClasses = classes
        lastSync = Date()
    }

    public func getCachedClasses() -> [TeacherClassSessionModel] {
        return cachedClasses
    }

    public func saveStudents(_ students: [RollCallStudent]) {
        cachedStudents = students
        lastSync = Date()
    }

    public func getCachedStudents() -> [RollCallStudent] {
        return cachedStudents
    }

    public func isCacheValid(timeout: TimeInterval = 3600) -> Bool {
        guard let lastSync = lastSync else { return false }
        return Date().timeIntervalSince(lastSync) < timeout && !cachedClasses.isEmpty
    }

    public func clear() {
        cachedClasses.removeAll()
        cachedStudents.removeAll()
        lastSync = nil
    }
}
