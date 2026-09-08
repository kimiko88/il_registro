import SwiftUI

public struct TeacherSubstitutionModel: Identifiable, Equatable {
    public let id: String
    public let hourText: String
    public let className: String
    public let absentTeacherAndLocation: String
    public let activityDetails: String

    public init(id: String = UUID().uuidString, hourText: String, className: String, absentTeacherAndLocation: String, activityDetails: String) {
        self.id = id
        self.hourText = hourText
        self.className = className
        self.absentTeacherAndLocation = absentTeacherAndLocation
        self.activityDetails = activityDetails
    }
}

public struct TeacherDailySubstitutionsView: View {
    public var substitutions: [TeacherSubstitutionModel]

    public init(substitutions: [TeacherSubstitutionModel] = []) {
        self.substitutions = substitutions
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if substitutions.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(substitutions) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.hourText)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.className)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.blue.opacity(0.2))
                                        .foregroundColor(.blue)
                                        .cornerRadius(4)
                                }
                                Text(item.absentTeacherAndLocation)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.activityDetails)
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
