import SwiftUI

public struct ParentBullyingReportModel: Identifiable, Equatable {
    public let id: String
    public let practiceCode: String
    public let typeAndDate: String
    public let managerAndAction: String
    public let statusText: String

    public init(id: String = UUID().uuidString, practiceCode: String, typeAndDate: String, managerAndAction: String, statusText: String = "In Carico") {
        self.id = id
        self.practiceCode = practiceCode
        self.typeAndDate = typeAndDate
        self.managerAndAction = managerAndAction
        self.statusText = statusText
    }
}

public struct ParentCyberbullyingReportView: View {
    public var reports: [ParentBullyingReportModel]

    public init(reports: [ParentBullyingReportModel] = []) {
        self.reports = reports
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if reports.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(reports) { rep in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(rep.practiceCode)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(rep.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(rep.typeAndDate)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(rep.managerAndAction)
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
