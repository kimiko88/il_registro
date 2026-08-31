import Foundation

public protocol ParentAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func fetchChildren(token: String) async throws -> [ParentChildModel]
    func justifyAbsence(token: String, absenceId: String, reason: String) async throws -> Bool
}

public class HttpParentAPIService: ParentAPIServiceProtocol {
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
            throw NSError(domain: "ParentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Credenziali non valide"])
        }

        if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any],
           let token = json["token"] as? String ?? json["access_token"] as? String {
            return token
        }
        throw NSError(domain: "ParentAPI", code: 500, userInfo: [NSLocalizedDescriptionKey: "Token mancante"])
    }

    public func fetchChildren(token: String) async throws -> [ParentChildModel] {
        let url = baseURL.appendingPathComponent("parents/my-children")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "ParentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }

        guard let jsonArray = try JSONSerialization.jsonObject(with: data) as? [[String: Any]] else {
            return []
        }

        return jsonArray.enumerated().map { index, dict in
            ParentChildModel(
                id: "\(dict["id"] ?? index)",
                firstName: (dict["first_name"] as? String) ?? (dict["name"] as? String) ?? "Figlio",
                lastName: (dict["last_name"] as? String) ?? "",
                className: (dict["class_name"] as? String) ?? "Classe"
            )
        }
    }

    public func justifyAbsence(token: String, absenceId: String, reason: String) async throws -> Bool {
        let url = baseURL.appendingPathComponent("attendance/\(absenceId)/justify")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["reason": reason]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (_, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            return false
        }
        return true
    }
}
