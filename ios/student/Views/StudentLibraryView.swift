import SwiftUI

public struct LibraryLoanModel: Identifiable, Equatable {
    public let id: String
    public let bookTitleAndAuthor: String
    public let shelfAndInventory: String
    public let loanDates: String
    public let status: String

    public init(id: String = UUID().uuidString, bookTitleAndAuthor: String, shelfAndInventory: String, loanDates: String, status: String = "In Prestito") {
        self.id = id
        self.bookTitleAndAuthor = bookTitleAndAuthor
        self.shelfAndInventory = shelfAndInventory
        self.loanDates = loanDates
        self.status = status
    }
}

public struct StudentLibraryView: View {
    public var loans: [LibraryLoanModel]

    public init(loans: [LibraryLoanModel] = []) {
        self.loans = loans
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if loans.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(loans) { loan in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(loan.bookTitleAndAuthor)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(loan.status)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(loan.shelfAndInventory)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(loan.loanDates)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
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
