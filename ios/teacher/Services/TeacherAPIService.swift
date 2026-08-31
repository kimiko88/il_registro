import Foundation

public protocol TeacherAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func signLessonHour(token: String, classId: String, topic: String) async throws -> Bool
    func submitGrade(token: String, studentId: String, grade: Double, type: String) async throws -> Bool
}

public class MockTeacherAPIService: TeacherAPIServiceProtocol {
    public init() {}

    public func login(email: String, password: String) async throws -> String {
        guard !email.isEmpty && password.count >= 6 else {
            throw NSError(domain: "TeacherAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Credenziali non valide"])
        }
        return "jwt_teacher_token_mock"
    }

    public func signLessonHour(token: String, classId: String, topic: String) async throws -> Bool {
        guard !token.isEmpty && !topic.isEmpty else {
            throw NSError(domain: "TeacherAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Argomento obbligatorio"])
        }
        return true
    }

    public func submitGrade(token: String, studentId: String, grade: Double, type: String) async throws -> Bool {
        guard !token.isEmpty && grade >= 1.0 && grade <= 10.0 else {
            throw NSError(domain: "TeacherAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Voto fuori dai limiti"])
        }
        return true
    }
}
