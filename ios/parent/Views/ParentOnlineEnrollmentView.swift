import SwiftUI

public struct ParentEnrollmentModel: Identifiable, Equatable {
    public let id: String
    public let applicationCode: String
    public let studentAndTrack: String
    public let branchAndLanguage: String
    public let statusText: String

    public init(id: String = UUID().uuidString, applicationCode: String, studentAndTrack: String, branchAndLanguage: String, statusText: String = "Accettata") {
        self.id = id
        self.applicationCode = applicationCode
        self.studentAndTrack = studentAndTrack
        self.branchAndLanguage = branchAndLanguage
        self.statusText = statusText
    }
}

public struct ParentOnlineEnrollmentView: View {
    public var applications: [ParentEnrollmentModel]

    public init(applications: [ParentEnrollmentModel] = []) {
        self.applications = applications
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if applications.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(applications) { app in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(app.applicationCode)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(app.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(app.studentAndTrack)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(app.branchAndLanguage)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("parent_dashboard_title", comment: ""))
        }
    }
}
