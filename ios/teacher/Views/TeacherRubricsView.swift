import SwiftUI

public struct TeacherRubricModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let descriptors: String
    public let departmentAndSubjects: String
    public let levelsBadgeText: String

    public init(id: String = UUID().uuidString, title: String, descriptors: String, departmentAndSubjects: String, levelsBadgeText: String = "4 Livelli") {
        self.id = id
        self.title = title
        self.descriptors = descriptors
        self.departmentAndSubjects = departmentAndSubjects
        self.levelsBadgeText = levelsBadgeText
    }
}

public struct TeacherRubricsView: View {
    public var rubrics: [TeacherRubricModel]

    public init(rubrics: [TeacherRubricModel] = []) {
        self.rubrics = rubrics
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if rubrics.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(rubrics) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.levelsBadgeText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.descriptors)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.departmentAndSubjects)
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
