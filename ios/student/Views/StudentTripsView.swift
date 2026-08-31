import SwiftUI

public struct StudentTripModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let period: String
    public let chaperones: String
    public let paymentStatus: String
    public let status: String

    public init(id: String = UUID().uuidString, title: String, period: String, chaperones: String, paymentStatus: String, status: String = "Approvato") {
        self.id = id
        self.title = title
        self.period = period
        self.chaperones = chaperones
        self.paymentStatus = paymentStatus
        self.status = status
    }
}

public struct StudentTripsView: View {
    public var trips: [StudentTripModel]

    public init(trips: [StudentTripModel] = []) {
        self.trips = trips
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if trips.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(trips) { trip in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(trip.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(trip.status)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(trip.period)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(trip.chaperones)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                                Text(trip.paymentStatus)
                                    .font(.caption2)
                                    .foregroundColor(.green)
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
