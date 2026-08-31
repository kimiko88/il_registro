import SwiftUI

public struct TeacherTripModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let classesAndDate: String
    public let guidesAndAuthorizations: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, classesAndDate: String, guidesAndAuthorizations: String, statusText: String = "Approvata") {
        self.id = id
        self.title = title
        self.classesAndDate = classesAndDate
        self.guidesAndAuthorizations = guidesAndAuthorizations
        self.statusText = statusText
    }
}

public struct TeacherTripsView: View {
    public var trips: [TeacherTripModel]

    public init(trips: [TeacherTripModel] = []) {
        self.trips = trips
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if trips.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(trips) { item in
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
                                Text(item.classesAndDate)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.guidesAndAuthorizations)
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
