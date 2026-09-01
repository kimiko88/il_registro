import SwiftUI

public struct StudentBadgeModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let level: String
    public let competence: String
    public let verificationCode: String

    public init(id: String = UUID().uuidString, title: String, level: String, competence: String, verificationCode: String) {
        self.id = id
        self.title = title
        self.level = level
        self.competence = competence
        self.verificationCode = verificationCode
    }
}

public struct StudentBadgesView: View {
    public var badges: [StudentBadgeModel]

    public init(badges: [StudentBadgeModel] = []) {
        self.badges = badges
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if badges.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(badges) { badge in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(badge.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(badge.level)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.yellow.opacity(0.3))
                                        .foregroundColor(.orange)
                                        .cornerRadius(4)
                                }
                                Text(badge.competence)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(badge.verificationCode)
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
