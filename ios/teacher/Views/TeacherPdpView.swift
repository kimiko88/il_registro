import SwiftUI

public struct TeacherPdpPlanModel: Identifiable, Equatable {
    public let id: String
    public let studentAndType: String
    public let compensatoryTools: String
    public let dispensatoryMeasures: String
    public let statusText: String

    public init(id: String = UUID().uuidString, studentAndType: String, compensatoryTools: String, dispensatoryMeasures: String, statusText: String = "Firmato") {
        self.id = id
        self.studentAndType = studentAndType
        self.compensatoryTools = compensatoryTools
        self.dispensatoryMeasures = dispensatoryMeasures
        self.statusText = statusText
    }
}

public struct TeacherPdpView: View {
    public var plans: [TeacherPdpPlanModel]

    public init(plans: [TeacherPdpPlanModel] = []) {
        self.plans = plans
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if plans.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(plans) { p in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(p.studentAndType)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(p.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(p.compensatoryTools)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(p.dispensatoryMeasures)
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
