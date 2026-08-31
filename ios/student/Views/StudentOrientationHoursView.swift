import SwiftUI

public struct OrientationModuleModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let hours: Int
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, hours: Int, statusText: String = "Convalidato") {
        self.id = id
        self.title = title
        self.hours = hours
        self.statusText = statusText
    }
}

public struct StudentOrientationHoursView: View {
    public var modules: [OrientationModuleModel]
    public var totalHours: Int
    public var targetHours: Int

    public init(modules: [OrientationModuleModel] = [], totalHours: Int = 30, targetHours: Int = 30) {
        self.modules = modules
        self.totalHours = totalHours
        self.targetHours = targetHours
    }

    public var progress: Double {
        guard targetHours > 0 else { return 0.0 }
        return min(max(Double(totalHours) / Double(targetHours), 0.0), 1.0)
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Ore \(totalHours)/\(targetHours)h")
                                .fontWeight(.bold)
                            Spacer()
                            Text("\(Int(progress * 100))%")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        ProgressView(value: progress)
                            .padding(.vertical, 2)

                        ForEach(modules) { mod in
                            Text("\(mod.title) (\(mod.hours)h) • \(mod.statusText)")
                                .font(.subheadline)
                                .foregroundColor(.blue)
                        }
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
