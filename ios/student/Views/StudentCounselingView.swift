import SwiftUI

public struct StudentCounselingModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let dateTime: String
    public let counselorAndLocation: String
    public let status: String

    public init(id: String = UUID().uuidString, title: String, dateTime: String, counselorAndLocation: String, status: String = "Confermato") {
        self.id = id
        self.title = title
        self.dateTime = dateTime
        self.counselorAndLocation = counselorAndLocation
        self.status = status
    }
}

public struct StudentCounselingView: View {
    public var appointments: [StudentCounselingModel]

    public init(appointments: [StudentCounselingModel] = []) {
        self.appointments = appointments
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if appointments.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(appointments) { appt in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(appt.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(appt.status)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(appt.dateTime)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(appt.counselorAndLocation)
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
