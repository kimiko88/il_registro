import SwiftUI

public struct StudentQuizModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let subjectAndDetails: String
    public let scoreAndAttempts: String
    public let status: String

    public init(id: String = UUID().uuidString, title: String, subjectAndDetails: String, scoreAndAttempts: String, status: String = "Aperto") {
        self.id = id
        self.title = title
        self.subjectAndDetails = subjectAndDetails
        self.scoreAndAttempts = scoreAndAttempts
        self.status = status
    }
}

public struct StudentELearningView: View {
    public var quizzes: [StudentQuizModel]
    public var onStartQuiz: ((String) -> Void)?

    public init(quizzes: [StudentQuizModel] = [], onStartQuiz: ((String) -> Void)? = nil) {
        self.quizzes = quizzes
        self.onStartQuiz = onStartQuiz
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if quizzes.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(quizzes) { quiz in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(quiz.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(quiz.status)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(quiz.subjectAndDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(quiz.scoreAndAttempts)
                                    .font(.caption)
                                    .foregroundColor(.secondary)

                                Button(NSLocalizedString("dashboard_title", comment: "")) {
                                    onStartQuiz?(quiz.id)
                                }
                                .buttonStyle(.borderedProminent)
                                .padding(.top, 4)
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
