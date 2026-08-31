import SwiftUI

public struct SecretaryStaffContractModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let classAndHours: String
    public let periodAndSidi: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, classAndHours: String, periodAndSidi: String, statusText: String = "Trasmesso a NoiPA") {
        self.id = id
        self.title = title
        self.classAndHours = classAndHours
        self.periodAndSidi = periodAndSidi
        self.statusText = statusText
    }
}

public struct SecretaryStaffContractsView: View {
    public var contracts: [SecretaryStaffContractModel]

    public init(contracts: [SecretaryStaffContractModel] = []) {
        self.contracts = contracts
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if contracts.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(contracts) { item in
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
                                Text(item.classAndHours)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.periodAndSidi)
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
