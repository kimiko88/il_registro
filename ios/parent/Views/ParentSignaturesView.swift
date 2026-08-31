import SwiftUI

public struct ParentSignatureDocumentModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let description: String
    public var isSigned: Bool

    public init(id: String = UUID().uuidString, title: String, description: String, isSigned: Bool = false) {
        self.id = id
        self.title = title
        self.description = description
        self.isSigned = isSigned
    }
}

public struct ParentSignaturesView: View {
    public var documents: [ParentSignatureDocumentModel]
    public var onSignDocument: ((String) -> Void)?

    public init(documents: [ParentSignatureDocumentModel] = [], onSignDocument: ((String) -> Void)? = nil) {
        self.documents = documents
        self.onSignDocument = onSignDocument
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if documents.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(documents) { doc in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(doc.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(doc.isSigned ? "Firmato" : "Richiesto")
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(doc.isSigned ? Color.green.opacity(0.2) : Color.red.opacity(0.2))
                                        .foregroundColor(doc.isSigned ? .green : .red)
                                        .cornerRadius(4)
                                }
                                Text(doc.description)
                                    .font(.caption)
                                    .foregroundColor(.secondary)

                                if !doc.isSigned {
                                    Button(NSLocalizedString("sign_hour", comment: "")) {
                                        onSignDocument?(doc.id)
                                    }
                                    .buttonStyle(.borderedProminent)
                                    .padding(.top, 2)
                                }
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
