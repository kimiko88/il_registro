import SwiftUI

public struct SecretaryInventoryAssetModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let inventoryNumberAndCategory: String
    public let locationAndFunding: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, inventoryNumberAndCategory: String, locationAndFunding: String, statusText: String = "Inventariato") {
        self.id = id
        self.title = title
        self.inventoryNumberAndCategory = inventoryNumberAndCategory
        self.locationAndFunding = locationAndFunding
        self.statusText = statusText
    }
}

public struct SecretaryInventoryView: View {
    public var items: [SecretaryInventoryAssetModel]

    public init(items: [SecretaryInventoryAssetModel] = []) {
        self.items = items
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if items.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(items) { item in
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
                                Text(item.inventoryNumberAndCategory)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.locationAndFunding)
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
