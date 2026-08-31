import SwiftUI

public struct SecretaryElectionModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let turnoutAndLists: String
    public let closingAndCrypto: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, turnoutAndLists: String, closingAndCrypto: String, statusText: String = "Seggio Aperto") {
        self.id = id
        self.title = title
        self.turnoutAndLists = turnoutAndLists
        self.closingAndCrypto = closingAndCrypto
        self.statusText = statusText
    }
}

public struct SecretaryElectionsView: View {
    public var elections: [SecretaryElectionModel]

    public init(elections: [SecretaryElectionModel] = []) {
        self.elections = elections
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if elections.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(elections) { item in
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
                                Text(item.turnoutAndLists)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.closingAndCrypto)
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
