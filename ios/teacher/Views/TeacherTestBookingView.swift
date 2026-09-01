import SwiftUI

public struct TeacherTestBookingModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let dateAndHours: String
    public let topics: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, dateAndHours: String, topics: String, statusText: String = "Nessun Conflitto") {
        self.id = id
        self.title = title
        self.dateAndHours = dateAndHours
        self.topics = topics
        self.statusText = statusText
    }
}

public struct TeacherTestBookingView: View {
    public var bookings: [TeacherTestBookingModel]

    public init(bookings: [TeacherTestBookingModel] = []) {
        self.bookings = bookings
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if bookings.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(bookings) { item in
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
                                Text(item.dateAndHours)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.topics)
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
