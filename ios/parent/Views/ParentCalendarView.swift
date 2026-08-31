import SwiftUI

public struct ParentCalendarEventModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let dateText: String
    public let details: String
    public let tag: String

    public init(id: String = UUID().uuidString, title: String, dateText: String, details: String, tag: String) {
        self.id = id
        self.title = title
        self.dateText = dateText
        self.details = details
        self.tag = tag
    }
}

public struct ParentCalendarView: View {
    public var events: [ParentCalendarEventModel]

    public init(events: [ParentCalendarEventModel] = []) {
        self.events = events
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if events.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(events) { ev in
                            VStack(alignment: .leading, spacing: 4) {
                                HStack {
                                    Text(ev.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(ev.dateText)
                                        .font(.caption2)
                                        .foregroundColor(.blue)
                                }
                                Text(ev.details)
                                    .font(.subheadline)
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
