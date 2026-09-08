import SwiftUI

public struct SecretarySafetyDrillModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let timeDetails: String
    public let peopleAndReport: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, timeDetails: String, peopleAndReport: String, statusText: String = "Riuscita") {
        self.id = id
        self.title = title
        self.timeDetails = timeDetails
        self.peopleAndReport = peopleAndReport
        self.statusText = statusText
    }
}

public struct SecretarySafetyView: View {
    public var drills: [SecretarySafetyDrillModel]

    public init(drills: [SecretarySafetyDrillModel] = []) {
        self.drills = drills
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if drills.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(drills) { item in
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
                                Text(item.timeDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.peopleAndReport)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("secretary_dashboard_title", comment: ""))
        }
    }
}
