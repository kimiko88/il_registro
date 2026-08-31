import SwiftUI

public struct TeacherShiftExchangeModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let subjectAndDates: String
    public let hoursAndStatus: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, subjectAndDates: String, hoursAndStatus: String, statusText: String = "Approvato") {
        self.id = id
        self.title = title
        self.subjectAndDates = subjectAndDates
        self.hoursAndStatus = hoursAndStatus
        self.statusText = statusText
    }
}

public struct TeacherShiftExchangeView: View {
    public var exchanges: [TeacherShiftExchangeModel]

    public init(exchanges: [TeacherShiftExchangeModel] = []) {
        self.exchanges = exchanges
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if exchanges.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(exchanges) { item in
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
                                Text(item.subjectAndDates)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.hoursAndStatus)
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
