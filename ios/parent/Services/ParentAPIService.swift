import Foundation

public protocol ParentAPIServiceProtocol {
    func login(email: String, password: String) async throws -> String
    func fetchChildren(token: String) async throws -> [ChildItemModel]
    func justifyAbsence(token: String, absenceId: String, note: String) async throws -> Bool
}

public class MockParentAPIService: ParentAPIServiceProtocol {
    public init() {}

    public func login(email: String, password: String) async throws -> String {
        guard !email.isEmpty && password.count >= 6 else {
            throw NSError(domain: "ParentAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Credenziali non valide"])
        }
        return "jwt_parent_token_mock"
    }

    public func fetchChildren(token: String) async throws -> [ChildItemModel] {
        guard !token.isEmpty else {
            throw NSError(domain: "ParentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }
        return [
            ChildItemModel(id: "c1", firstName: "Marco", lastName: "Rossi", className: "Classe 2A"),
            ChildItemModel(id: "c2", firstName: "Giulia", lastName: "Rossi", className: "Classe 4B")
        ]
    }

    public func justifyAbsence(token: String, absenceId: String, note: String) async throws -> Bool {
        guard !token.isEmpty && !note.isEmpty else {
            throw NSError(domain: "ParentAPI", code: 400, userInfo: [NSLocalizedDescriptionKey: "Nota obbligatoria"])
        }
        return true
    }
}
