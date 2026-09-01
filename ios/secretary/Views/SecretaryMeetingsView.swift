import SwiftUI

public struct SecretaryMeetingModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let dateAndLocation: String
    public let agendaDetails: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, dateAndLocation: String, agendaDetails: String, statusText: String = "Convocato") {
        self.id = id
        self.title = title
        self.dateAndLocation = dateAndLocation
        self.agendaDetails = agendaDetails
        self.statusText = statusText
    }
}

public struct SecretaryMeetingsView: View {
    public var meetings: [SecretaryMeetingModel]

    public init(meetings: [SecretaryMeetingModel] = []) {
        self.meetings = meetings
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if meetings.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(meetings) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.blue.opacity(0.2))
                                        .foregroundColor(.blue)
                                        .cornerRadius(4)
                                }
                                Text(item.dateAndLocation)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.agendaDetails)
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
