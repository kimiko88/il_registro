import SwiftUI

public struct ParentWorkshopModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let teacherAndSchedule: String
    public let hoursAndCost: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, teacherAndSchedule: String, hoursAndCost: String, statusText: String = "Iscritto") {
        self.id = id
        self.title = title
        self.teacherAndSchedule = teacherAndSchedule
        self.hoursAndCost = hoursAndCost
        self.statusText = statusText
    }
}

public struct ParentWorkshopsEnrollmentView: View {
    public var workshops: [ParentWorkshopModel]

    public init(workshops: [ParentWorkshopModel] = []) {
        self.workshops = workshops
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if workshops.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(workshops) { w in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(w.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(w.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(w.teacherAndSchedule)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(w.hoursAndCost)
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
