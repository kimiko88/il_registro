import SwiftUI

public struct TeacherLabBookingModel: Identifiable, Equatable {
    public let id: String
    public let labName: String
    public let className: String
    public let dateTimeSlot: String
    public let topic: String
    public let status: String

    public init(id: String = UUID().uuidString, labName: String, className: String, dateTimeSlot: String, topic: String, status: String) {
        self.id = id
        self.labName = labName
        self.className = className
        self.dateTimeSlot = dateTimeSlot
        self.topic = topic
        self.status = status
    }
}

public struct TeacherLabBookingView: View {
    public var bookings: [TeacherLabBookingModel]

    public init(bookings: [TeacherLabBookingModel] = []) {
        self.bookings = bookings
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if bookings.isEmpty {
                        Text(NSLocalizedString("select_class", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(bookings) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.labName)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.status)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text("\(item.className) • \(item.dateTimeSlot)")
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.topic)
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
