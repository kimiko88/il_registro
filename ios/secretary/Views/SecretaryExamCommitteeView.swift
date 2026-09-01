import SwiftUI

public struct SecretaryExamCommitteeModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let presidentDetails: String
    public let membersAndCandidates: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, presidentDetails: String, membersAndCandidates: String, statusText: String = "Insediata") {
        self.id = id
        self.title = title
        self.presidentDetails = presidentDetails
        self.membersAndCandidates = membersAndCandidates
        self.statusText = statusText
    }
}

public struct SecretaryExamCommitteeView: View {
    public var committees: [SecretaryExamCommitteeModel]

    public init(committees: [SecretaryExamCommitteeModel] = []) {
        self.committees = committees
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if committees.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(committees) { item in
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
                                Text(item.presidentDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.membersAndCandidates)
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
