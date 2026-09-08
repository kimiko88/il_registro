import SwiftUI

public struct ParentReportCardModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let summaryAndAverage: String
    public let statusTag: String

    public init(id: String = UUID().uuidString, title: String, summaryAndAverage: String, statusTag: String = "Sigillata (QR)") {
        self.id = id
        self.title = title
        self.summaryAndAverage = summaryAndAverage
        self.statusTag = statusTag
    }
}

public struct ParentReportCardsArchiveView: View {
    public var reportCards: [ParentReportCardModel]
    public var onDownloadPdf: ((String) -> Void)?

    public init(reportCards: [ParentReportCardModel] = [], onDownloadPdf: ((String) -> Void)? = nil) {
        self.reportCards = reportCards
        self.onDownloadPdf = onDownloadPdf
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if reportCards.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(reportCards) { rc in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(rc.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(rc.statusTag)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(rc.summaryAndAverage)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)

                                Button(action: { onDownloadPdf?(rc.id) }) {
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
