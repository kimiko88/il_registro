import SwiftUI

public struct TeacherExamKioskModel: Identifiable, Equatable {
    public let id: String
    public let sessionTitle: String
    public let policyDetails: String
    public let statsDetails: String
    public let lockedCountText: String

    public init(id: String = UUID().uuidString, sessionTitle: String, policyDetails: String, statsDetails: String, lockedCountText: String) {
        self.id = id
        self.sessionTitle = sessionTitle
        self.policyDetails = policyDetails
        self.statsDetails = statsDetails
        self.lockedCountText = lockedCountText
    }
}

public struct TeacherExamKioskView: View {
    public var kiosks: [TeacherExamKioskModel]

    public init(kiosks: [TeacherExamKioskModel] = []) {
        self.kiosks = kiosks
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if kiosks.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(kiosks) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.sessionTitle)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.lockedCountText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.policyDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.statsDetails)
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
