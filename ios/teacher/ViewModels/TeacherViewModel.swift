import Foundation

public struct RollCallStudent: Identifiable, Equatable {
    public let id: String
    public let name: String
    public var status: String // presente, assente, ritardo

    public init(id: String, name: String, status: String = "presente") {
        self.id = id
        self.name = name
        self.status = status
    }
}

public struct GradeItemProposal: Identifiable, Equatable {
    public let id: String
    public let studentId: String
    public let grade: Double
    public let type: String
    public let comment: String

    public init(id: String, studentId: String, grade: Double, type: String, comment: String = "") {
        self.id = id
        self.studentId = studentId
        self.grade = grade
        self.type = type
        self.comment = comment
    }
}

public class TeacherViewModel: ObservableObject {
    @Published public var isHourSigned: Bool = false
    @Published public var lessonTopic: String = ""
    @Published public var students: [RollCallStudent] = []
    @Published public var grades: [GradeItemProposal] = []

    public init() {
        loadData()
    }

    public func loadData() {
        students = [
            RollCallStudent(id: "s1", name: "Banchi Andrea", status: "presente"),
            RollCallStudent(id: "s2", name: "Bianchi Elena", status: "presente"),
            RollCallStudent(id: "s3", name: "Ferrari Matteo", status: "assente"),
            RollCallStudent(id: "s4", name: "Rossi Sofia", status: "ritardo")
        ]
    }

    public func signLesson(topic: String) -> Bool {
        guard !topic.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else { return false }
        lessonTopic = topic
        isHourSigned = true
        return true
    }

    public func updateAttendance(studentId: String, status: String) -> Bool {
        guard let index = students.firstIndex(where: { $0.id == studentId }) else { return false }
        students[index].status = status
        return true
    }

    public func insertGrade(studentId: String, grade: Double, type: String) -> Bool {
        guard grade >= 1.0 && grade <= 10.0 else { return false }
        let newGrade = GradeItemProposal(id: UUID().uuidString, studentId: studentId, grade: grade, type: type)
        grades.append(newGrade)
        return true
    }
}
