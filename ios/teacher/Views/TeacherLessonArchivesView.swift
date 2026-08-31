import SwiftUI

public struct TeacherLessonArchiveModel: Identifiable, Equatable {
    public let id: String
    public let lessonTitle: String
    public let topicDetails: String
    public let dateAndHomework: String
    public let statusText: String

    public init(id: String = UUID().uuidString, lessonTitle: String, topicDetails: String, dateAndHomework: String, statusText: String = "Firmata") {
        self.id = id
        self.lessonTitle = lessonTitle
        self.topicDetails = topicDetails
        self.dateAndHomework = dateAndHomework
        self.statusText = statusText
    }
}

public struct TeacherLessonArchivesView: View {
    public var lessons: [TeacherLessonArchiveModel]

    public init(lessons: [TeacherLessonArchiveModel] = []) {
        self.lessons = lessons
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if lessons.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(lessons) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.lessonTitle)
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
                                Text(item.topicDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.dateAndHomework)
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
