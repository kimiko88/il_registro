import SwiftUI

public struct TeacherCoTeachingModel: Identifiable, Equatable {
    public let id: String
    public let subjectTitle: String
    public let teachersAndRoles: String
    public let classAndDetails: String
    public let statusText: String

    public init(id: String = UUID().uuidString, subjectTitle: String, teachersAndRoles: String, classAndDetails: String, statusText: String = "Firmata") {
        self.id = id
        self.subjectTitle = subjectTitle
        self.teachersAndRoles = teachersAndRoles
        self.classAndDetails = classAndDetails
        self.statusText = statusText
    }
}

public struct TeacherCoTeachingView: View {
    public var sessions: [TeacherCoTeachingModel]

    public init(sessions: [TeacherCoTeachingModel] = []) {
        self.sessions = sessions
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if sessions.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(sessions) { s in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(s.subjectTitle)
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
                                Text(s.teachersAndRoles)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(s.classAndDetails)
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
