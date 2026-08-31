import SwiftUI

public struct SecretaryInsuranceClaimModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let injuredAndDate: String
    public let protocolsAndPolicy: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, injuredAndDate: String, protocolsAndPolicy: String, statusText: String = "Inviato a INAIL") {
        self.id = id
        self.title = title
        self.injuredAndDate = injuredAndDate
        self.protocolsAndPolicy = protocolsAndPolicy
        self.statusText = statusText
    }
}

public struct SecretaryInsuranceClaimsView: View {
    public var claims: [SecretaryInsuranceClaimModel]

    public init(claims: [SecretaryInsuranceClaimModel] = []) {
        self.claims = claims
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if claims.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(claims) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.injuredAndDate)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.protocolsAndPolicy)
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
