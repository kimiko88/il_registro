import Foundation

public protocol StudentAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func fetchGrades(token: String) async throws -> [GradeItemModel]
    func fetchHomework(token: String) async throws -> [HomeworkItemModel]
}

public class MockStudentAPIService: StudentAPIServiceProtocol {
    public init() {}

    public func login(email: String, password: String) async throws -> String {
        guard !email.isEmpty && password.count >= 6 else {
            throw NSError(domain: "StudentAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Credenziali non valide"])
        }
        return "jwt_student_token_mock"
    }

    public func fetchGrades(token: String) async throws -> [GradeItemModel] {
        guard !token.isEmpty else {
            throw NSError(domain: "StudentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }
        return [
            GradeItemModel(id: "1", subject: "Matematica", grade: 8.5, weight: 1.0, type: "Scritto", date: "30 Ago"),
            GradeItemModel(id: "2", subject: "Italiano", grade: 8.0, weight: 1.0, type: "Tema", date: "28 Ago")
        ]
    }

    public func fetchHomework(token: String) async throws -> [HomeworkItemModel] {
        guard !token.isEmpty else {
            throw NSError(domain: "StudentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }
        return [
            HomeworkItemModel(id: "1", subject: "Matematica", taskDescription: "Disequazioni", dueDate: "Domani", isCompleted: false)
        ]
    }
}
