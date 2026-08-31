import Foundation

public struct GradeItemModel: Identifiable, Equatable {
    public let id: String
    public let subject: String
    public let grade: Double
    public let weight: Double
    public let type: String
    public let date: String

    public init(id: String, subject: String, grade: Double, weight: Double = 1.0, type: String, date: String) {
        self.id = id
        self.subject = subject
        self.grade = grade
        self.weight = weight
        self.type = type
        self.date = date
    }
}

public struct HomeworkItemModel: Identifiable, Equatable {
    public let id: String
    public let subject: String
    public let taskDescription: String
    public let dueDate: String
    public var isCompleted: Bool

    public init(id: String, subject: String, taskDescription: String, dueDate: String, isCompleted: Bool = false) {
        self.id = id
        self.subject = subject
        self.taskDescription = taskDescription
        self.dueDate = dueDate
        self.isCompleted = isCompleted
    }
}

public class StudentViewModel: ObservableObject {
    @Published public var grades: [GradeItemModel] = []
    @Published public var homework: [HomeworkItemModel] = []

    public init() {
        loadData()
    }

    public func loadData() {
        grades = [
            GradeItemModel(id: "1", subject: "Matematica", grade: 8.5, weight: 1.0, type: "Scritto", date: "30 Ago"),
            GradeItemModel(id: "2", subject: "Matematica", grade: 7.0, weight: 1.0, type: "Orale", date: "15 Ago"),
            GradeItemModel(id: "3", subject: "Italiano", grade: 8.0, weight: 1.0, type: "Tema", date: "28 Ago"),
            GradeItemModel(id: "4", subject: "Inglese", grade: 9.0, weight: 1.0, type: "Pratico", date: "25 Ago"),
            GradeItemModel(id: "5", subject: "Fisica", grade: 7.0, weight: 1.0, type: "Scritto", date: "22 Ago")
        ]

        homework = [
            HomeworkItemModel(id: "1", subject: "Matematica", taskDescription: "Disequazioni", dueDate: "Domani", isCompleted: false),
            HomeworkItemModel(id: "2", subject: "Fisica", taskDescription: "Relazione moto", dueDate: "Tra 2 giorni", isCompleted: false),
            HomeworkItemModel(id: "3", subject: "Italiano", taskDescription: "Capitolo 8", dueDate: "Tra 3 giorni", isCompleted: true)
        ]
    }

    public func calculateGPA() -> Double {
        guard !grades.isEmpty else { return 0.0 }
        let totalWeighted = grades.reduce(0.0) { $0 + ($1.grade * $1.weight) }
        let totalWeight = grades.reduce(0.0) { $0 + $1.weight }
        guard totalWeight > 0 else { return 0.0 }
        return (totalWeighted / totalWeight * 10).rounded() / 10
    }

    public func toggleHomework(id: String) -> Bool {
        guard let index = homework.firstIndex(where: { $0.id == id }) else { return false }
        homework[index].isCompleted.toggle()
        return homework[index].isCompleted
    }
}
