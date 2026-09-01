import SwiftUI

public struct TeacherColloquioBookingModel: Identifiable, Equatable {
    public let id: String
    public let parentAndStudent: String
    public let timeSlot: String
    public let location: String
    public let reason: String

    public init(id: String = UUID().uuidString, parentAndStudent: String, timeSlot: String, location: String, reason: String) {
        self.id = id
        self.parentAndStudent = parentAndStudent
        self.timeSlot = timeSlot
        self.location = location
        self.reason = reason
    }
}

public struct TeacherColloquiView: View {
    public var colloqui: [TeacherColloquioBookingModel]

    public init(colloqui: [TeacherColloquioBookingModel] = []) {
        self.colloqui = colloqui
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if colloqui.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(colloqui) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.parentAndStudent)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.timeSlot)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.blue.opacity(0.2))
                                        .foregroundColor(.blue)
                                        .cornerRadius(4)
                                }
                                Text(item.location)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.reason)
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
