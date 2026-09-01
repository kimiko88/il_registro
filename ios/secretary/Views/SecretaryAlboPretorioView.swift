import SwiftUI

public struct SecretaryAlboActModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let actCodeAndCategory: String
    public let datesAndHash: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, actCodeAndCategory: String, datesAndHash: String, statusText: String = "In Pubblicazione") {
        self.id = id
        self.title = title
        self.actCodeAndCategory = actCodeAndCategory
        self.datesAndHash = datesAndHash
        self.statusText = statusText
    }
}

public struct SecretaryAlboPretorioView: View {
    public var acts: [SecretaryAlboActModel]

    public init(acts: [SecretaryAlboActModel] = []) {
        self.acts = acts
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if acts.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(acts) { item in
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
                                Text(item.actCodeAndCategory)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.datesAndHash)
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
