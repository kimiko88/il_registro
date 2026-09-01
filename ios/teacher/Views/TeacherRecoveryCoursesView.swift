import SwiftUI

public struct TeacherRecoveryCourseModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let enrolledDetails: String
    public let examDateDetails: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, enrolledDetails: String, examDateDetails: String, statusText: String = "In Corso") {
        self.id = id
        self.title = title
        self.enrolledDetails = enrolledDetails
        self.examDateDetails = examDateDetails
        self.statusText = statusText
    }
}

public struct TeacherRecoveryCoursesView: View {
    public var courses: [TeacherRecoveryCourseModel]

    public init(courses: [TeacherRecoveryCourseModel] = []) {
        self.courses = courses
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if courses.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(courses) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
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
                                Text(item.enrolledDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.secondary)
                                Text(item.examDateDetails)
                                    .font(.caption)
                                    .foregroundColor(.blue)
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
