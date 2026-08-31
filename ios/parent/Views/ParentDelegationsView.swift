import SwiftUI

public struct ParentDelegateModel: Identifiable, Equatable {
    public let id: String
    public let delegateNameAndRelation: String
    public let documentAndExpiry: String
    public let statusDetails: String
    public let statusText: String

    public init(id: String = UUID().uuidString, delegateNameAndRelation: String, documentAndExpiry: String, statusDetails: String, statusText: String = "Valida") {
        self.id = id
        self.delegateNameAndRelation = delegateNameAndRelation
        self.documentAndExpiry = documentAndExpiry
        self.statusDetails = statusDetails
        self.statusText = statusText
    }
}

public struct ParentDelegationsView: View {
    public var delegates: [ParentDelegateModel]

    public init(delegates: [ParentDelegateModel] = []) {
        self.delegates = delegates
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if delegates.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(delegates) { d in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(d.delegateNameAndRelation)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(d.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(d.documentAndExpiry)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(d.statusDetails)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("parent_dashboard_title", comment: ""))
        }
    }
}
