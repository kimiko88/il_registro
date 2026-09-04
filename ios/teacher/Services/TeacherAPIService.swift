import Foundation

public protocol TeacherAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func fetchClasses(token: String) async throws -> [TeacherClassSessionModel]
    func fetchClassStudents(token: String, classId: String) async throws -> [RollCallStudent]
    func signLesson(token: String, classId: String, topic: String) async throws -> TeacherClassSessionModel
    func submitRollCall(token: String, classId: String, records: [TeacherRollCallRecord]) async throws -> Bool
}

public class HttpTeacherAPIService: TeacherAPIServiceProtocol {
    private let baseURL: URL

    public init(baseURL: URL? = nil) {
        if let baseURL = baseURL {
            self.baseURL = baseURL
        } else if let url = URL(string: AppConfig.baseURL) {
            self.baseURL = url
        } else {
            self.baseURL = URL(string: "https://registro-backend-fdu2.onrender.com/api/v1")!
        }
    }

    public func login(email: String, password: String) async throws -> String {
        let cleanEmail = email.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        let url = baseURL.appendingPathComponent("auth/login")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["email": cleanEmail, "password": password]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "TeacherAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Credenziali errate"])
        }

        if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any] {
            let userObj = json["user"] as? [String: Any]
            let role = userObj?["role"] as? String
            guard role == "teacher" || role == "coordinator" else {
                throw NSError(domain: "TeacherAPI", code: 403, userInfo: [NSLocalizedDescriptionKey: "Accesso non consentito: questo account non appartiene a un docente."])
            }
            if let token = json["token"] as? String ?? json["access_token"] as? String {
                return token
            }
        }
        throw NSError(domain: "TeacherAPI", code: 500, userInfo: [NSLocalizedDescriptionKey: "Token mancante"])
    }

    public func fetchClasses(token: String) async throws -> [TeacherClassSessionModel] {
        let url = baseURL.appendingPathComponent("teacher/classes")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        guard let (data, response) = try? await URLSession.shared.data(for: request),
              let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode),
              let jsonArray = try? JSONSerialization.jsonObject(with: data) as? [[String: Any]] else {
            return []
        }

        return jsonArray.enumerated().map { index, dict in
            TeacherClassSessionModel(
                id: (dict["id"] as? String) ?? "\(index)",
                className: (dict["name"] as? String) ?? (dict["class_name"] as? String) ?? "Classe",
                subject: (dict["subject"] as? String) ?? (dict["subject_name"] as? String) ?? "Materia",
                isSigned: false,
                lessonTopic: ""
            )
        }
    }

    public func fetchClassStudents(token: String, classId: String) async throws -> [RollCallStudent] {
        let url = baseURL.appendingPathComponent("classes/\(classId)")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        guard let (data, response) = try? await URLSession.shared.data(for: request),
              let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode),
              let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
              let studentsArray = json["students"] as? [[String: Any]] else {
            return []
        }

        return studentsArray.enumerated().map { index, dict in
            let firstName = dict["first_name"] as? String ?? ""
            let lastName = dict["last_name"] as? String ?? ""
            let fullName = "\(firstName) \(lastName)".trimmingCharacters(in: .whitespacesAndNewlines)
            return RollCallStudent(
                id: (dict["id"] as? String) ?? "\(index)",
                name: fullName.isEmpty ? ((dict["name"] as? String) ?? "Studente") : fullName,
                status: "presente"
            )
        }
    }

    public func signLesson(token: String, classId: String, topic: String) async throws -> TeacherClassSessionModel {
        let url = baseURL.appendingPathComponent("lessons")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["class_id": classId, "topic": topic]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "TeacherAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Errore firma lezione"])
        }

        guard let dict = try JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            return TeacherClassSessionModel(id: classId, className: "Classe", subject: "Materia", isSigned: true, lessonTopic: topic)
        }

        return TeacherClassSessionModel(
            id: (dict["id"] as? String) ?? classId,
            className: (dict["class_name"] as? String) ?? "Classe",
            subject: (dict["subject_name"] as? String) ?? "Materia",
            isSigned: true,
            lessonTopic: topic
        )
    }

    public func submitRollCall(token: String, classId: String, records: [TeacherRollCallRecord]) async throws -> Bool {
        let url = baseURL.appendingPathComponent("attendance/mark-bulk")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let recordsArray = records.map { ["student_id": $0.studentId, "status": $0.status, "reason": $0.note] }
        let body: [String: Any] = ["class_id": classId, "records": recordsArray]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (_, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            return false
        }
        return true
    }
}
