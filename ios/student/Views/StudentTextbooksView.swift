import SwiftUI

public struct StudentTextbookModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let authorAndPublisher: String
    public let isbnAndPrice: String
    public let status: String
    public let isOwned: Bool

    public init(id: String = UUID().uuidString, title: String, authorAndPublisher: String, isbnAndPrice: String, status: String = "Da Acquistare", isOwned: Bool = false) {
        self.id = id
        self.title = title
        self.authorAndPublisher = authorAndPublisher
        self.isbnAndPrice = isbnAndPrice
        self.status = status
        self.isOwned = isOwned
    }
}

public struct StudentTextbooksView: View {
    public var textbooks: [StudentTextbookModel]

    public init(textbooks: [StudentTextbookModel] = []) {
        self.textbooks = textbooks
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if textbooks.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(textbooks) { book in
                            VStack(alignment: .leading, spacing: 4) {
                                HStack {
                                    Text(book.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(book.status)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(book.isOwned ? Color.green.opacity(0.2) : Color.blue.opacity(0.2))
                                        .foregroundColor(book.isOwned ? .green : .blue)
                                        .cornerRadius(4)
                                }
                                Text(book.authorAndPublisher)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                                Text(book.isbnAndPrice)
                                    .font(.caption)
                                    .foregroundColor(.blue)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
