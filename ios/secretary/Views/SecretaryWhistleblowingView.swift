import SwiftUI

public struct WhistleblowingReportItem: Identifiable, Equatable {
    public let id: String
    public let code: String
    public let status: String
    public let date: String

    public init(id: String = UUID().uuidString, code: String, status: String, date: String) {
        self.id = id
        self.code = code
        self.status = status
        self.date = date
    }
}

public struct SecretaryWhistleblowingView: View {
    public var reports: [WhistleblowingReportItem]

    public init(reports: [WhistleblowingReportItem] = []) {
        self.reports = reports
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("audit_logs", comment: ""))) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                                .fontWeight(.bold)
                            Spacer()
                            Text("ANAC")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                    }
                    .padding(.vertical, 4)
                }

                if !reports.isEmpty {
                    Section(header: Text("Reports")) {
                        ForEach(reports) { item in
                            HStack {
                                Text(item.code)
                                Spacer()
                                Text(item.status)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("audit_logs", comment: ""))
        }
    }
}
