import SwiftUI

public struct TeacherLeaveModel: Identifiable, Equatable {
    public let id: String
    public let leaveTitle: String
    public let dateTimeDetails: String
    public let substituteAndProtocol: String
    public let statusText: String

    public init(id: String = UUID().uuidString, leaveTitle: String, dateTimeDetails: String, substituteAndProtocol: String, statusText: String = "Approvato") {
        self.id = id
        self.leaveTitle = leaveTitle
        self.dateTimeDetails = dateTimeDetails
        self.substituteAndProtocol = substituteAndProtocol
        self.statusText = statusText
    }
}

public struct TeacherLeavesView: View {
    public var leaves: [TeacherLeaveModel]

    public init(leaves: [TeacherLeaveModel] = []) {
        self.leaves = leaves
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if leaves.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(leaves) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.leaveTitle)
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
                                Text(item.dateTimeDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.substituteAndProtocol)
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
