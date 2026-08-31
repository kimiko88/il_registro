import SwiftUI

public struct StudentCertificationModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let issuerAndLevel: String
    public let protocolAndDate: String
    public let status: String

    public init(id: String = UUID().uuidString, title: String, issuerAndLevel: String, protocolAndDate: String, status: String = "Validata") {
        self.id = id
        self.title = title
        self.issuerAndLevel = issuerAndLevel
        self.protocolAndDate = protocolAndDate
        self.status = status
    }
}

public struct StudentCertificationsView: View {
    public var certifications: [StudentCertificationModel]

    public init(certifications: [StudentCertificationModel] = []) {
        self.certifications = certifications
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if certifications.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(certifications) { cert in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(cert.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(cert.status)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(cert.issuerAndLevel)
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
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
