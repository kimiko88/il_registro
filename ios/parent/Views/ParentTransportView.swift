import SwiftUI

public struct ParentTransportLineModel: Identifiable, Equatable {
    public let id: String
    public let lineTitle: String
    public let departureStop: String
    public let returnStop: String
    public let driverAndVehicle: String
    public let statusText: String

    public init(id: String = UUID().uuidString, lineTitle: String, departureStop: String, returnStop: String, driverAndVehicle: String, statusText: String = "Regolare") {
        self.id = id
        self.lineTitle = lineTitle
        self.departureStop = departureStop
        self.returnStop = returnStop
        self.driverAndVehicle = driverAndVehicle
        self.statusText = statusText
    }
}

public struct ParentTransportView: View {
    public var lines: [ParentTransportLineModel]

    public init(lines: [ParentTransportLineModel] = []) {
        self.lines = lines
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if lines.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(lines) { line in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(line.lineTitle)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(line.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(line.departureStop)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(line.returnStop)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(line.driverAndVehicle)
                                    .font(.caption)
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
