import SwiftUI

public struct SchoolCalendarEventModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let dateRange: String
    public let tag: String
    public let notes: String

    public init(id: String = UUID().uuidString, title: String, dateRange: String, tag: String, notes: String = "") {
        self.id = id
        self.title = title
        self.dateRange = dateRange
        self.tag = tag
        self.notes = notes
    }
}

public struct StudentSchoolCalendarView: View {
    public var events: [SchoolCalendarEventModel]

    public init(events: [SchoolCalendarEventModel] = []) {
        self.events = events
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("agenda_title", comment: ""))) {
                    if events.isEmpty {
                        Text(NSLocalizedString("agenda_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(events) { item in
                            VStack(alignment: .leading, spacing: 4) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.tag)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.blue.opacity(0.2))
                                        .foregroundColor(.blue)
                                        .cornerRadius(4)
                                }
                                Text(item.dateRange)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                if !item.notes.isEmpty {
                                    Text(item.notes)
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("agenda_title", comment: ""))
        }
    }
}
