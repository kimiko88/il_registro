import SwiftUI

public struct TeacherSupportStudentModel: Identifiable, Equatable {
    public let id: String
    public let studentAndClass: String
    public let peiType: String
    public let diaryUpdateDetails: String
    public let hoursBadgeText: String

    public init(id: String = UUID().uuidString, studentAndClass: String, peiType: String, diaryUpdateDetails: String, hoursBadgeText: String = "18 h/sett") {
        self.id = id
        self.studentAndClass = studentAndClass
        self.peiType = peiType
        self.diaryUpdateDetails = diaryUpdateDetails
        self.hoursBadgeText = hoursBadgeText
    }
}

public struct TeacherSupportView: View {
    public var students: [TeacherSupportStudentModel]

    public init(students: [TeacherSupportStudentModel] = []) {
        self.students = students
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if students.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(students) { s in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(s.studentAndClass)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(s.hoursBadgeText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.blue.opacity(0.2))
                                        .foregroundColor(.blue)
                                        .cornerRadius(4)
                                }
                                Text(s.peiType)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(s.diaryUpdateDetails)
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
