import SwiftUI

public struct TeacherCreditModel: Identifiable, Equatable {
    public let id: String
    public let studentAndGpa: String
    public let proposalDetails: String
    public let bandBadgeText: String

    public init(id: String = UUID().uuidString, studentAndGpa: String, proposalDetails: String, bandBadgeText: String = "10 - 11 Punti") {
        self.id = id
        self.studentAndGpa = studentAndGpa
        self.proposalDetails = proposalDetails
        self.bandBadgeText = bandBadgeText
    }
}

public struct TeacherCreditsView: View {
    public var credits: [TeacherCreditModel]

    public init(credits: [TeacherCreditModel] = []) {
        self.credits = credits
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if credits.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(credits) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.studentAndGpa)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.bandBadgeText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.proposalDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
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
