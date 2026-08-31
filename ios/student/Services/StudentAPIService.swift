import Foundation

public protocol StudentAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func fetchGrades(token: String) async throws -> [GradeItemModel]
    func fetchHomework(token: String) async throws -> [HomeworkItemModel]
}

public class HttpStudentAPIService: StudentAPIServiceProtocol {
    private let baseURL: URL

    public init(baseURL: URL = URL(string: "https://api.scuola.registro.it/api/v1")!) {
        self.baseURL = baseURL
    }

    public func login(email: String, password: String) async throws -> String {
        let url = baseURL.appendingPathComponent("auth/login")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["email": email, "password": password]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "StudentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Credenziali non valide o errore server"])
        }

        if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any],
           let token = json["token"] as? String ?? json["access_token"] as? String {
            return token
        }
        throw NSError(domain: "StudentAPI", code: 500, userInfo: [NSLocalizedDescriptionKey: "Token non presente nella risposta"])
    }

    public func fetchGrades(token: String) async throws -> [GradeItemModel] {
        let url = baseURL.appendingPathComponent("grades/my-grades")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "StudentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }

        guard let jsonArray = try JSONSerialization.jsonObject(with: data) as? [[String: Any]] else {
            return []
        }

        return jsonArray.enumerated().map { index, dict in
            GradeItemModel(
                id: "\(dict["id"] ?? index)",
                subject: (dict["subject_name"] as? String) ?? (dict["subject"] as? String) ?? "Materia",
                grade: (dict["grade"] as? Double) ?? 0.0,
                weight: (dict["weight"] as? Double) ?? 1.0,
                type: (dict["type"] as? String) ?? "Orale",
                date: (dict["date"] as? String) ?? ""
            )
        }
    }

    public func fetchHomework(token: String) async throws -> [HomeworkItemModel] {
        let url = baseURL.appendingPathComponent("agenda/my-homework")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "StudentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }

        guard let jsonArray = try JSONSerialization.jsonObject(with: data) as? [[String: Any]] else {
            return []
        }

        return jsonArray.enumerated().map { index, dict in
            HomeworkItemModel(
                id: "\(dict["id"] ?? index)",
                subject: (dict["subject_name"] as? String) ?? (dict["subject"] as? String) ?? "",
                taskDescription: (dict["description"] as? String) ?? (dict["title"] as? String) ?? "",
                dueDate: (dict["due_date"] as? String) ?? "",
                isCompleted: (dict["is_completed"] as? Bool) ?? false
            )
        }
    }
}
