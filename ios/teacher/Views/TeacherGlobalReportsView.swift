import SwiftUI

public struct TeacherGlobalReportModel: Identifiable, Equatable {
    public let id: String
    public let studentAndTitle: String
    public let profileDetails: String
    public let orientationAdvice: String
    public let statusText: String

    public init(id: String = UUID().uuidString, studentAndTitle: String, profileDetails: String, orientationAdvice: String, statusText: String = "Deliberato") {
        self.id = id
        self.studentAndTitle = studentAndTitle
        self.profileDetails = profileDetails
        self.orientationAdvice = orientationAdvice
        self.statusText = statusText
    }
}

public struct TeacherGlobalReportsView: View {
    public var reports: [TeacherGlobalReportModel]

    public init(reports: [TeacherGlobalReportModel] = []) {
        self.reports = reports
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if reports.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(reports) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.studentAndTitle)
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
                                Text(item.profileDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.orientationAdvice)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("teacher_dashboard_title", comment: ""))
        }
    }
}
