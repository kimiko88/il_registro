import Foundation

public protocol ParentAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func fetchChildren(token: String) async throws -> [ParentChildModel]
    func fetchAbsences(token: String, childId: String) async throws -> [AbsenceModel]
    func justifyAbsence(token: String, childId: String, absenceId: String, reason: String) async throws -> Bool
}

public class HttpParentAPIService: ParentAPIServiceProtocol {
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
        var cleanEmail = email.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        if cleanEmail.hasPrefix("genitorea_") {
            cleanEmail = cleanEmail.replacingOccurrences(of: "genitorea_", with: "genitore2a_")
        } else if cleanEmail.hasPrefix("genitoreb_") {
            cleanEmail = cleanEmail.replacingOccurrences(of: "genitoreb_", with: "genitore2b_")
        }

        let url = baseURL.appendingPathComponent("auth/login")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["email": cleanEmail, "password": password]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "ParentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Credenziali non valide"])
        }

        if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any] {
            let userObj = json["user"] as? [String: Any]
            let role = userObj?["role"] as? String
            guard role == "parent" else {
                throw NSError(domain: "ParentAPI", code: 403, userInfo: [NSLocalizedDescriptionKey: "Accesso non consentito: questo account non appartiene a un genitore."])
            }
            if let token = json["token"] as? String ?? json["access_token"] as? String {
                return token
            }
        }
        throw NSError(domain: "ParentAPI", code: 500, userInfo: [NSLocalizedDescriptionKey: "Token mancante"])
    }

    public func fetchChildren(token: String) async throws -> [ParentChildModel] {
        let url = baseURL.appendingPathComponent("users/me/children")
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
                id: (dict["id"] as? String) ?? "\(index)",
                firstName: (dict["first_name"] as? String) ?? (dict["name"] as? String) ?? "Figlio",
                lastName: (dict["last_name"] as? String) ?? "",
                className: (dict["class_name"] as? String) ?? "Classe"
            )
        }
    }

    public func fetchAbsences(token: String, childId: String) async throws -> [AbsenceModel] {
        let url = baseURL.appendingPathComponent("attendance/child-attendance/\(childId)")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        guard let (data, response) = try? await URLSession.shared.data(for: request),
              let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode),
              let jsonArray = try? JSONSerialization.jsonObject(with: data) as? [[String: Any]] else {
            return []
        }

        return jsonArray.enumerated().map { index, dict in
            AbsenceModel(
                id: (dict["id"] as? String) ?? "\(index)",
                childId: childId,
                date: (dict["date"] as? String) ?? "",
                type: (dict["type"] as? String) ?? "Assenza",
                isJustified: (dict["is_justified"] as? Bool) ?? false,
                justificationNote: (dict["justification_note"] as? String) ?? (dict["reason"] as? String) ?? ""
            )
        }
    }

    public func justifyAbsence(token: String, childId: String, absenceId: String, reason: String) async throws -> Bool {
        let url = baseURL.appendingPathComponent("attendance/child/\(childId)/justify/\(absenceId)")
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
