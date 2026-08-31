import SwiftUI

public struct TeacherHomeworkCorrectionModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let studentAndFile: String
    public let dateAndStatus: String
    public let countBadgeText: String

    public init(id: String = UUID().uuidString, title: String, studentAndFile: String, dateAndStatus: String, countBadgeText: String = "22 / 24") {
        self.id = id
        self.title = title
        self.studentAndFile = studentAndFile
        self.dateAndStatus = dateAndStatus
        self.countBadgeText = countBadgeText
    }
}

public struct TeacherHomeworkCorrectionsView: View {
    public var corrections: [TeacherHomeworkCorrectionModel]
    public var onOpenHomework: ((String) -> Void)?

    public init(corrections: [TeacherHomeworkCorrectionModel] = [], onOpenHomework: ((String) -> Void)? = nil) {
        self.corrections = corrections
        self.onOpenHomework = onOpenHomework
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if corrections.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(corrections) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.countBadgeText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.blue.opacity(0.2))
                                        .foregroundColor(.blue)
                                        .cornerRadius(4)
                                }
                                Text(item.studentAndFile)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.dateAndStatus)
                                    .font(.caption)
                                    .foregroundColor(.secondary)

                                Button(action: { onOpenHomework?(item.id) }) {
                                    Label(NSLocalizedString("teacher_dashboard_title", comment: ""), systemImage: "pencil.and.outline")
                                }
                                .buttonStyle(.borderedProminent)
                                .padding(.top, 2)
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
