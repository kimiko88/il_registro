import Foundation

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
    @Published public var certificates: [CertificateItemModel] = []

    public init() {
        loadData()
    }

    public func loadData() {
        users = [
            ManagedUserModel(id: "u1", name: "Prof.ssa Maria Rossi", role: "Docente"),
            ManagedUserModel(id: "u2", name: "Prof. Marco Bianchi", role: "Docente"),
            ManagedUserModel(id: "u3", name: "Mario Rossi (2B)", role: "Studente"),
            ManagedUserModel(id: "u4", name: "Giuseppe Rossi", role: "Genitore")
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

    public func issueCertificate(title: String) -> CertificateItemModel {
        let newCert = CertificateItemModel(id: UUID().uuidString, title: title, status: "pronto", pdfUrl: "/api/v1/cert/gen_\(UUID().uuidString).pdf")
        certificates.append(newCert)
        return newCert
    }
}
