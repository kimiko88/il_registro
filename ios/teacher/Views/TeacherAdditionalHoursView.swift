import SwiftUI

public struct TeacherAdditionalHourModel: Identifiable, Equatable {
    public let id: String
    public let totalHoursTitle: String
    public let details: String
    public let compensationDetails: String
    public let statusText: String

    public init(id: String = UUID().uuidString, totalHoursTitle: String, details: String, compensationDetails: String, statusText: String = "Validate") {
        self.id = id
        self.totalHoursTitle = totalHoursTitle
        self.details = details
        self.compensationDetails = compensationDetails
        self.statusText = statusText
    }
}

public struct TeacherAdditionalHoursView: View {
    public var summaries: [TeacherAdditionalHourModel]

    public init(summaries: [TeacherAdditionalHourModel] = []) {
        self.summaries = summaries
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if summaries.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(summaries) { s in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(s.totalHoursTitle)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(s.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(s.details)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(s.compensationDetails)
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
