import Foundation

public protocol SecretaryAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func fetchUsers(token: String) async throws -> [ManagedUserModel]
    func issueCertificate(token: String, studentId: String, type: String) async throws -> CertificateItemModel
}

public class MockSecretaryAPIService: SecretaryAPIServiceProtocol {
    public init() {}

    public func login(email: String, password: String) async throws -> String {
        guard !email.isEmpty && password.count >= 6 else {
            throw NSError(domain: "SecretaryAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Credenziali non valide"])
        }
        return "jwt_secretary_token_mock"
    }

    public func fetchUsers(token: String) async throws -> [ManagedUserModel] {
        guard !token.isEmpty else {
            throw NSError(domain: "SecretaryAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }
        return [
            ManagedUserModel(id: "u1", name: "Prof.ssa Maria Rossi", role: "Docente"),
            ManagedUserModel(id: "u2", name: "Prof. Marco Bianchi", role: "Docente")
        ]
    }

    public func issueCertificate(token: String, studentId: String, type: String) async throws -> CertificateItemModel {
        guard !token.isEmpty && !studentId.isEmpty else {
            throw NSError(domain: "SecretaryAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Dati non validi"])
        }
        return CertificateItemModel(
            id: UUID().uuidString,
            title: "Certificato \(type)",
            status: "pronto",
            pdfUrl: "/api/v1/certificates/cert_\(UUID().uuidString).pdf"
        )
    }
}
