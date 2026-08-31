import SwiftUI

public struct ParentAbsenceMonitoringView: View {
    public var studentName: String
    public var absentHours: Int
    public var maxAllowedAbsentHours: Int
    public var totalYearHours: Int

    public init(studentName: String = "", absentHours: Int = 0, maxAllowedAbsentHours: Int = 247, totalYearHours: Int = 990) {
        self.studentName = studentName
        self.absentHours = absentHours
        self.maxAllowedAbsentHours = maxAllowedAbsentHours
        self.totalYearHours = totalYearHours
    }

    public var progress: Double {
        guard maxAllowedAbsentHours > 0 else { return 0.0 }
        return min(max(Double(absentHours) / Double(maxAllowedAbsentHours), 0.0), 1.0)
    }

    public var isCritical: Bool {
        return progress > 0.8
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if !studentName.isEmpty {
                        Text(studentName)
                            .font(.headline)
                    }

                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text(NSLocalizedString("pending_justifications", comment: ""))
                                .fontWeight(.bold)
                            Spacer()
                            Text("\(absentHours) / \(maxAllowedAbsentHours) h")
                                .font(.caption2)
                                .foregroundColor(isCritical ? .red : .green)
                        }
                        ProgressView(value: progress)
                            .padding(.vertical, 2)
                            .tint(isCritical ? .red : .green)
                    }
                    .padding(.vertical, 4)

                    HStack {
                        Text("Monte Ore Totale")
                        Spacer()
                        Text("\(totalYearHours) h")
                            .fontWeight(.bold)
                    }
                }
            }
            .navigationTitle(NSLocalizedString("parent_dashboard_title", comment: ""))
        }
    }
}
