import SwiftUI

public struct SecretaryInvoicingFlowModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let collectedAndTransactions: String
    public let reasons: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, collectedAndTransactions: String, reasons: String, statusText: String = "Quadratura 100%") {
        self.id = id
        self.title = title
        self.collectedAndTransactions = collectedAndTransactions
        self.reasons = reasons
        self.statusText = statusText
    }
}

public struct SecretaryInvoicingView: View {
    public var flows: [SecretaryInvoicingFlowModel]

    public init(flows: [SecretaryInvoicingFlowModel] = []) {
        self.flows = flows
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if flows.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(flows) { item in
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
                                Text(item.collectedAndTransactions)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.reasons)
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
