import SwiftUI

public struct TeacherGroupModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let teacherAndSchedule: String
    public let roomAndRegister: String
    public let countBadgeText: String

    public init(id: String = UUID().uuidString, title: String, teacherAndSchedule: String, roomAndRegister: String, countBadgeText: String = "18 Studenti") {
        self.id = id
        self.title = title
        self.teacherAndSchedule = teacherAndSchedule
        self.roomAndRegister = roomAndRegister
        self.countBadgeText = countBadgeText
    }
}

public struct TeacherGroupsView: View {
    public var groups: [TeacherGroupModel]

    public init(groups: [TeacherGroupModel] = []) {
        self.groups = groups
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if groups.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(groups) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.countBadgeText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.teacherAndSchedule)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.roomAndRegister)
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
