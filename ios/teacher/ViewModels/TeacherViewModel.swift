import Foundation

public struct TeacherClassSessionModel: Identifiable, Equatable {
    public let id: String
    public let className: String
    public let subject: String
    public var isSigned: Bool
    public var lessonTopic: String

    public init(id: String, className: String, subject: String, isSigned: Bool = false, lessonTopic: String = "") {
        self.id = id
        self.className = className
        self.subject = subject
        self.isSigned = isSigned
        self.lessonTopic = lessonTopic
    }
}

public struct TeacherRollCallRecord: Identifiable, Equatable {
    public let id: String
    public let studentId: String
    public var status: String
    public var note: String

    public init(id: String = UUID().uuidString, studentId: String, status: String, note: String = "") {
        self.id = id
        self.studentId = studentId
        self.status = status
        self.note = note
    }
}

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
    @Published public var currentSession: TeacherClassSessionModel? = nil
    @Published public var isLoading: Bool = false
    @Published public var errorMessage: String? = nil

    private let apiService: TeacherAPIServiceProtocol

    public init(apiService: TeacherAPIServiceProtocol = HttpTeacherAPIService()) {
        self.apiService = apiService
    }

    public func signLessonViaApi(token: String, classId: String, topic: String) async -> Bool {
        await MainActor.run {
            self.isLoading = true
            self.errorMessage = nil
        }
        do {
            let session = try await apiService.signLesson(token: token, classId: classId, topic: topic)
            await MainActor.run {
                self.currentSession = session
                self.isHourSigned = session.isSigned
                self.lessonTopic = session.lessonTopic
                self.isLoading = false
            }
            return true
        } catch {
            await MainActor.run {
                self.errorMessage = error.localizedDescription
                self.isLoading = false
            }
            return false
        }
    }

    public func loadSampleSession() {
        currentSession = TeacherClassSessionModel(id: "3A", className: "Classe 3A", subject: "Matematica", isSigned: false, lessonTopic: "")
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
        currentSession?.isSigned = true
        currentSession?.lessonTopic = topic
        return true
    }

    public func updateAttendance(studentId: String, status: String) -> Bool {
        guard let index = students.firstIndex(where: { $0.id == studentId }) else { return false }
        students[index].status = status
        return true
    }

    public func toggleAttendance(studentId: String, status: String) -> Bool {
        return updateAttendance(studentId: studentId, status: status)
    }

    public func insertGrade(studentId: String, grade: Double, type: String) -> Bool {
        guard grade >= 1.0 && grade <= 10.0 else { return false }
        let newGrade = GradeItemProposal(id: UUID().uuidString, studentId: studentId, grade: grade, type: type)
        grades.append(newGrade)
        return true
    }
}
