import Foundation

public struct ChildItemModel: Identifiable, Equatable {
    public let id: String
    public let firstName: String
    public let lastName: String
    public let className: String

    public init(id: String, firstName: String, lastName: String, className: String) {
        self.id = id
        self.firstName = firstName
        self.lastName = lastName
        self.className = className
    }
}

public struct ParentChildModel: Identifiable, Equatable {
    public let id: String
    public let firstName: String
    public let lastName: String
    public let className: String

    public init(id: String, firstName: String, lastName: String, className: String) {
        self.id = id
        self.firstName = firstName
        self.lastName = lastName
        self.className = className
    }
}

public struct AbsenceModel: Identifiable, Equatable {
    public let id: String
    public let childId: String
    public let date: String
    public let type: String
    public var isJustified: Bool
    public var justificationNote: String

    public init(id: String, childId: String, date: String, type: String, isJustified: Bool = false, justificationNote: String = "") {
        self.id = id
        self.childId = childId
        self.date = date
        self.type = type
        self.isJustified = isJustified
        self.justificationNote = justificationNote
    }
}

public class ParentViewModel: ObservableObject {
    @Published public var children: [ChildItemModel] = []
    @Published public var selectedChildId: String = ""
    @Published public var absences: [AbsenceModel] = []
    @Published public var isLoading: Bool = false
    @Published public var errorMessage: String? = nil

    private let apiService: ParentAPIServiceProtocol

    public init(apiService: ParentAPIServiceProtocol = HttpParentAPIService()) {
        self.apiService = apiService
    }

    public func loadFromDatabase(token: String) async {
        await MainActor.run {
            self.isLoading = true
            self.errorMessage = nil
        }
        do {
            let fetchedChildren = try await apiService.fetchChildren(token: token)
            await MainActor.run {
                self.children = fetchedChildren.map {
                    ChildItemModel(id: $0.id, firstName: $0.firstName, lastName: $0.lastName, className: $0.className)
                }
                if let first = self.children.first {
                    self.selectedChildId = first.id
                }
                self.isLoading = false
            }
        } catch {
            await MainActor.run {
                self.errorMessage = error.localizedDescription
                self.isLoading = false
            }
        }
    }

    public func loadSampleData() {
        loadData()
    }

    public func loadData() {
        children = [
            ChildItemModel(id: "c1", firstName: "Marco", lastName: "Rossi", className: "Classe 2A"),
            ChildItemModel(id: "c2", firstName: "Giulia", lastName: "Rossi", className: "Classe 4B")
        ]
        selectedChildId = "c1"

        absences = [
            AbsenceModel(id: "a1", childId: "c1", date: "2026-08-26", type: "Assenza", isJustified: false),
            AbsenceModel(id: "a2", childId: "c1", date: "2026-08-18", type: "Ritardo", isJustified: false),
            AbsenceModel(id: "a3", childId: "c2", date: "2026-08-20", type: "Assenza", isJustified: false)
        ]
    }

    public func selectChild(id: String) {
        selectedChildId = id
    }

    public func getAbsencesForSelectedChild() -> [AbsenceModel] {
        return absences.filter { $0.childId == selectedChildId }
    }

    public func justifyAbsence(id: String, reason: String) -> Bool {
        guard let index = absences.firstIndex(where: { $0.id == id }) else { return false }
        absences[index].isJustified = true
        absences[index].justificationNote = reason
        return true
    }
}
