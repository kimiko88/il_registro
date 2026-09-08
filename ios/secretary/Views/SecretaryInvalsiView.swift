import SwiftUI

public struct SecretaryInvalsiSessionModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let studentsAndLabs: String
    public let proctorsAndStatus: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, studentsAndLabs: String, proctorsAndStatus: String, statusText: String = "Completata") {
        self.id = id
        self.title = title
        self.studentsAndLabs = studentsAndLabs
        self.proctorsAndStatus = proctorsAndStatus
        self.statusText = statusText
    }
}

public struct SecretaryInvalsiView: View {
    public var sessions: [SecretaryInvalsiSessionModel]

    public init(sessions: [SecretaryInvalsiSessionModel] = []) {
        self.sessions = sessions
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if sessions.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(sessions) { item in
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
                                Text(item.studentsAndLabs)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.proctorsAndStatus)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("secretary_dashboard_title", comment: ""))
        }
    }
}
