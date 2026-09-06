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

public struct StudentAttendanceRecordModel: Identifiable, Equatable {
    public let id: String
    public let date: String
    public let type: String
    public let isJustified: Bool

    public init(id: String = UUID().uuidString, date: String, type: String, isJustified: Bool = true) {
        self.id = id
        self.date = date
        self.type = type
        self.isJustified = isJustified
    }
}

extension GradeItemModel {
    public var localizedType: String {
        switch type.trimmingCharacters(in: .whitespacesAndNewlines).lowercased() {
        case "written", "scritto":
            return studentLocalizedString("grade_type_written")
        case "oral", "orale":
            return studentLocalizedString("grade_type_oral")
        case "practical", "pratico":
            return studentLocalizedString("grade_type_practical")
        case "test":
            return studentLocalizedString("grade_type_test")
        case "project", "progetto":
            return studentLocalizedString("grade_type_project")
        case "lab", "laboratory", "laboratorio":
            return studentLocalizedString("grade_type_lab")
        default:
            return type
        }
    }
}

extension StudentAttendanceRecordModel {
    public var localizedType: String {
        switch type.trimmingCharacters(in: .whitespacesAndNewlines).lowercased() {
        case "present", "presenza":
            return studentLocalizedString("attendance_present")
        case "absent", "assenza":
            return studentLocalizedString("attendance_absent")
        case "late", "ritardo":
            return studentLocalizedString("attendance_late")
        case "early_exit", "uscita anticipata", "uscita_anticipata":
            return studentLocalizedString("attendance_early_exit")
        default:
            return type
        }
    }
}

public class StudentViewModel: ObservableObject {
    @Published public var grades: [GradeItemModel] = []
    @Published public var homework: [HomeworkItemModel] = []
    @Published public var attendance: [StudentAttendanceRecordModel] = []
    @Published public var isLoading: Bool = false
    @Published public var errorMessage: String? = nil

    private let apiService: StudentAPIServiceProtocol

    public init(apiService: StudentAPIServiceProtocol = HttpStudentAPIService()) {
        self.apiService = apiService
        // Le collezioni partono vuote in produzione; i dati vengono caricati dal server via loadFromDatabase(token:)
    }

    public func loadFromDatabase(token: String) async {
        await MainActor.run {
            self.isLoading = true
            self.errorMessage = nil
        }
        do {
            async let fetchedGrades = apiService.fetchGrades(token: token)
            async let fetchedHomework = apiService.fetchHomework(token: token)
            async let fetchedAttendance = apiService.fetchAttendance(token: token)

            let (gradesResult, homeworkResult, attendanceResult) = try await (fetchedGrades, fetchedHomework, fetchedAttendance)

            await MainActor.run {
                self.grades = gradesResult
                self.homework = homeworkResult
                self.attendance = attendanceResult
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

        attendance = [
            StudentAttendanceRecordModel(id: "1", date: "26 Ago", type: "Assenza", isJustified: true),
            StudentAttendanceRecordModel(id: "2", date: "18 Ago", type: "Ritardo", isJustified: true)
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
