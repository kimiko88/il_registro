import SwiftUI

public struct SecretaryStorageModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let usageRatio: Double
    public let description: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, usageRatio: Double, description: String, statusText: String = "Capienza Ottimale") {
        self.id = id
        self.title = title
        self.usageRatio = usageRatio
        self.description = description
        self.statusText = statusText
    }
}

public struct SecretaryMaterialsView: View {
    public var storageItems: [SecretaryStorageModel]

    public init(storageItems: [SecretaryStorageModel] = []) {
        self.storageItems = storageItems
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if storageItems.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(storageItems) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.statusText)
                                        .font(.caption2)
                                        .foregroundColor(.blue)
                                }
                                ProgressView(value: item.usageRatio)
                                    .padding(.vertical, 2)
                                Text(item.description)
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
