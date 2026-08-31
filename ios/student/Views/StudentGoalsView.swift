import SwiftUI

public struct StudentGoalModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let subjectAndTarget: String
    public let progress: Double
    public let progressText: String
    public let teacherFeedback: String

    public init(id: String = UUID().uuidString, title: String, subjectAndTarget: String, progress: Double, progressText: String, teacherFeedback: String) {
        self.id = id
        self.title = title
        self.subjectAndTarget = subjectAndTarget
        self.progress = progress
        self.progressText = progressText
        self.teacherFeedback = teacherFeedback
    }
}

public struct StudentGoalsView: View {
    public var goals: [StudentGoalModel]

    public init(goals: [StudentGoalModel] = []) {
        self.goals = goals
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if goals.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(goals) { goal in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(goal.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(goal.progressText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.blue.opacity(0.2))
                                        .foregroundColor(.blue)
                                        .cornerRadius(4)
                                }
                                Text(goal.subjectAndTarget)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)

                                ProgressView(value: goal.progress)
                                    .padding(.vertical, 2)

                                Text(goal.teacherFeedback)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
