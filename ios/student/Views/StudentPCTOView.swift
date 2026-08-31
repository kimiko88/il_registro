import SwiftUI

public struct PctoExperienceModel: Identifiable, Equatable {
    public let id: String
    public let companyName: String
    public let tutorName: String
    public let hoursCompleted: String
    public let status: String

    public init(id: String = UUID().uuidString, companyName: String, tutorName: String, hoursCompleted: String, status: String) {
        self.id = id
        self.companyName = companyName
        self.tutorName = tutorName
        self.hoursCompleted = hoursCompleted
        self.status = status
    }
}

public struct StudentPCTOView: View {
    public var completedHours: Double
    public var totalHours: Double
    public var experiences: [PctoExperienceModel]

    public init(completedHours: Double = 120, totalHours: Double = 150, experiences: [PctoExperienceModel] = []) {
        self.completedHours = completedHours
        self.totalHours = totalHours
        self.experiences = experiences
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    VStack(alignment: .leading, spacing: 8) {
                        HStack {
                            Text("PCTO:")
                            Spacer()
                            Text("\(Int(completedHours)) / \(Int(totalHours)) h")
                                .fontWeight(.bold)
                                .foregroundColor(.purple)
                        }
                        ProgressView(value: completedHours, total: totalHours)
                            .tint(.purple)
                    }
                    .padding(.vertical, 4)
                }

                if !experiences.isEmpty {
                    Section(header: Text("Esperienze")) {
                        ForEach(experiences) { item in
                            VStack(alignment: .leading, spacing: 4) {
                                Text(item.companyName)
                                    .fontWeight(.bold)
                                Text(item.tutorName)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                                Text(item.hoursCompleted)
                                    .font(.caption2)
                                    .foregroundColor(.green)
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
