import SwiftUI

public struct CertificateProtocolModel: Identifiable, Equatable {
    public let id: String
    public let certificateType: String
    public let applicantAndUse: String
    public let protocolAndDate: String
    public let sealStatus: String

    public init(id: String = UUID().uuidString, certificateType: String, applicantAndUse: String, protocolAndDate: String, sealStatus: String = "Sigillo QR") {
        self.id = id
        self.certificateType = certificateType
        self.applicantAndUse = applicantAndUse
        self.protocolAndDate = protocolAndDate
        self.sealStatus = sealStatus
    }
}

public struct SecretaryCertificatesProtocolView: View {
    public var certificates: [CertificateProtocolModel]

    public init(certificates: [CertificateProtocolModel] = []) {
        self.certificates = certificates
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("certificate_requests", comment: ""))) {
                    if certificates.isEmpty {
                        Text(NSLocalizedString("certificate_requests", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(certificates) { cert in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(cert.certificateType)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(cert.sealStatus)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(cert.applicantAndUse)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(cert.protocolAndDate)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("secretary_dashboard_title", comment: ""))
        }
    }
}
