import Foundation

public struct SecretaryUserModel: Identifiable, Equatable {
    public let id: String
    public let firstName: String
    public let lastName: String
    public let email: String
    public let role: String

    public init(id: String, firstName: String, lastName: String, email: String, role: String) {
        self.id = id
        self.firstName = firstName
        self.lastName = lastName
        self.email = email
        self.role = role
    }
}

public struct ManagedUserModel: Identifiable, Equatable {
    public let id: String
    public let name: String
    public let role: String

    public init(id: String, name: String, role: String) {
        self.id = id
        self.name = name
        self.role = role
    }
}

public struct SecretaryClassModel: Identifiable, Equatable {
    public let id: String
    public let name: String

    public init(id: String, name: String) {
        self.id = id
        self.name = name
    }
}

public struct CertificateItemModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public var status: String
    public var pdfUrl: String

    public init(id: String, title: String, status: String = "pronto", pdfUrl: String = "") {
        self.id = id
        self.title = title
        self.status = status
        self.pdfUrl = pdfUrl
    }
}

public class SecretaryViewModel: ObservableObject {
    @Published public var users: [ManagedUserModel] = []
    @Published public var classes: [SecretaryClassModel] = []
    @Published public var certificates: [CertificateItemModel] = []
    @Published public var isLoading: Bool = false
    @Published public var errorMessage: String? = nil

    private let apiService: SecretaryAPIServiceProtocol

    public init(apiService: SecretaryAPIServiceProtocol = HttpSecretaryAPIService()) {
        self.apiService = apiService
    }

    public func loadFromDatabase(token: String) async {
        await MainActor.run {
            self.isLoading = true
            self.errorMessage = nil
        }
        do {
            async let fetchedUsers = apiService.fetchUsers(token: token)
            async let fetchedClasses = apiService.fetchClasses(token: token)
            async let fetchedCertificates = apiService.fetchCertificates(token: token)

            let (u, cl, cert) = try await (fetchedUsers, fetchedClasses, fetchedCertificates)

            await MainActor.run {
                self.users = u.map {
                    let fullName = "\($0.firstName) \($0.lastName)".trimmingCharacters(in: .whitespacesAndNewlines)
                    return ManagedUserModel(id: $0.id, name: fullName.isEmpty ? $0.email : fullName, role: $0.role)
                }
                self.classes = cl
                self.certificates = cert
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
        users = [
            ManagedUserModel(id: "u1", name: "Prof.ssa Maria Rossi", role: "Docente"),
            ManagedUserModel(id: "u2", name: "Prof. Marco Bianchi", role: "Docente")
        ]
        classes = [
            SecretaryClassModel(id: "c1", name: "Classe 1A"),
            SecretaryClassModel(id: "c2", name: "Classe 2A")
        ]
        certificates = [
            CertificateItemModel(id: "c1", title: "Certificato di Iscrizione e Frequenza", status: "pronto", pdfUrl: "/api/v1/cert/c1.pdf"),
            CertificateItemModel(id: "c2", title: "Certificato con Valutazioni", status: "pronto", pdfUrl: "/api/v1/cert/c2.pdf")
        ]
    }

    public func loadData() {
        users = [
            ManagedUserModel(id: "u1", name: "Prof.ssa Maria Rossi", role: "Docente"),
            ManagedUserModel(id: "u2", name: "Prof. Marco Bianchi", role: "Docente"),
            ManagedUserModel(id: "u3", name: "Mario Rossi (2B)", role: "Studente"),
            ManagedUserModel(id: "u4", name: "Giuseppe Rossi", role: "Genitore")
        ]

        classes = [
            SecretaryClassModel(id: "c1", name: "Classe 1A"),
            SecretaryClassModel(id: "c2", name: "Classe 2A")
        ]

        certificates = [
            CertificateItemModel(id: "c1", title: "Certificato di Iscrizione e Frequenza", status: "pronto", pdfUrl: "/api/v1/cert/c1.pdf"),
            CertificateItemModel(id: "c2", title: "Certificato con Valutazioni", status: "pronto", pdfUrl: "/api/v1/cert/c2.pdf")
        ]
    }

    public func addUser(name: String, role: String) -> Bool {
        guard !name.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else { return false }
        users.append(ManagedUserModel(id: UUID().uuidString, name: name, role: role))
        return true
    }

    public func addUser(firstName: String, lastName: String, email: String, role: String) -> Bool {
        return addUser(name: "\(firstName) \(lastName)", role: role)
    }

    public func issueCertificate(title: String) -> CertificateItemModel {
        let newCert = CertificateItemModel(id: UUID().uuidString, title: title, status: "pronto", pdfUrl: "/api/v1/cert/gen_\(UUID().uuidString).pdf")
        certificates.append(newCert)
        return newCert
    }

    public func requestCertificate(studentId: String, type: String) -> CertificateItemModel {
        return issueCertificate(title: "Certificato \(type)")
    }
}
