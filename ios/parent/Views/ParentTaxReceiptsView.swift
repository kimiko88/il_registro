import SwiftUI

public struct ParentTaxReceiptModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let totalAndDetails: String
    public let itemsSummary: String
    public let tag: String

    public init(id: String = UUID().uuidString, title: String, totalAndDetails: String, itemsSummary: String, tag: String = "Detraibile 19%") {
        self.id = id
        self.title = title
        self.totalAndDetails = totalAndDetails
        self.itemsSummary = itemsSummary
        self.tag = tag
    }
}

public struct ParentTaxReceiptsView: View {
    public var receipts: [ParentTaxReceiptModel]
    public var onDownloadPdf: ((String) -> Void)?

    public init(receipts: [ParentTaxReceiptModel] = [], onDownloadPdf: ((String) -> Void)? = nil) {
        self.receipts = receipts
        self.onDownloadPdf = onDownloadPdf
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if receipts.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(receipts) { r in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(r.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(r.tag)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(r.totalAndDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(r.itemsSummary)
                                    .font(.caption)
                                    .foregroundColor(.secondary)

                                Button(action: { onDownloadPdf?(r.id) }) {
                                    Label("PDF", systemImage: "arrow.down.doc.fill")
                                }
                                .buttonStyle(.bordered)
                                .padding(.top, 2)
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
