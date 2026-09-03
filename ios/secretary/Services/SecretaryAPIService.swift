import Foundation

public protocol SecretaryAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func fetchUsers(token: String) async throws -> [SecretaryUserModel]
    func requestCertificate(token: String, studentId: String, type: String) async throws -> String
}

public class HttpSecretaryAPIService: SecretaryAPIServiceProtocol {
    private let baseURL: URL

    public init(baseURL: URL = URL(string: "https://registro-backend-fdu2.onrender.com/api/v1")!) {
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
            throw NSError(domain: "SecretaryAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Credenziali errate"])
        }

        if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any] {
            let userObj = json["user"] as? [String: Any]
            let role = userObj?["role"] as? String
            let allowedRoles = ["secretary", "admin", "superadmin", "principal", "vice_principal"]
            guard let role = role, allowedRoles.contains(role) else {
                throw NSError(domain: "SecretaryAPI", code: 403, userInfo: [NSLocalizedDescriptionKey: "Accesso non consentito: questo account non ha i permessi di segreteria o amministrazione."])
            }
            if let token = json["token"] as? String ?? json["access_token"] as? String {
                return token
            }
        }
        throw NSError(domain: "SecretaryAPI", code: 500, userInfo: [NSLocalizedDescriptionKey: "Token mancante"])
    }

    public func fetchUsers(token: String) async throws -> [SecretaryUserModel] {
        let url = baseURL.appendingPathComponent("users")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "SecretaryAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }

        let jsonObject = try JSONSerialization.jsonObject(with: data)
        let jsonArray: [[String: Any]]
        if let dict = jsonObject as? [String: Any], let users = dict["users"] as? [[String: Any]] {
            jsonArray = users
        } else if let array = jsonObject as? [[String: Any]] {
            jsonArray = array
        } else {
            return []
        }

        return jsonArray.enumerated().map { index, dict in
            SecretaryUserModel(
                id: "\(dict["id"] ?? index)",
                firstName: (dict["first_name"] as? String) ?? (dict["name"] as? String) ?? "Utente",
                lastName: (dict["last_name"] as? String) ?? "",
                email: (dict["email"] as? String) ?? "",
                role: (dict["role"] as? String) ?? "staff"
            )
        }
    }

    public func requestCertificate(token: String, studentId: String, type: String) async throws -> String {
        let url = baseURL.appendingPathComponent("certificates")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["student_id": studentId, "type": type]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "SecretaryAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Errore emissione certificato"])
        }

        if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any],
           let id = json["id"] as? String {
            return id
        }
        return "cert_ok"
    }
}
