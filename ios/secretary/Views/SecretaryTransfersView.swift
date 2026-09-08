import SwiftUI

public struct SecretaryTransferModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let destinationAndClass: String
    public let datesAndProtocol: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, destinationAndClass: String, datesAndProtocol: String, statusText: String = "Rilasciato") {
        self.id = id
        self.title = title
        self.destinationAndClass = destinationAndClass
        self.datesAndProtocol = datesAndProtocol
        self.statusText = statusText
    }
}

public struct SecretaryTransfersView: View {
    public var transfers: [SecretaryTransferModel]

    public init(transfers: [SecretaryTransferModel] = []) {
        self.transfers = transfers
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if transfers.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(transfers) { item in
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
                                Text(item.destinationAndClass)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.datesAndProtocol)
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
