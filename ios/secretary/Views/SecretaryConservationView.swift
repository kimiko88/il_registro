import SwiftUI

public struct SecretaryConservationPackageModel: Identifiable, Equatable {
    public let id: String
    public let packageTitle: String
    public let hashDetails: String
    public let conservatorAndReport: String
    public let statusText: String

    public init(id: String = UUID().uuidString, packageTitle: String, hashDetails: String, conservatorAndReport: String, statusText: String = "Conservato") {
        self.id = id
        self.packageTitle = packageTitle
        self.hashDetails = hashDetails
        self.conservatorAndReport = conservatorAndReport
        self.statusText = statusText
    }
}

public struct SecretaryConservationView: View {
    public var packages: [SecretaryConservationPackageModel]

    public init(packages: [SecretaryConservationPackageModel] = []) {
        self.packages = packages
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if packages.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(packages) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.packageTitle)
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
                                Text(item.hashDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.conservatorAndReport)
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
