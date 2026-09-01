import SwiftUI

public struct TeacherSeatingLayoutModel: Identifiable, Equatable {
    public let id: String
    public let layoutTitle: String
    public let studentGroupDetails: String
    public let inclusionNotes: String
    public let stationsCountText: String

    public init(id: String = UUID().uuidString, layoutTitle: String, studentGroupDetails: String, inclusionNotes: String, stationsCountText: String = "24 Postazioni") {
        self.id = id
        self.layoutTitle = layoutTitle
        self.studentGroupDetails = studentGroupDetails
        self.inclusionNotes = inclusionNotes
        self.stationsCountText = stationsCountText
    }
}

public struct TeacherSeatingPlanView: View {
    public var layouts: [TeacherSeatingLayoutModel]

    public init(layouts: [TeacherSeatingLayoutModel] = []) {
        self.layouts = layouts
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if layouts.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(layouts) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.layoutTitle)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.stationsCountText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.studentGroupDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.inclusionNotes)
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
