import SwiftUI

public struct SecretaryTimetableModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let structureDetails: String
    public let teacherDetails: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, structureDetails: String, teacherDetails: String, statusText: String = "Pubblicato") {
        self.id = id
        self.title = title
        self.structureDetails = structureDetails
        self.teacherDetails = teacherDetails
        self.statusText = statusText
    }
}

public struct SecretaryTimetableView: View {
    public var timetables: [SecretaryTimetableModel]

    public init(timetables: [SecretaryTimetableModel] = []) {
        self.timetables = timetables
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if timetables.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(timetables) { item in
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
                                Text(item.structureDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.teacherDetails)
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
