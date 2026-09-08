import SwiftUI

public struct TeacherFreeActivityModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let dateAndHour: String
    public let notes: String
    public let hoursBadgeText: String

    public init(id: String = UUID().uuidString, title: String, dateAndHour: String, notes: String, hoursBadgeText: String = "1 h") {
        self.id = id
        self.title = title
        self.dateAndHour = dateAndHour
        self.notes = notes
        self.hoursBadgeText = hoursBadgeText
    }
}

public struct TeacherFreeActivitiesView: View {
    public var activities: [TeacherFreeActivityModel]

    public init(activities: [TeacherFreeActivityModel] = []) {
        self.activities = activities
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if activities.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(activities) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.hoursBadgeText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.dateAndHour)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.notes)
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
